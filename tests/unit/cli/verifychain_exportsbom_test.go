package cli_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// signingKey is a stable 32-byte key used across verify-chain tests.
var signingKey = []byte("test-signing-key-1234567890abcd!")

func signingKeyB64() string {
	return base64.StdEncoding.EncodeToString(signingKey)
}

func writeKeyringFile(t *testing.T, dir string, keys [][]byte) string {
	t.Helper()
	path := filepath.Join(dir, ".skpm-keys")
	var lines []string
	for _, k := range keys {
		lines = append(lines, base64.StdEncoding.EncodeToString(k))
	}
	require.NoError(t, os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644))
	return path
}

func signedLockfile(t *testing.T, dir string, entries []lockfile.SkillLock) string {
	t.Helper()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	for _, sl := range entries {
		if sl.Signature == "" && sl.SHA256 != "" {
			sl.Signature = fmt.Sprintf("%s", cli.SkillSignature(signingKey, sl.Name, sl.Version, sl.SHA256))
		}
		lf.Upsert(sl)
	}
	require.NoError(t, lf.Write(lp))
	return lp
}

// ── verify-chain command ──────────────────────────────────────────────────────

func TestVerifyChainOK(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("a", 64)
	sig := cli.SkillSignature(signingKey, "auth-skill", "1.0.0", sha)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: sha, Signature: sig})
	require.NoError(t, lf.Write(lp))

	keysFile := writeKeyringFile(t, dir, [][]byte{signingKey})

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"verify-chain", "auth-skill", "--lock", lp, "--keys-file", keysFile})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "ok")
}

func TestVerifyChainInvalidSignature(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("a", 64)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: sha, Signature: "badsig=="})
	require.NoError(t, lf.Write(lp))

	keysFile := writeKeyringFile(t, dir, [][]byte{signingKey})

	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-chain", "auth-skill", "--lock", lp, "--keys-file", keysFile})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed verification")
}

func TestVerifyChainUnsigned(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	keysFile := writeKeyringFile(t, dir, [][]byte{signingKey})

	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-chain", "auth-skill", "--lock", lp, "--keys-file", keysFile})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed verification")
}

func TestVerifyChainAll(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("a", 64)
	lp := signedLockfile(t, dir, []lockfile.SkillLock{
		{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: sha},
		{Name: "lint-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64)},
	})
	keysFile := writeKeyringFile(t, dir, [][]byte{signingKey})

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"verify-chain", "--all", "--lock", lp, "--keys-file", keysFile})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "2 skill(s) verified")
}

func TestVerifyChainNoKeysError(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))
	// No keyring file.

	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-chain", "--all", "--lock", lp, "--keys-file", filepath.Join(dir, ".skpm-keys")})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no trusted keys")
}

func TestVerifyChainSkillNotFound(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))
	keysFile := writeKeyringFile(t, dir, [][]byte{signingKey})

	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-chain", "ghost", "--lock", lp, "--keys-file", keysFile})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ghost")
}

func TestVerifyChainRequiresNameOrAll(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))
	keysFile := writeKeyringFile(t, dir, [][]byte{signingKey})

	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-chain", "--lock", lp, "--keys-file", keysFile})
	err := root.Execute()
	require.Error(t, err)
}

// ── export-sbom command ───────────────────────────────────────────────────────

func TestExportSBOMPrintsToStdout(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "auth-skill",
		Version: "1.2.3",
		Source:  "https://registry.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"export-sbom", "--lock", lp})
	require.NoError(t, root.Execute())

	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &doc))
	assert.Equal(t, "agentskills-sbom", doc["format"])
	components, ok := doc["components"].([]any)
	require.True(t, ok)
	require.Len(t, components, 1)
	c := components[0].(map[string]any)
	assert.Equal(t, "auth-skill", c["name"])
	assert.Equal(t, "1.2.3", c["version"])
}

func TestExportSBOMWritesFile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	outPath := filepath.Join(dir, "sbom.json")
	root := cli.NewRootCmd()
	root.SetArgs([]string{"export-sbom", "--lock", lp, "--out", outPath})
	require.NoError(t, root.Execute())

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "auth-skill")
}

func TestExportSBOMIncludesSignature(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("a", 64)
	sig := cli.SkillSignature(signingKey, "auth-skill", "1.0.0", sha)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: sha, Signature: sig})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"export-sbom", "--lock", lp})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), sig)
}

func TestExportSBOMEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"export-sbom", "--lock", lp})
	require.NoError(t, root.Execute())

	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &doc))
	components := doc["components"].([]any)
	assert.Empty(t, components)
}
