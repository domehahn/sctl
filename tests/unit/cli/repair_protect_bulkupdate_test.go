package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── repair: command registered ────────────────────────────────────────────────

func TestRepairCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "repair" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("fix"))
		}
	}
	assert.True(t, found, "repair command should be registered on root")
}

func TestRepairCleanLockfile(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "clean-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"repair", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No issues")
}

func TestRepairDetectsMissingSHA256(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "no-hash",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		// SHA256 intentionally empty
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"repair", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "SHA256")
}

func TestRepairDetectsMissingSource(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "no-source",
		Version: "1.0.0",
		SHA256:  strings.Repeat("b", 64),
		// Source intentionally empty
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"repair", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "source")
}

func TestRepairFixRemovesProblematicEntries(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "good-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lf.Upsert(lockfile.SkillLock{
		Name:    "bad-skill",
		Version: "1.0.0",
		// missing Source and SHA256
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"repair", "--lock", lp, "--fix"})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	// Read back and verify bad-skill was removed.
	fixed, err := lockfile.Read(lp)
	require.NoError(t, err)
	_, found := fixed.Find("good-skill")
	assert.True(t, found, "good-skill should remain")
	_, found = fixed.Find("bad-skill")
	assert.False(t, found, "bad-skill should be removed")
}

func TestRepairJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "broken",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		// missing SHA256
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"repair", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var issues []map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &issues))
	assert.NotEmpty(t, issues)
	assert.Contains(t, issues[0]["skill"], "broken")
}

// ── protect / unprotect ───────────────────────────────────────────────────────

func TestProtectCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	var foundProtect, foundUnprotect bool
	for _, sub := range root.Commands() {
		if sub.Use == "protect <name>" {
			foundProtect = true
			assert.NotNil(t, sub.Flags().Lookup("lock-dir"))
		}
		if sub.Use == "unprotect <name>" {
			foundUnprotect = true
		}
	}
	assert.True(t, foundProtect, "protect command should be registered")
	assert.True(t, foundUnprotect, "unprotect command should be registered")
}

func TestProtectRequiresOneArg(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"protect"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestProtectAndUnprotect(t *testing.T) {
	dir := t.TempDir()

	// Protect a skill.
	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"protect", "my-skill", "--lock-dir", dir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Protected")

	// Verify the protect file was created.
	protectPath := filepath.Join(dir, ".skpm-protect")
	assert.FileExists(t, protectPath)

	// Unprotect the skill.
	buf.Reset()
	root = cli.NewRootCmd()
	root.SetArgs([]string{"unprotect", "my-skill", "--lock-dir", dir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Unprotected")

	// Re-protecting after unprotect should succeed.
	root = cli.NewRootCmd()
	root.SetArgs([]string{"protect", "my-skill", "--lock-dir", dir})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())
}

func TestProtectIdempotent(t *testing.T) {
	dir := t.TempDir()

	run := func() {
		root := cli.NewRootCmd()
		root.SetArgs([]string{"protect", "skill-a", "--lock-dir", dir})
		root.SetOut(&bytes.Buffer{})
		require.NoError(t, root.Execute())
	}
	run()
	run()

	// File should list skill-a exactly once.
	data, err := os.ReadFile(filepath.Join(dir, ".skpm-protect"))
	require.NoError(t, err)
	count := strings.Count(string(data), "skill-a")
	assert.Equal(t, 1, count, "skill should appear exactly once in protect file")
}

func TestUnprotectNonProtectedSkill(t *testing.T) {
	dir := t.TempDir()

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unprotect", "ghost-skill", "--lock-dir", dir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "not protected")
}

func TestProtectDryRun(t *testing.T) {
	dir := t.TempDir()

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"protect", "skill-x", "--lock-dir", dir, "--dry-run"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Would protect")

	// File should not have been created.
	_, err := os.Stat(filepath.Join(dir, ".skpm-protect"))
	assert.True(t, os.IsNotExist(err), "protect file should not be created on dry-run")
}

// ── bulk-update ───────────────────────────────────────────────────────────────

func TestBulkUpdateCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "bulk-update <pattern>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("protect-file"))
		}
	}
	assert.True(t, found, "bulk-update command should be registered on root")
}

func TestBulkUpdateRequiresPattern(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"bulk-update"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestBulkUpdateNoMatchPrints(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "unrelated-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Write a minimal manifest so manifest.Read succeeds.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte("skills: []\n"), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"bulk-update", "lint-*", "--lock", lp})
	root.SetOut(&buf)
	_ = root.Execute()
	t.Logf("output: %s", buf.String())
}

func TestBulkUpdateSkipsProtected(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "protected-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Add to protect store.
	protectPath := filepath.Join(dir, ".skpm-protect")
	require.NoError(t, os.WriteFile(protectPath, []byte("protected:\n  - protected-skill\n"), 0o644))

	// Write minimal manifest so manifest.Read succeeds.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte("skills: []\n"), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"bulk-update", "*", "--lock", lp, "--protect-file", protectPath})
	root.SetOut(&buf)
	_ = root.Execute()
	assert.Contains(t, buf.String(), "protected")
}
