//go:build e2e

package integration

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/attestation"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/require"
)

func TestToolchainGoldenPath(t *testing.T) {
	skcr := os.Getenv("SKCR_BINARY")
	skil := os.Getenv("SKIL_BINARY")
	require.NotEmpty(t, skcr, "SKCR_BINARY is required")
	require.NotEmpty(t, skil, "SKIL_BINARY is required")
	root, err := filepath.Abs("../..")
	require.NoError(t, err)
	work := t.TempDir()
	build := func(output, target string) {
		c := exec.Command("go", "build", "-trimpath", "-o", output, target)
		c.Dir = root
		out, err := c.CombinedOutput()
		require.NoError(t, err, "%s", out)
	}
	skpm := filepath.Join(work, "skpm")
	build(skpm, "./cmd/skpm")
	runtime := filepath.Join(work, "digest-runtime")
	build(runtime, "./toolchain-conformance/runtime")
	env := append(os.Environ(), "XDG_CONFIG_HOME="+filepath.Join(work, "config"), "XDG_CACHE_HOME="+filepath.Join(work, "cache"), "SKPM_REGISTRY_TOKEN=", "GIT_COMMIT=")
	run := func(cwd, bin string, args ...string) []byte {
		c := exec.Command(bin, args...)
		c.Dir = cwd
		c.Env = env
		out, err := c.CombinedOutput()
		if err != nil {
			for i, a := range args {
				if a == "--output" && i+1 < len(args) {
					report, _ := os.ReadFile(args[i+1])
					t.Logf("failure report: %s", report)
				}
			}
		}
		require.NoError(t, err, "%s %v\n%s", bin, args, out)
		return out
	}
	run(work, skcr, "scaffold", "skill", "toolchain-demo", "--output-dir", work, "--version", "1.0.0", "--description", "Compute a SHA256 digest of supplied text without host tools.", "--owner", "conformance")
	source := filepath.Join(work, "toolchain-demo")
	compiledRoot := filepath.Join(work, "compiled")
	run(work, skcr, "compile", source, "--target", "skil", "--output", compiledRoot, "--require-lossless")
	compiled := filepath.Join(compiledRoot, "toolchain-demo")
	run(work, skil, "validate", compiled)
	run(work, skil, "scan", compiled, "--static-only")
	run(work, skil, "verify", compiled)
	run(work, skil, "assure", compiled, "--runtime-command", runtime, "--format", "json", "--output", filepath.Join(work, "assurance.json"))
	run(work, skpm, "package", compiled, "--output-dir", filepath.Join(work, "packages"))
	pkg := filepath.Join(work, "packages", "toolchain-demo-1.0.0.zip")
	published, err := os.ReadFile(pkg)
	require.NoError(t, err)
	sum := sha256.Sum256(published)
	digest := hex.EncodeToString(sum[:])
	// Packaging metadata changes content identity. Attest the final archive, never
	// relabel a directory attestation with an unrelated ZIP digest.
	key := filepath.Join(work, "key.pem")
	var public struct {
		KeyID     string `json:"key_id"`
		PublicKey string `json:"public_key"`
	}
	require.NoError(t, json.Unmarshal(run(work, skil, "key", "generate", "--output", key), &public))
	evidence := filepath.Join(work, "attestation.json")
	run(work, skil, "attest", pkg, "--static-only", "--signing-key", key, "--output", evidence)
	predicate, err := os.ReadFile(evidence)
	require.NoError(t, err)
	require.NoError(t, attestation.BindSubject(predicate, digest))
	_, err = attestation.Verify(predicate, map[string]string{public.KeyID: public.PublicKey})
	require.NoError(t, err)
	base, cleanup := startSkillForge(t, findSkillForgeRepo(t))
	defer cleanup()
	configDir := filepath.Join(work, "config", "skpm")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	config := "default_registry: test\ncache_dir: " + filepath.Join(work, "cache") + "\nregistries:\n  test:\n    type: skillforge\n    namespace: default\n    url: " + base + "\ntrusted_signers:\n  " + public.KeyID + ": " + public.PublicKey + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(config), 0600))
	run(work, skpm, "publish", compiled, "--no-tag", "--no-push", "--no-changelog", "--output-dir", filepath.Join(work, "republished"))
	republished, err := os.ReadFile(filepath.Join(work, "republished", "toolchain-demo-1.0.0.zip"))
	require.NoError(t, err)
	require.Equal(t, published, republished)
	run(work, skpm, "attest", "toolchain-demo@1.0.0", "--file", evidence)
	consumer := filepath.Join(work, "consumer")
	require.NoError(t, os.Mkdir(consumer, 0755))
	run(consumer, skpm, "init")
	run(consumer, skpm, "add", "toolchain-demo@1.0.0", "--source", "test", "--no-install")
	run(consumer, skpm, "lock")
	lockPath := filepath.Join(consumer, "agent-skills.lock")
	locked, err := os.ReadFile(lockPath)
	require.NoError(t, err)
	lf, err := lockfile.Read(lockPath)
	require.NoError(t, err)
	require.Len(t, lf.Skills, 1)
	require.Equal(t, digest, lf.Skills[0].SHA256)
	require.Equal(t, "1.0.0", lf.Skills[0].Version)
	run(consumer, skpm, "install", "--frozen-lockfile")
	run(consumer, skpm, "verify")
	run(consumer, skpm, "attestations", "toolchain-demo@1.0.0", "--verify")
	after, err := os.ReadFile(lockPath)
	require.NoError(t, err)
	require.Equal(t, locked, after)
	// Locate installation from actual extracted manifest, and compare every archive
	// entry, not just SKILL.md or registry metadata.
	installed := ""
	require.NoError(t, filepath.WalkDir(consumer, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Name() == "skill.yaml" {
			installed = filepath.Dir(p)
		}
		return nil
	}))
	require.NotEmpty(t, installed)
	zr, err := zip.NewReader(bytes.NewReader(published), int64(len(published)))
	require.NoError(t, err)
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		r, err := f.Open()
		require.NoError(t, err)
		want, err := io.ReadAll(r)
		require.NoError(t, err)
		require.NoError(t, r.Close())
		got, err := os.ReadFile(filepath.Join(installed, f.Name))
		require.NoError(t, err)
		require.Equal(t, want, got, f.Name)
	}
	run(consumer, skil, "verify", installed)
	run(consumer, skil, "assure", installed, "--runtime-command", runtime, "--format", "json", "--output", filepath.Join(work, "installed-assurance.json"))
	// Wrong signed subject must fail before an HTTP upload.
	wrong := exec.Command(skpm, "attest", "toolchain-demo@1.0.0", "--file", evidence, "--digest", strings.Repeat("b", 64))
	wrong.Dir = consumer
	wrong.Env = env
	out, err := wrong.CombinedOutput()
	require.Error(t, err)
	require.Contains(t, string(out), "ATTESTATION_SUBJECT_MISMATCH")
	t.Logf("AUTHORING → ASSURANCE → PACKAGE → REGISTRY → FROZEN INSTALL → INTEGRITY: %s", digest)
}
