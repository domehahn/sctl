package skill_test

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── test helpers ──────────────────────────────────────────────────────────

func writeSkillFixture(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	return dir
}

var validSkillFiles = map[string]string{
	"SKILL.md":     "# My Skill\nDoes things.",
	"VERSION":      "1.2.3",
	"CHANGELOG.md": "# Changelog\n\n## 1.2.3\n\n- Initial release\n",
	"skill.yaml": `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
  - gitlab-duo
`,
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func assertErrorField(t *testing.T, res *skill.ValidationResult, field string) {
	t.Helper()
	for _, e := range res.Errors {
		if e.Field == field {
			return
		}
	}
	t.Errorf("expected error for field %q, got: %+v", field, res.Errors)
}

// ── Validator tests ───────────────────────────────────────────────────────

func TestValidateValidSkill(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	res, err := skill.NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Empty(t, res.Errors)
}

func TestValidateMissingSkillMD(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "SKILL.md")
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "SKILL.md")
}

func TestValidateEmptySkillMD(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["SKILL.md"] = ""
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "SKILL.md")
}

func TestValidateMissingVERSION(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "VERSION")
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "VERSION")
}

func TestValidateInvalidSemver(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "not-a-version"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "VERSION")
}

func TestValidateVersionMismatch(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "2.0.0"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateUnknownPlatform(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = "name: my-skill\nversion: \"1.2.3\"\ncompatible_with:\n  - unknown-agent\n"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateMissingChangelog(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "CHANGELOG.md")
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.NotEmpty(t, res.Warnings)
}

func TestValidateChangelogVPrefix(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["CHANGELOG.md"] = "# Changelog\n\n## v1.2.3\n\n- added\n"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Empty(t, res.Warnings)
}

func TestValidateSkillYAMLMissingName(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = "version: \"1.2.3\"\ncompatible_with:\n  - claude-code\n"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateSkillYAMLParseError(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = "invalid: [yaml: {{"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateEmptyCompatibleWith(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = "name: my-skill\nversion: \"1.2.3\"\ncompatible_with: []\n"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateMissingDescription(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = "name: my-skill\nversion: \"1.2.3\"\ncompatible_with:\n  - claude-code\n"
	res, err := skill.NewValidator().Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

// ── Packager tests ────────────────────────────────────────────────────────

func TestPackageValidSkill(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	outDir := t.TempDir()
	result, err := skill.NewPackager().Package(context.Background(), dir, outDir)
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.2.3", result.Version)
	assert.NotEmpty(t, result.SHA256)
	assert.FileExists(t, result.OutputPath)

	zr, err := zip.OpenReader(result.OutputPath)
	require.NoError(t, err)
	defer zr.Close()

	fileNames := make(map[string]bool)
	var manifest skill.SkillManifest
	for _, f := range zr.File {
		fileNames[f.Name] = true
		if f.Name == "manifest.json" {
			rc, _ := f.Open()
			json.NewDecoder(rc).Decode(&manifest)
			rc.Close()
		}
	}
	assert.True(t, fileNames["SKILL.md"])
	assert.True(t, fileNames["manifest.json"])
	assert.Equal(t, "my-skill", manifest.Name)
}

func TestPackageInvalidSkill(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "SKILL.md")
	_, err := skill.NewPackager().Package(context.Background(), writeSkillFixture(t, files), t.TempDir())
	require.Error(t, err)
	var vfe *skill.ValidationFailedError
	assert.ErrorAs(t, err, &vfe)
}

func TestPackageOutputDefaultsToSkillDir(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	result, err := skill.NewPackager().Package(context.Background(), dir, "")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "my-skill-1.2.3.zip"), result.OutputPath)
}

func TestPackageExcludesGitDir(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".git"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".git", "HEAD"), []byte("ref: main"), 0o644))

	result, err := skill.NewPackager().Package(context.Background(), dir, t.TempDir())
	require.NoError(t, err)

	zr, err := zip.OpenReader(result.OutputPath)
	require.NoError(t, err)
	defer zr.Close()
	for _, f := range zr.File {
		assert.NotContains(t, f.Name, ".git/")
	}
}

func TestPackageCreatesMissingOutputDir(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	outDir := filepath.Join(t.TempDir(), "dist")
	result, err := skill.NewPackager().Package(context.Background(), dir, outDir)
	require.NoError(t, err)
	assert.FileExists(t, result.OutputPath)
}

func TestShouldSkip(t *testing.T) {
	assert.True(t, skill.ShouldSkip(".git", ".git"))
	assert.True(t, skill.ShouldSkip(".DS_Store", ".DS_Store"))
	assert.True(t, skill.ShouldSkip("file.tmp", "file.tmp"))
	assert.True(t, skill.ShouldSkip("sub", ".git/sub"))
	assert.False(t, skill.ShouldSkip("SKILL.md", "SKILL.md"))
}
