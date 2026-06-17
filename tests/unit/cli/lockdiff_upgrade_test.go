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

// ── lock-diff: command registered ────────────────────────────────────────────

func TestLockDiffCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "lock-diff <before.lock> <after.lock>" {
			found = true
			break
		}
	}
	assert.True(t, found, "lock-diff command should be registered on root")
}

// ── lock-diff: diff logic ─────────────────────────────────────────────────────

func writeLock(t *testing.T, dir, name string, lf *lockfile.LockFile) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, lf.Write(p))
	return p
}

func TestLockDiffDetectsAddedSkill(t *testing.T) {
	dir := t.TempDir()

	before := lockfile.New()
	beforePath := writeLock(t, dir, "before.lock", before)

	after := lockfile.New()
	after.Upsert(lockfile.SkillLock{
		Name: "new-skill", Version: "1.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("a", 64),
	})
	afterPath := writeLock(t, dir, "after.lock", after)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-diff", beforePath, afterPath})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "+")
	assert.Contains(t, buf.String(), "new-skill")
}

func TestLockDiffDetectsRemovedSkill(t *testing.T) {
	dir := t.TempDir()

	before := lockfile.New()
	before.Upsert(lockfile.SkillLock{
		Name: "old-skill", Version: "1.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("a", 64),
	})
	beforePath := writeLock(t, dir, "before.lock", before)

	after := lockfile.New()
	afterPath := writeLock(t, dir, "after.lock", after)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-diff", beforePath, afterPath})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "-")
	assert.Contains(t, buf.String(), "old-skill")
}

func TestLockDiffDetectsVersionChange(t *testing.T) {
	dir := t.TempDir()

	before := lockfile.New()
	before.Upsert(lockfile.SkillLock{
		Name: "my-skill", Version: "1.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("a", 64),
	})
	beforePath := writeLock(t, dir, "before.lock", before)

	after := lockfile.New()
	after.Upsert(lockfile.SkillLock{
		Name: "my-skill", Version: "2.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("b", 64),
	})
	afterPath := writeLock(t, dir, "after.lock", after)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-diff", beforePath, afterPath})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "~")
	assert.Contains(t, out, "my-skill")
	assert.Contains(t, out, "1.0.0")
	assert.Contains(t, out, "2.0.0")
}

func TestLockDiffNoChanges(t *testing.T) {
	dir := t.TempDir()

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name: "stable-skill", Version: "1.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("a", 64),
	})
	p1 := writeLock(t, dir, "a.lock", lf)
	p2 := writeLock(t, dir, "b.lock", lf)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-diff", p1, p2})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "No changes")
}

func TestLockDiffNonExistentBeforeIsEmptyLockfile(t *testing.T) {
	dir := t.TempDir()

	after := lockfile.New()
	after.Upsert(lockfile.SkillLock{
		Name: "brand-new", Version: "1.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("a", 64),
	})
	afterPath := writeLock(t, dir, "after.lock", after)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-diff", filepath.Join(dir, "nonexistent.lock"), afterPath})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "brand-new")
	assert.Contains(t, buf.String(), "+")
}

func TestLockDiffJSONOutput(t *testing.T) {
	dir := t.TempDir()

	before := lockfile.New()
	beforePath := writeLock(t, dir, "before.lock", before)

	after := lockfile.New()
	after.Upsert(lockfile.SkillLock{
		Name: "added-skill", Version: "1.0.0",
		Source: "https://r.example.com", SHA256: strings.Repeat("a", 64),
	})
	afterPath := writeLock(t, dir, "after.lock", after)

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-diff", beforePath, afterPath, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), `"kind"`)
	assert.Contains(t, buf.String(), `"added"`)
	assert.Contains(t, buf.String(), `"added-skill"`)
}

// ── upgrade: command registered ───────────────────────────────────────────────

func TestUpgradeCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "upgrade" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("all"))
			assert.NotNil(t, sub.Flags().Lookup("skill"))
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("install"))
			break
		}
	}
	assert.True(t, found, "upgrade command should be registered on root")
}

func TestUpgradeFailsWithoutManifest(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"upgrade", "--all"})
	assert.Error(t, root.Execute())
}
