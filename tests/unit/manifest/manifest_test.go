package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	mf := manifest.New()
	assert.Equal(t, 1, mf.Version)
	assert.Empty(t, mf.Skills)
}

func TestReadWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifest.DefaultFilename)

	mf := manifest.New()
	mf.Upsert(manifest.SkillEntry{Name: "my-skill", Version: "1.0.0", Source: "company-gitlab"})
	require.NoError(t, mf.Write(path))

	loaded, err := manifest.Read(path)
	require.NoError(t, err)
	assert.Equal(t, 1, loaded.Version)
	require.Len(t, loaded.Skills, 1)
	assert.Equal(t, "my-skill", loaded.Skills[0].Name)
	assert.Equal(t, "1.0.0", loaded.Skills[0].Version)
}

func TestUpsertInsert(t *testing.T) {
	mf := manifest.New()
	mf.Upsert(manifest.SkillEntry{Name: "a", Version: "1.0.0"})
	mf.Upsert(manifest.SkillEntry{Name: "b", Version: "2.0.0"})
	assert.Len(t, mf.Skills, 2)
}

func TestUpsertUpdate(t *testing.T) {
	mf := manifest.New()
	mf.Upsert(manifest.SkillEntry{Name: "a", Version: "1.0.0"})
	mf.Upsert(manifest.SkillEntry{Name: "a", Version: "1.1.0"})
	assert.Len(t, mf.Skills, 1)
	assert.Equal(t, "1.1.0", mf.Skills[0].Version)
}

func TestRemove(t *testing.T) {
	mf := manifest.New()
	mf.Upsert(manifest.SkillEntry{Name: "a"})
	mf.Upsert(manifest.SkillEntry{Name: "b"})

	assert.True(t, mf.Remove("a"))
	assert.False(t, mf.Remove("a"))
	assert.Len(t, mf.Skills, 1)
}

func TestReadMissing(t *testing.T) {
	_, err := manifest.Read("/nonexistent/manifest.yaml")
	assert.Error(t, err)
}

func TestReadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifest.DefaultFilename)
	require.NoError(t, os.WriteFile(path, []byte("invalid: [yaml: {{"), 0o644))
	_, err := manifest.Read(path)
	assert.Error(t, err)
}

func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifest.DefaultFilename)
	mf := manifest.New()
	require.NoError(t, mf.Write(path))

	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err))
}

func TestVersionDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifest.DefaultFilename)
	require.NoError(t, os.WriteFile(path, []byte("skills: []"), 0o644))
	mf, err := manifest.Read(path)
	require.NoError(t, err)
	assert.Equal(t, 1, mf.Version)
}
