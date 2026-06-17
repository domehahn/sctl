package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── move command ──────────────────────────────────────────────────────────────

func TestMoveChangesNamespace(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64), Namespace: "security"})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"move", "auth-skill", "platform", "--lock", lp})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	sl, ok := updated.Find("auth-skill")
	require.True(t, ok)
	assert.Equal(t, "platform", sl.Namespace)
}

func TestMoveSkillNotFound(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"move", "missing-skill", "ns", "--lock", lp})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing-skill")
}

func TestMoveAlreadyInNamespace(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64), Namespace: "platform"})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"move", "auth-skill", "platform", "--lock", lp})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "already in namespace")
}

// ── prune-protected command ───────────────────────────────────────────────────

func TestPruneProtectedRemovesStaleEntry(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	// Protect both skills (one of which is not in the lockfile).
	root := cli.NewRootCmd()
	root.SetArgs([]string{"protect", "auth-skill", "--lock-dir", dir})
	require.NoError(t, root.Execute())
	root2 := cli.NewRootCmd()
	root2.SetArgs([]string{"protect", "ghost-skill", "--lock-dir", dir})
	require.NoError(t, root2.Execute())

	var buf strings.Builder
	root3 := cli.NewRootCmd()
	root3.SetOut(&buf)
	root3.SetArgs([]string{"prune-protected", "--lock", lp, "--lock-dir", dir})
	require.NoError(t, root3.Execute())

	out := buf.String()
	assert.Contains(t, out, "ghost-skill")
	assert.Contains(t, out, "Pruned 1")

	// ghost-skill must be gone from the protect file.
	data, err := os.ReadFile(filepath.Join(dir, ".skpm-protect"))
	require.NoError(t, err)
	assert.NotContains(t, string(data), "ghost-skill")
	assert.Contains(t, string(data), "auth-skill")
}

func TestPruneProtectedNothingToRemove(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"protect", "auth-skill", "--lock-dir", dir})
	require.NoError(t, root.Execute())

	var buf strings.Builder
	root2 := cli.NewRootCmd()
	root2.SetOut(&buf)
	root2.SetArgs([]string{"prune-protected", "--lock", lp, "--lock-dir", dir})
	require.NoError(t, root2.Execute())
	assert.Contains(t, buf.String(), "No stale")
}

func TestPruneProtectedEmptyProtectFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))
	// No protect file created — prune-protected should handle its absence gracefully.

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"prune-protected", "--lock", lp, "--lock-dir", dir})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No protected")
}

// ── tag command ───────────────────────────────────────────────────────────────

func TestTagAddAppendsTag(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	sl, ok := updated.Find("auth-skill")
	require.True(t, ok)
	assert.Equal(t, "production", sl.Metadata["tags"])
}

func TestTagAddSecondTag(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root.Execute())
	root2 := cli.NewRootCmd()
	root2.SetArgs([]string{"tag", "add", "auth-skill", "stable", "--lock", lp})
	require.NoError(t, root2.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	sl, ok := updated.Find("auth-skill")
	require.True(t, ok)
	tags := sl.Metadata["tags"]
	assert.Contains(t, tags, "production")
	assert.Contains(t, tags, "stable")
}

func TestTagAddDuplicate(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root.Execute())

	var buf strings.Builder
	root2 := cli.NewRootCmd()
	root2.SetOut(&buf)
	root2.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root2.Execute())
	assert.Contains(t, buf.String(), "already has tag")
}

func TestTagRemoveDeletesTag(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root.Execute())
	root2 := cli.NewRootCmd()
	root2.SetArgs([]string{"tag", "remove", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root2.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	sl, ok := updated.Find("auth-skill")
	require.True(t, ok)
	assert.Empty(t, sl.Metadata["tags"])
}

func TestTagRemoveNotFound(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"tag", "remove", "auth-skill", "missing-tag", "--lock", lp})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "does not have tag")
}

func TestTagListShowsAllTags(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root.Execute())

	var buf strings.Builder
	root2 := cli.NewRootCmd()
	root2.SetOut(&buf)
	root2.SetArgs([]string{"tag", "list", "--lock", lp})
	require.NoError(t, root2.Execute())
	assert.Contains(t, buf.String(), "production")
}

func TestTagListFilterByTag(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "lint-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "auth-skill", "production", "--lock", lp})
	require.NoError(t, root.Execute())

	var buf strings.Builder
	root2 := cli.NewRootCmd()
	root2.SetOut(&buf)
	root2.SetArgs([]string{"tag", "list", "--tag", "production", "--lock", lp})
	require.NoError(t, root2.Execute())
	out := buf.String()
	assert.Contains(t, out, "auth-skill")
	assert.NotContains(t, out, "lint-skill")
}

func TestTagListEmpty(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"tag", "list", "--lock", lp})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No tags")
}

func TestTagSkillNotFound(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"tag", "add", "ghost", "production", "--lock", lp})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ghost")
}
