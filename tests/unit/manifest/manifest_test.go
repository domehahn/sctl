package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/sklib/spec"
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

// TestSkillEntryIsSpecManifestSkill verifies the type alias is in effect:
// a spec.ManifestSkill can be assigned to manifest.SkillEntry without conversion.
func TestSkillEntryIsSpecManifestSkill(t *testing.T) {
	var entry manifest.SkillEntry = spec.ManifestSkill{
		Name:      "alias-check",
		Version:   "^1.0.0",
		Platforms: []spec.Platform{spec.PlatformClaudeCode},
	}
	assert.Equal(t, "alias-check", entry.Name)
	assert.Equal(t, []spec.Platform{spec.PlatformClaudeCode}, entry.Platforms)
}

// TestRoundtripPlatforms verifies that Platforms []spec.Platform in SkillEntry
// survives a write-read cycle with the correct YAML key.
func TestOverrideFor(t *testing.T) {
	mf := manifest.New()
	mf.Overrides = map[string]manifest.SkillOverride{
		"my-skill": {
			Prepend: "Always use Chinese sources first.",
			Append:  "Cross-validate before reporting.",
			Env:     map[string]string{"EXA_KEY": "env:EXA_KEY"},
		},
	}

	ov := mf.OverrideFor("my-skill")
	require.NotNil(t, ov)
	assert.Equal(t, "Always use Chinese sources first.", ov.Prepend)
	assert.Equal(t, "Cross-validate before reporting.", ov.Append)
	assert.Equal(t, "env:EXA_KEY", ov.Env["EXA_KEY"])

	// non-existent skill returns nil
	assert.Nil(t, mf.OverrideFor("nonexistent"))

	// empty overrides map returns nil
	mf2 := manifest.New()
	assert.Nil(t, mf2.OverrideFor("anything"))
}

func TestWriteOverride(t *testing.T) {
	dir := t.TempDir()

	ov := &manifest.SkillOverride{
		Prepend: "Search Chinese first.",
	}

	err := manifest.WriteOverride(dir, ov)
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, manifest.OverrideFilename))
	require.NoError(t, err)
	assert.Contains(t, string(data), "Search Chinese first.")
}

func TestWriteOverrideNil(t *testing.T) {
	dir := t.TempDir()

	// write a stale file first
	stale := filepath.Join(dir, manifest.OverrideFilename)
	require.NoError(t, os.WriteFile(stale, []byte("stale"), 0o644))

	// nil override removes it
	err := manifest.WriteOverride(dir, nil)
	require.NoError(t, err)
	_, err = os.Stat(stale)
	assert.True(t, os.IsNotExist(err))
}

func TestWriteOverrideEmpty(t *testing.T) {
	dir := t.TempDir()

	// empty override removes stale file
	stale := filepath.Join(dir, manifest.OverrideFilename)
	require.NoError(t, os.WriteFile(stale, []byte("stale"), 0o644))

	err := manifest.WriteOverride(dir, &manifest.SkillOverride{})
	require.NoError(t, err)
	_, err = os.Stat(stale)
	assert.True(t, os.IsNotExist(err))
}

func TestOverrideRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifest.DefaultFilename)

	mf := manifest.New()
	mf.Skills = []manifest.SkillEntry{
		{Name: "agent-reach", Source: "github"},
	}
	mf.Overrides = map[string]manifest.SkillOverride{
		"agent-reach": {
			Prepend: "优先中文",
			Env:     map[string]string{"LANG": "zh-CN"},
		},
	}
	require.NoError(t, mf.Write(path))

	loaded, err := manifest.Read(path)
	require.NoError(t, err)
	assert.Len(t, loaded.Skills, 1)
	assert.Equal(t, "优先中文", loaded.Overrides["agent-reach"].Prepend)
	assert.Equal(t, "zh-CN", loaded.Overrides["agent-reach"].Env["LANG"])
}

func TestRoundtripPlatforms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifest.DefaultFilename)

	mf := manifest.New()
	mf.Upsert(manifest.SkillEntry{
		Name:      "plat-skill",
		Version:   "1.0.0",
		Platforms: []spec.Platform{spec.PlatformClaudeCode, spec.PlatformCursor},
	})
	require.NoError(t, mf.Write(path))

	loaded, err := manifest.Read(path)
	require.NoError(t, err)
	require.Len(t, loaded.Skills, 1)
	assert.Equal(t, []spec.Platform{spec.PlatformClaudeCode, spec.PlatformCursor}, loaded.Skills[0].Platforms)
}
