package lockfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/sklib/spec"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "github",
		DownloadURL: "https://github.com/org/skills/releases/download/v1.0.0/my-skill-1.0.0.zip",
		SHA256:      "abc123",
		InstalledTo: []string{"skills/my-skill"},
	})

	require.NoError(t, lf.Write(path))

	loaded, err := lockfile.Read(path)
	require.NoError(t, err)
	assert.Equal(t, 1, loaded.Version)
	require.Len(t, loaded.Skills, 1)
	assert.Equal(t, "my-skill", loaded.Skills[0].Name)
	assert.Equal(t, "1.0.0", loaded.Skills[0].Version)
}

func TestUpsert(t *testing.T) {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.0.0"})
	lf.Upsert(lockfile.SkillLock{Name: "skill-b", Version: "2.0.0"})
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.1.0"})

	assert.Len(t, lf.Skills, 2)
	s, ok := lf.Find("skill-a")
	require.True(t, ok)
	assert.Equal(t, "1.1.0", s.Version)
}

func TestFind(t *testing.T) {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "x", Version: "0.1.0"})

	s, ok := lf.Find("x")
	require.True(t, ok)
	assert.Equal(t, "0.1.0", s.Version)

	_, ok = lf.Find("nonexistent")
	assert.False(t, ok)
}

func TestRemove(t *testing.T) {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "a"})
	lf.Upsert(lockfile.SkillLock{Name: "b"})

	assert.True(t, lf.Remove("a"))
	assert.False(t, lf.Remove("a"))
	assert.Len(t, lf.Skills, 1)
	assert.Equal(t, "b", lf.Skills[0].Name)
}

func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)

	lf := lockfile.New()
	require.NoError(t, lf.Write(path))

	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err))
}

func TestReadMissing(t *testing.T) {
	_, err := lockfile.Read("/nonexistent/path/agent-skills.lock")
	assert.Error(t, err)
}

func TestReadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)
	require.NoError(t, os.WriteFile(path, []byte("invalid: [yaml: {{"), 0o644))
	_, err := lockfile.Read(path)
	assert.Error(t, err)
}

func TestWriteCreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", lockfile.DefaultFilename)
	lf := lockfile.New()
	require.NoError(t, lf.Write(path))
	assert.FileExists(t, path)
}

func TestVersionDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)
	require.NoError(t, os.WriteFile(path, []byte("skills: []"), 0o644))
	lf, err := lockfile.Read(path)
	require.NoError(t, err)
	assert.Equal(t, 1, lf.Version)
}

func TestUnsupportedVersionRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)
	require.NoError(t, os.WriteFile(path, []byte("version: 99\nskills: []"), 0o644))
	_, err := lockfile.Read(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported lockfile schema version")
}

// TestRoundtripResolvedAt verifies that ResolvedAt/GeneratedBy survive a write-read
// cycle with the correct YAML keys (resolved_at / generated_by, not generated_at).
func TestRoundtripResolvedAt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)

	lf := lockfile.New()
	lf.ResolvedAt = "2026-01-01T00:00:00Z"
	lf.GeneratedBy = "skpm"
	require.NoError(t, lf.Write(path))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	raw := string(data)
	assert.Contains(t, raw, "resolved_at:")
	assert.Contains(t, raw, "generated_by:")
	assert.NotContains(t, raw, "generated_at:")

	loaded, err := lockfile.Read(path)
	require.NoError(t, err)
	assert.Equal(t, "2026-01-01T00:00:00Z", loaded.ResolvedAt)
	assert.Equal(t, "skpm", loaded.GeneratedBy)
}

// TestCompatibleWithPlatforms verifies that CompatibleWith []spec.Platform survives
// a write-read cycle and is still typed as []spec.Platform (not []string).
func TestCompatibleWithPlatforms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockfile.DefaultFilename)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:           "plat-skill",
		Version:        "1.0.0",
		Source:         "github",
		SHA256:         "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		CompatibleWith: []spec.Platform{spec.PlatformClaudeCode, spec.PlatformCursor},
	})
	require.NoError(t, lf.Write(path))

	loaded, err := lockfile.Read(path)
	require.NoError(t, err)
	require.Len(t, loaded.Skills, 1)
	assert.Equal(t, []spec.Platform{spec.PlatformClaudeCode, spec.PlatformCursor}, loaded.Skills[0].CompatibleWith)
}

// TestSkillLockIsSpecLockedSkill verifies the type alias is in effect:
// a spec.LockedSkill can be assigned to lockfile.SkillLock without conversion.
func TestSkillLockIsSpecLockedSkill(t *testing.T) {
	var sl lockfile.SkillLock = spec.LockedSkill{Name: "alias-check", Version: "0.1.0", Source: "local"}
	assert.Equal(t, "alias-check", sl.Name)
}
