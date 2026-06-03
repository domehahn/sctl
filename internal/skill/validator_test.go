package skill

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	"SKILL.md":    "# My Skill\nDoes things.",
	"VERSION":     "1.2.3",
	"CHANGELOG.md": "# Changelog\n\n## 1.2.3\n\n- Initial release\n",
	"skill.yaml": `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
  - gitlab-duo
`,
}

func TestValidateValidSkill(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	v := NewValidator()
	res, err := v.Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Empty(t, res.Errors)
}

func TestValidateMissingSkillMD(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "SKILL.md")
	dir := writeSkillFixture(t, files)

	res, err := NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "SKILL.md")
}

func TestValidateMissingVERSION(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "VERSION")
	dir := writeSkillFixture(t, files)

	res, err := NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "VERSION")
}

func TestValidateInvalidSemver(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "not-a-version"
	dir := writeSkillFixture(t, files)

	res, err := NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "VERSION")
}

func TestValidateVersionMismatch(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "2.0.0"
	dir := writeSkillFixture(t, files)

	res, err := NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateUnknownPlatform(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
compatible_with:
  - unknown-agent
`
	dir := writeSkillFixture(t, files)

	res, err := NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "skill.yaml")
}

func TestValidateMissingChangelog(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "CHANGELOG.md")
	dir := writeSkillFixture(t, files)

	res, err := NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.NotEmpty(t, res.Warnings)
}

func assertErrorField(t *testing.T, res *ValidationResult, field string) {
	t.Helper()
	for _, e := range res.Errors {
		if e.Field == field {
			return
		}
	}
	t.Errorf("expected error for field %q, got errors: %+v", field, res.Errors)
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
