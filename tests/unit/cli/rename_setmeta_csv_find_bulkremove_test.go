package cli_test

import (
	"encoding/csv"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── rename command ────────────────────────────────────────────────────────────

func TestRenameChangesName(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "old-auth", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64), Namespace: "security"})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"rename", "old-auth", "new-auth", "--lock", lp})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	_, oldExists := updated.Find("old-auth")
	assert.False(t, oldExists, "old name should be gone")
	sl, newExists := updated.Find("new-auth")
	require.True(t, newExists, "new name should exist")
	assert.Equal(t, "1.0.0", sl.Version)
	assert.Equal(t, "security", sl.Namespace)
}

func TestRenameOldNotFound(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"rename", "ghost", "new-name", "--lock", lp})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ghost")
}

func TestRenameNewAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "old-auth", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "new-auth", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"rename", "old-auth", "new-auth", "--lock", lp})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// ── set-meta command ──────────────────────────────────────────────────────────

func TestSetMetaSetsArbitraryKey(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"set-meta", "auth-skill", "owner", "platform-team", "--lock", lp})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	sl, ok := updated.Find("auth-skill")
	require.True(t, ok)
	assert.Equal(t, "platform-team", sl.Metadata["owner"])
}

func TestSetMetaIdempotent(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64),
		Metadata: map[string]string{"owner": "platform-team"},
	})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"set-meta", "auth-skill", "owner", "platform-team", "--lock", lp})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "already")
}

func TestSetMetaSkillNotFound(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"set-meta", "ghost", "owner", "team", "--lock", lp})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ghost")
}

// ── export-csv command ────────────────────────────────────────────────────────

func TestExportCSVHeaders(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.2.3", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64), Namespace: "security"})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"export-csv", "--lock", lp})
	require.NoError(t, root.Execute())

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(records), 2)

	assert.Equal(t, []string{"name", "version", "source", "sha256", "namespace", "tags"}, records[0])
	assert.Equal(t, "auth-skill", records[1][0])
	assert.Equal(t, "1.2.3", records[1][1])
	assert.Equal(t, "security", records[1][4])
}

func TestExportCSVEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"export-csv", "--lock", lp})
	require.NoError(t, root.Execute())

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 1, "only the header row")
}

// ── find command ──────────────────────────────────────────────────────────────

func TestFindMatchesGlob(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-login", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "auth-logout", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "lint-check", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("c", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"find", "auth-*", "--lock", lp})
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "auth-login")
	assert.Contains(t, out, "auth-logout")
	assert.NotContains(t, out, "lint-check")
}

func TestFindNoMatch(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"find", "xyz-*", "--lock", lp})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No skills match")
}

func TestFindFilterByNamespace(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64), Namespace: "security"})
	lf.Upsert(lockfile.SkillLock{Name: "auth-lite", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64), Namespace: "platform"})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"find", "auth-*", "--namespace", "security", "--lock", lp})
	require.NoError(t, root.Execute())
	out := buf.String()
	assert.Contains(t, out, "auth-skill")
	assert.NotContains(t, out, "auth-lite")
}

// ── bulk-remove command ───────────────────────────────────────────────────────

func TestBulkRemoveMatchingSkills(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "old-auth", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "old-lint", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "keep-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("c", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"bulk-remove", "old-*", "--lock", lp, "--lock-dir", dir})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	_, hasOldAuth := updated.Find("old-auth")
	_, hasOldLint := updated.Find("old-lint")
	_, hasKeep := updated.Find("keep-skill")
	assert.False(t, hasOldAuth)
	assert.False(t, hasOldLint)
	assert.True(t, hasKeep)
}

func TestBulkRemoveSkipsProtected(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "old-auth", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "old-lint", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64)})
	require.NoError(t, lf.Write(lp))

	// Protect old-auth.
	protectRoot := cli.NewRootCmd()
	protectRoot.SetArgs([]string{"protect", "old-auth", "--lock-dir", dir})
	require.NoError(t, protectRoot.Execute())

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"bulk-remove", "old-*", "--lock", lp, "--lock-dir", dir})
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "protected")
	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	_, hasOldAuth := updated.Find("old-auth")
	_, hasOldLint := updated.Find("old-lint")
	assert.True(t, hasOldAuth, "protected skill should remain")
	assert.False(t, hasOldLint, "unprotected skill should be removed")
}

func TestBulkRemoveDryRun(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "old-auth", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"bulk-remove", "old-*", "--lock", lp, "--lock-dir", dir, "--dry-run"})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	_, still := updated.Find("old-auth")
	assert.True(t, still, "dry-run must not remove skill")
}

func TestBulkRemoveNoMatch(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"bulk-remove", "xyz-*", "--lock", lp, "--lock-dir", dir})
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No skills match")
}
