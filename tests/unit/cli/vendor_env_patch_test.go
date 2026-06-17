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

// ── vendor: command registered ────────────────────────────────────────────────

func TestVendorCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "vendor" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("dir"))
			assert.NotNil(t, sub.Flags().Lookup("clean"))
			break
		}
	}
	assert.True(t, found, "vendor command should be registered on root")
}

func TestVendorCopiesInstalledSkills(t *testing.T) {
	dir := t.TempDir()

	// Create a fake installed skill directory.
	skillDir := filepath.Join(dir, ".claude", "skills", "my-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# My Skill"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "skill.yaml"), []byte("name: my-skill\n"), 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	vendorDir := filepath.Join(dir, "vendor", "skills")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"vendor", "--lock", lockPath, "--dir", vendorDir})
	require.NoError(t, root.Execute())

	// Verify the vendor directory contains the skill files.
	assert.FileExists(t, filepath.Join(vendorDir, "my-skill", "SKILL.md"))
	assert.FileExists(t, filepath.Join(vendorDir, "my-skill", "skill.yaml"))
}

func TestVendorSkipsSkillWithNoInstallPath(t *testing.T) {
	dir := t.TempDir()

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "ghost-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{"/nonexistent/path"},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	vendorDir := filepath.Join(dir, "vendor", "skills")

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"vendor", "--lock", lockPath, "--dir", vendorDir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "skipping")
}

func TestVendorCleanRemovesExistingDir(t *testing.T) {
	dir := t.TempDir()

	skillDir := filepath.Join(dir, ".claude", "skills", "my-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Skill"), 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	vendorDir := filepath.Join(dir, "vendor", "skills")
	require.NoError(t, os.MkdirAll(vendorDir, 0o755))
	staleFile := filepath.Join(vendorDir, "stale-file.txt")
	require.NoError(t, os.WriteFile(staleFile, []byte("old"), 0o644))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"vendor", "--lock", lockPath, "--dir", vendorDir, "--clean"})
	require.NoError(t, root.Execute())

	// Stale file should be gone after --clean.
	assert.NoFileExists(t, staleFile)
	assert.FileExists(t, filepath.Join(vendorDir, "my-skill", "SKILL.md"))
}

func TestVendorDryRunDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	vendorDir := filepath.Join(dir, "vendor", "skills")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"vendor", "--dry-run", "--lock", lockPath, "--dir", vendorDir})
	require.NoError(t, root.Execute())

	_, err := os.Stat(vendorDir)
	assert.True(t, os.IsNotExist(err), "dry-run should not create vendor directory")
}

// ── env --export ──────────────────────────────────────────────────────────────

func writeTestConfig(t *testing.T, content string) {
	t.Helper()
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, "skpm")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(content), 0o600))
	t.Setenv("XDG_CONFIG_HOME", dir)
}

func TestEnvExportPrintsTokens(t *testing.T) {
	writeTestConfig(t, `
default_registry: my-reg
registries:
  my-reg:
    type: generic
    url: https://registry.example.com
    token: supersecret
`)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"env", "--export"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "export SKPM_TOKEN_MY_REG=")
	assert.Contains(t, out, "supersecret")
}

func TestEnvExportPrintsCommentWhenNoTokens(t *testing.T) {
	writeTestConfig(t, `
registries:
  no-token-reg:
    type: generic
    url: https://registry.example.com
`)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"env", "--export"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "# no registry tokens")
}

func TestEnvExportFlagRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	for _, sub := range root.Commands() {
		if sub.Use == "env" {
			assert.NotNil(t, sub.Flags().Lookup("export"), "env should have --export flag")
			return
		}
	}
	t.Error("env command not found")
}

// ── patch: command registered ─────────────────────────────────────────────────

func TestPatchCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "patch <skill-name>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("file"))
			assert.NotNil(t, sub.Flags().Lookup("strip"))
			assert.NotNil(t, sub.Flags().Lookup("reverse"))
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			break
		}
	}
	assert.True(t, found, "patch command should be registered on root")
}

func TestPatchFailsWhenSkillNotInLockfile(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	patchPath := filepath.Join(dir, "fix.patch")
	require.NoError(t, os.WriteFile(patchPath, []byte("--- a/file\n+++ b/file\n"), 0o644))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"patch", "nonexistent-skill", "--lock", lockPath, "--file", patchPath})
	assert.Error(t, root.Execute())
}

func TestPatchFailsWhenPatchFileNotFound(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "my-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"patch", "my-skill", "--lock", lockPath, "--file", "/nonexistent/fix.patch"})
	assert.Error(t, root.Execute())
}

func TestPatchFailsWhenNoInstallDir(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{"/nonexistent/path"},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	patchPath := filepath.Join(dir, "fix.patch")
	require.NoError(t, os.WriteFile(patchPath, []byte("--- a/file\n+++ b/file\n"), 0o644))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"patch", "my-skill", "--lock", lockPath, "--file", patchPath})
	assert.Error(t, root.Execute())
}
