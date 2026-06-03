package lockfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultFilename)

	lf := New()
	lf.Upsert(SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "github",
		SourceURL:   "https://github.com/org/skills/releases/download/v1.0.0/my-skill-1.0.0.zip",
		SHA256:      "abc123",
		InstalledTo: []string{"skills/my-skill"},
	})

	require.NoError(t, lf.Write(path))

	loaded, err := Read(path)
	require.NoError(t, err)
	assert.Equal(t, currentVersion, loaded.Version)
	require.Len(t, loaded.Skills, 1)
	assert.Equal(t, "my-skill", loaded.Skills[0].Name)
	assert.Equal(t, "1.0.0", loaded.Skills[0].Version)
}

func TestUpsert(t *testing.T) {
	lf := New()
	lf.Upsert(SkillLock{Name: "skill-a", Version: "1.0.0"})
	lf.Upsert(SkillLock{Name: "skill-b", Version: "2.0.0"})
	lf.Upsert(SkillLock{Name: "skill-a", Version: "1.1.0"})

	assert.Len(t, lf.Skills, 2)
	s, ok := lf.Find("skill-a")
	require.True(t, ok)
	assert.Equal(t, "1.1.0", s.Version)
}

func TestFind(t *testing.T) {
	lf := New()
	lf.Upsert(SkillLock{Name: "x", Version: "0.1.0"})

	s, ok := lf.Find("x")
	require.True(t, ok)
	assert.Equal(t, "0.1.0", s.Version)

	_, ok = lf.Find("nonexistent")
	assert.False(t, ok)
}

func TestRemove(t *testing.T) {
	lf := New()
	lf.Upsert(SkillLock{Name: "a"})
	lf.Upsert(SkillLock{Name: "b"})

	assert.True(t, lf.Remove("a"))
	assert.False(t, lf.Remove("a"))
	assert.Len(t, lf.Skills, 1)
	assert.Equal(t, "b", lf.Skills[0].Name)
}

func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultFilename)

	lf := New()
	require.NoError(t, lf.Write(path))

	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err), "tmp file should not exist after write")
}

func TestReadMissing(t *testing.T) {
	_, err := Read("/nonexistent/path/agent-skills.lock")
	assert.Error(t, err)
}
