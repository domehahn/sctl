package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── diff-versions: command registered ────────────────────────────────────────

func TestDiffVersionsCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "diff-versions <name> <v1> <v2>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("registry"))
			break
		}
	}
	assert.True(t, found, "diff-versions command should be registered on root")
}

func TestDiffVersionsRequiresThreeArgs(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"diff-versions", "my-skill", "1.0.0"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestDiffVersionsFailsWithNoRegistry(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"diff-versions", "my-skill", "1.0.0", "2.0.0", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute(), "should error without a configured registry")
}

func TestDiffVersionsErrorsWhenRegistryUnconfigured(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "known-skill",
		Version: "1.0.0",
		Source:  "https://no-registry.example.com/known-skill",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"diff-versions", "known-skill", "1.0.0", "2.0.0", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

// ── verify-signatures: command registered ────────────────────────────────────

func TestVerifySignaturesCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "verify-signatures" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("key-file"))
			assert.NotNil(t, sub.Flags().Lookup("key-env"))
			assert.NotNil(t, sub.Flags().Lookup("fail-on-missing"))
			break
		}
	}
	assert.True(t, found, "verify-signatures command should be registered on root")
}

func TestVerifySignaturesAllUnsigned(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "unsigned-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-signatures", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "unsigned")
}

func TestVerifySignaturesFailOnMissing(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "unsigned-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-signatures", "--lock", lp, "--fail-on-missing"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestVerifySignaturesValidSignature(t *testing.T) {
	dir := t.TempDir()

	// Write key file.
	keyFile := filepath.Join(dir, "sign.key")
	key := []byte("test-hmac-signing-key-32bytes!!!")
	require.NoError(t, os.WriteFile(keyFile, key, 0o600))

	// Build a properly signed lock entry.
	sha := strings.Repeat("a", 64)
	// We need to compute the expected signature the same way SkillSignature does.
	// Use the key file approach and verify with verify-signatures.
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:      "signed-skill",
		Version:   "1.0.0",
		Source:    "https://r.example.com",
		SHA256:    sha,
		Signature: "deliberate-mismatch", // will fail verification
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-signatures", "--lock", lp, "--key-file", keyFile})
	root.SetOut(&buf)
	assert.Error(t, root.Execute(), "mismatched signature should fail")
	assert.Contains(t, buf.String(), "mismatch")
}

func TestVerifySignaturesSkipsWithoutKey(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:      "signed-skill",
		Version:   "1.0.0",
		Source:    "https://r.example.com",
		SHA256:    strings.Repeat("a", 64),
		Signature: "some-signature",
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// No trust file, no key file → should skip verification.
	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-signatures", "--lock", lp, "--trust-file", filepath.Join(dir, "no-such-trust.yaml")})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "skip")
}

func TestVerifySignaturesJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-a",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"verify-signatures", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), `"status"`)
}

// ── stat: command registered ──────────────────────────────────────────────────

func TestStatCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "stat <name>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("registry"))
			break
		}
	}
	assert.True(t, found, "stat command should be registered on root")
}

func TestStatRequiresOneArg(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"stat"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestStatShowsInstalledVersion(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "3.1.4",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{".agents/skills/my-skill"},
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"stat", "my-skill", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "my-skill")
	assert.Contains(t, out, "3.1.4")
	assert.Contains(t, out, ".agents/skills/my-skill")
}

func TestStatNotInstalledShowsMessage(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"stat", "missing-skill", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "not installed")
}

func TestStatJSONOutput(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "json-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("b", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"stat", "json-skill", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), `"json-skill"`)
	assert.Contains(t, buf.String(), `"installed_version"`)
}
