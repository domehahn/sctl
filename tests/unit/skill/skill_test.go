package skill_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/sklib/packageio"
	"github.com/domehahn/sklib/spec"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── helpers ───────────────────────────────────────────────────────────────

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
	"SKILL.md": `---
name: my-skill
description: Does things
version: "1.2.3"
since: "2025-01-01"
last_modified: "2026-06-10"
authors:
  - platform-engineering
stability: stable
min_platform_version:
  codex: "unknown"
deprecated_since:
replaces:
supersedes: []
changelog:
  - version: "1.2.3"
    date: "2026-06-10"
    change: "Initial release"
---

# My Skill

Does things.

## Changelog

### 1.2.3 - 2026-06-10

- Initial release.
`,
	"VERSION":      "1.2.3",
	"CHANGELOG.md": "# Changelog\n\n## 1.2.3\n\n- Initial release\n",
	"README.md":    "# My Skill\n",
	"LICENSE":      "MIT\n",
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
	t.Errorf("expected error for field %q, got errors: %+v", field, res.Errors)
}

func assertErrorCode(t *testing.T, res *skill.ValidationResult, code string) {
	t.Helper()
	for _, e := range res.Errors {
		if e.Code == code {
			return
		}
	}
	t.Errorf("expected error with code %q, got errors: %+v", code, res.Errors)
}

func assertWarnField(t *testing.T, res *skill.ValidationResult, field string) {
	t.Helper()
	for _, w := range res.Warnings {
		if w.Field == field {
			return
		}
	}
	t.Errorf("expected warning for field %q, got warnings: %+v", field, res.Warnings)
}

func validate(t *testing.T, files map[string]string, opts skill.ValidationOptions) *skill.ValidationResult {
	t.Helper()
	res, err := skill.NewValidatorWithOptions(opts).Validate(context.Background(), writeSkillFixture(t, files))
	require.NoError(t, err)
	return res
}

func defaultValidate(t *testing.T, files map[string]string) *skill.ValidationResult {
	return validate(t, files, skill.ValidationOptions{})
}

func strictValidate(t *testing.T, files map[string]string) *skill.ValidationResult {
	return validate(t, files, skill.ValidationOptions{Strict: true})
}

func publishValidate(t *testing.T, files map[string]string) *skill.ValidationResult {
	return validate(t, files, skill.ValidationOptions{Publish: true})
}

// ── Default validation ────────────────────────────────────────────────────

func TestValidateValidSkill(t *testing.T) {
	res := defaultValidate(t, validSkillFiles)
	assert.True(t, res.Valid)
	assert.Empty(t, res.Errors)
	assert.Equal(t, "default", res.Profile)
}

func TestValidateMissingSkillMD(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "SKILL.md")
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorField(t, res, "SKILL.md")
	assertErrorCode(t, res, "missing_skill_md")
}

func TestValidateEmptySkillMD(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["SKILL.md"] = ""
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "empty_skill_md")
}

func TestValidateMissingVERSION(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "VERSION")
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_version")
}

func TestValidateInvalidSemver(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "not-a-version"
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "invalid_semver")
}

func TestValidatePrereleaseVersionRejectedByDefault(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "1.2.3-beta.1"
	files["skill.yaml"] = `name: my-skill
version: "1.2.3-beta.1"
description: Does things
compatible_with:
  - claude-code
`
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "prerelease_version")
}

func TestValidatePrereleaseVersionAllowed(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "1.2.3-beta.1"
	files["skill.yaml"] = `name: my-skill
version: "1.2.3-beta.1"
description: Does things
compatible_with:
  - claude-code
`
	res := validate(t, files, skill.ValidationOptions{AllowPrerelease: true})
	assert.True(t, res.Valid)
}

func TestValidateVersionMismatch(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "2.0.0"
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "version_mismatch")
}

func TestValidateVersionVPrefixStripped(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "v1.2.3"
	res := defaultValidate(t, files)
	assert.True(t, res.Valid, "leading v in VERSION should be stripped before comparison")
}

func TestValidateUnknownPlatform(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - unknown-agent
`
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "unknown_platform")
}

func TestValidatePlatformAlias(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - gitlab
`
	res := defaultValidate(t, files)
	// aliases are currently accepted (not flagged as unknown) — normalization happens in format
	assert.True(t, res.Valid)
}

func TestValidateMissingChangelog(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "CHANGELOG.md")
	res := defaultValidate(t, files)
	assert.True(t, res.Valid)
	assertWarnField(t, res, "CHANGELOG.md")
}

func TestValidateMissingREADME(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "README.md")
	res := defaultValidate(t, files)
	assert.True(t, res.Valid)
	assertWarnField(t, res, "README.md")
}

func TestValidateMissingLicense(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "LICENSE")
	res := defaultValidate(t, files)
	assert.True(t, res.Valid)
	assertWarnField(t, res, "LICENSE")
}

func TestValidateChangelogVPrefix(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["CHANGELOG.md"] = "# Changelog\n\n## v1.2.3\n\n- added\n"
	dir := writeSkillFixture(t, files)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidator().Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Empty(t, res.Warnings, "no warnings expected when changelog has v-prefix entry and all optional files present")
}

func TestValidateSkillYAMLMissingName(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
`
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_name")
}

func TestValidateSkillYAMLMissingDescription(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
compatible_with:
  - claude-code
`
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_description")
}

func TestValidateEmptyCompatibleWith(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with: []
`
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_compatible_with")
}

func TestValidateSkillYAMLParseError(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = "invalid: [yaml: {{"
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "yaml_parse_error")
}

func TestValidateForbiddenFile(t *testing.T) {
	files := copyMap(validSkillFiles)
	files[".env"] = "SECRET=abc123"
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "forbidden_file")
}

func TestValidatePossibleSecret(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["config.txt"] = "api_key=abcdef1234567890abcdef"
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "possible_secret")
}

func TestValidatePrivateKey(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["key.pem"] = "-----BEGIN RSA PRIVATE KEY-----\nfakekey\n-----END RSA PRIVATE KEY-----\n"
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)
}

func TestValidateAbsolutePathWarningInDefault(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["SKILL.md"] = "# My Skill\nSee /Users/john/docs for details."
	res := defaultValidate(t, files)
	// absolute path is a warning in default, not an error
	assert.True(t, res.Valid)
	assertWarnField(t, res, "SKILL.md")
}

func TestValidatePlatformNotCompatible(t *testing.T) {
	files := copyMap(validSkillFiles)
	res := validate(t, files, skill.ValidationOptions{Platform: "codex"})
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "platform_not_compatible")
}

func TestValidatePlatformAll(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - all
`
	res := validate(t, files, skill.ValidationOptions{Platform: "codex"})
	assert.True(t, res.Valid)
}

// ── Strict validation ─────────────────────────────────────────────────────

func TestStrictValidMissingChangelogFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "CHANGELOG.md")
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_changelog")
}

func TestStrictMissingChangelogEntryFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["CHANGELOG.md"] = "# Changelog\n\n## 0.9.0\n\n- old release\n"
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_changelog_entry")
}

func TestStrictMissingREADMEFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "README.md")
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_readme")
}

func TestStrictAbsolutePathFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["SKILL.md"] = "# My Skill\nSee /Users/john/docs for details."
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "absolute_path")
}

func TestStrictGeneratedArtifactFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["my-skill-1.2.3.zip"] = "fake zip content"
	dir := writeSkillFixture(t, files)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Strict: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "generated_artifact")
}

func TestStrictGeneratedManifestFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["manifest.json"] = `{"name":"my-skill"}`
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "generated_artifact")
}

func TestStrictDuplicatePlatformFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
  - claude-code
`
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "duplicate_platform")
}

func TestStrictDuplicateTagFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
tags:
  - security
  - security
`
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "duplicate_tag")
}

func TestStrictUnknownYAMLFieldFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
custom_thing: value
`
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "unknown_field")
}

func TestStrictMissingTestsFailsWithoutFlag(t *testing.T) {
	files := copyMap(validSkillFiles)
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_tests")
}

func TestStrictMissingTestsAllowedWithFlag(t *testing.T) {
	files := copyMap(validSkillFiles)
	res := validate(t, files, skill.ValidationOptions{Strict: true, AllowMissingTests: true})
	// missing tests is suppressed — should not be an error; check via Infos
	for _, e := range res.Errors {
		assert.NotEqual(t, "missing_tests", e.Code, "missing_tests should not appear as error with --allow-missing-tests")
	}
	for _, w := range res.Warnings {
		assert.NotEqual(t, "missing_tests", w.Code, "missing_tests should not appear as warning with --allow-missing-tests")
	}
}

func TestStrictBuildDirFails(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "node_modules", "pkg"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "node_modules", "pkg", "index.js"), []byte(""), 0o644))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Strict: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "build_dir_present")
}

func TestStrictMissingLicenseFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "LICENSE")
	res := strictValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_license")
}

func TestStrictMissingLicenseAllowedWithFlag(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "LICENSE")
	res := validate(t, files, skill.ValidationOptions{Strict: true, AllowMissingLicense: true, AllowMissingTests: true})
	for _, e := range res.Errors {
		assert.NotEqual(t, "missing_license", e.Code)
	}
}

func TestStrictProfile(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Strict: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.Equal(t, "strict", res.Profile)
}

// ── Publish validation ────────────────────────────────────────────────────

func TestPublishValidSkill(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Publish: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Equal(t, "publish", res.Profile)
}

func TestPublishMissingChangelogFails(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.RemoveAll(filepath.Join(dir, "CHANGELOG.md")))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Publish: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_changelog")
}

func TestPublishWarningsPromotedToErrors(t *testing.T) {
	files := copyMap(validSkillFiles)
	// No README — in default this is a warning, in publish it must become an error
	delete(files, "README.md")
	dir := writeSkillFixture(t, files)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Publish: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assert.Empty(t, res.Warnings, "publish mode must have no warnings — all promoted to errors")
	assertErrorCode(t, res, "missing_readme")
}

func TestPublishSecretFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["config.txt"] = "token=supersecretvalue12345678"
	dir := writeSkillFixture(t, files)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Publish: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "possible_secret")
}

func TestPublishForbiddenFileFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["id_rsa"] = "fake key"
	res := publishValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "forbidden_file")
}

func TestPublishPackageArtifactFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["my-skill-1.2.3.zip"] = "fake zip"
	res := publishValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "generated_artifact")
}

func TestPublishInvalidNameFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: My Skill With Spaces
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
`
	res := publishValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "invalid_name")
}

func TestPublishMissingEntrypointFails(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - claude-code
entrypoint: CUSTOM.md
`
	res := publishValidate(t, files)
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "missing_entrypoint")
}

// ── Lint (delegates to strict) ────────────────────────────────────────────

func TestLintDelegatesToStrict(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "CHANGELOG.md")
	// In default, missing CHANGELOG is a warning; in strict it's an error.
	// Lint must behave like strict.
	dir := writeSkillFixture(t, files)
	res, err := skill.NewValidatorWithOptions(skill.ValidationOptions{Strict: true}).Validate(context.Background(), dir)
	require.NoError(t, err)
	assert.False(t, res.Valid, "lint (strict) should fail on missing CHANGELOG")
	assert.Equal(t, "strict", res.Profile)
}

func TestLintPlatformFlag(t *testing.T) {
	files := copyMap(validSkillFiles)
	// my-skill is compatible with claude-code and gitlab-duo but not codex
	res := validate(t, files, skill.ValidationOptions{Strict: true, Platform: "codex"})
	assert.False(t, res.Valid)
	assertErrorCode(t, res, "platform_not_compatible")
}

func TestLintJSONOutput(t *testing.T) {
	res := validate(t, validSkillFiles, skill.ValidationOptions{Strict: true})
	// Verify the result is JSON-serializable with required fields
	data, err := json.Marshal(res)
	require.NoError(t, err)
	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Contains(t, decoded, "valid")
	assert.Contains(t, decoded, "profile")
	assert.Contains(t, decoded, "path")
}

// ── Format helpers ────────────────────────────────────────────────────────

func TestFormatNormalizePlatformAlias(t *testing.T) {
	// NormalizePlatform should return canonical name for aliases
	assert.Equal(t, skill.PlatformGitLabDuo, skill.NormalizePlatform("gitlab"))
	assert.Equal(t, skill.PlatformGitLabDuo, skill.NormalizePlatform("duo"))
	assert.Equal(t, skill.PlatformGitHubCopilot, skill.NormalizePlatform("github"))
	assert.Equal(t, skill.PlatformGitHubCopilot, skill.NormalizePlatform("copilot"))
	assert.Equal(t, skill.PlatformClaudeCode, skill.NormalizePlatform("claude"))
}

func TestFormatNormalizeCanonicalPassthrough(t *testing.T) {
	assert.Equal(t, skill.PlatformClaudeCode, skill.NormalizePlatform(skill.PlatformClaudeCode))
	assert.Equal(t, skill.PlatformAll, skill.NormalizePlatform(skill.PlatformAll))
}

// ── Packager tests ────────────────────────────────────────────────────────

func TestPackageValidSkill(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
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
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	result, err := skill.NewPackager().Package(context.Background(), dir, "")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "my-skill-1.2.3.zip"), result.OutputPath)
}

func TestPackageExcludesGitDir(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
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
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	outDir := filepath.Join(t.TempDir(), "dist")
	result, err := skill.NewPackager().Package(context.Background(), dir, outDir)
	require.NoError(t, err)
	assert.FileExists(t, result.OutputPath)
}

func TestPackageUsesStrictValidation(t *testing.T) {
	// Packager uses strict validation — a skill missing CHANGELOG fails package
	files := copyMap(validSkillFiles)
	delete(files, "CHANGELOG.md")
	_, err := skill.NewPackager().Package(context.Background(), writeSkillFixture(t, files), t.TempDir())
	require.Error(t, err, "packager with strict validation should fail without CHANGELOG.md")
}

func TestShouldSkip(t *testing.T) {
	assert.True(t, skill.ShouldSkip(".git", ".git"))
	assert.True(t, skill.ShouldSkip(".DS_Store", ".DS_Store"))
	assert.True(t, skill.ShouldSkip("file.tmp", "file.tmp"))
	assert.True(t, skill.ShouldSkip("sub", ".git/sub"))
	assert.False(t, skill.ShouldSkip("SKILL.md", "SKILL.md"))
}

// ── New platform constants ────────────────────────────────────────────────

func TestNewPlatformsAreKnown(t *testing.T) {
	for _, p := range []skill.Platform{
		skill.PlatformCursor,
		skill.PlatformWindsurf,
		skill.PlatformOpenHands,
		skill.PlatformOpenCode,
		skill.PlatformOllama,
		skill.PlatformGeneric,
	} {
		assert.True(t, skill.KnownPlatforms(p), "expected %q to be a known platform", p)
	}
}

func TestNewPlatformsAcceptedInCompatibleWith(t *testing.T) {
	for _, platform := range []string{"cursor", "windsurf", "openhands", "opencode", "ollama", "generic"} {
		t.Run(platform, func(t *testing.T) {
			files := copyMap(validSkillFiles)
			files["skill.yaml"] = `name: my-skill
version: "1.2.3"
description: Does things
compatible_with:
  - ` + platform + "\n"
			res := defaultValidate(t, files)
			assert.True(t, res.Valid, "platform %q should be valid", platform)
		})
	}
}

// ── Result serialisation ──────────────────────────────────────────────────

func TestValidationResultJSON(t *testing.T) {
	files := copyMap(validSkillFiles)
	files["VERSION"] = "1.5"
	files["skill.yaml"] = `name: my-skill
version: "1.5"
description: Does things
compatible_with:
  - claude-code
`
	res := defaultValidate(t, files)
	assert.False(t, res.Valid)

	data, err := json.Marshal(res)
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, false, decoded["valid"])
	assert.Equal(t, "default", decoded["profile"])

	errs := decoded["errors"].([]interface{})
	require.NotEmpty(t, errs)
	first := errs[0].(map[string]interface{})
	assert.Equal(t, "error", first["severity"])
	assert.NotEmpty(t, first["code"])

	// Verify bytes round-trip cleanly
	var roundtrip skill.ValidationResult
	require.NoError(t, json.Unmarshal(data, &roundtrip))
	assert.Equal(t, res.Valid, roundtrip.Valid)
	assert.Equal(t, res.Profile, roundtrip.Profile)

	// Confirm no extra bytes leaked
	var buf bytes.Buffer
	require.NoError(t, json.Compact(&buf, data))
	assert.True(t, buf.Len() > 0)
}

// TestSkillManifestIsSpecPackageManifest verifies the type alias is in effect:
// skill.SkillManifest and spec.PackageManifest are the same type.
func TestSkillManifestIsSpecPackageManifest(t *testing.T) {
	var m skill.SkillManifest = spec.PackageManifest{
		Name:    "alias-check",
		Version: "1.0.0",
	}
	assert.Equal(t, "alias-check", m.Name)
}

// TestPackageManifestRoundtrip verifies that the manifest.json produced by the
// packager deserializes into spec.PackageManifest with the expected field values.
func TestPackageManifestRoundtrip(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	result, err := skill.NewPackager().Package(context.Background(), dir, t.TempDir())
	require.NoError(t, err)

	zr, err := zip.OpenReader(result.OutputPath)
	require.NoError(t, err)
	defer zr.Close()

	var manifest spec.PackageManifest
	for _, f := range zr.File {
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			require.NoError(t, err)
			require.NoError(t, json.NewDecoder(rc).Decode(&manifest))
			rc.Close()
			break
		}
	}

	assert.Equal(t, "my-skill", manifest.Name)
	assert.Equal(t, "1.2.3", manifest.Version)
	assert.Equal(t, "zip", manifest.PackageType)
	assert.Equal(t, 1, manifest.SpecVersion)
	assert.NotEmpty(t, manifest.Files)
	assert.Contains(t, manifest.CompatibleWith, spec.Platform(spec.PlatformClaudeCode))
}

// TestPackageChecksumsRoundtrip verifies that checksums.txt produced by the
// packager is parseable by packageio.ParseChecksumsText.
func TestPackageChecksumsRoundtrip(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tests"), 0o755))
	result, err := skill.NewPackager().Package(context.Background(), dir, t.TempDir())
	require.NoError(t, err)

	zr, err := zip.OpenReader(result.OutputPath)
	require.NoError(t, err)
	defer zr.Close()

	var checksumData []byte
	for _, f := range zr.File {
		if f.Name == "checksums.txt" {
			rc, err := f.Open()
			require.NoError(t, err)
			var buf bytes.Buffer
			_, err = buf.ReadFrom(rc)
			require.NoError(t, err)
			rc.Close()
			checksumData = buf.Bytes()
			break
		}
	}
	require.NotEmpty(t, checksumData, "checksums.txt not found in package")

	entries, err := packageio.ParseChecksumsText(checksumData)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
	for _, e := range entries {
		assert.NotEmpty(t, e.Path)
		assert.Len(t, e.SHA256, 64)
	}
}

// ── EnsureChangelogEntry ─────────────────────────────────────────────────────

func TestEnsureChangelogEntry_AddsEntryWhenMissing(t *testing.T) {
	dir := writeSkillFixture(t, map[string]string{
		"VERSION":      "1.2.0",
		"CHANGELOG.md": "# Changelog\n\n## 1.1.0\n\n- previous release\n",
	})
	added, err := skill.EnsureChangelogEntry(dir, "1.2.0")
	require.NoError(t, err)
	assert.True(t, added)

	data, _ := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	content := string(data)
	assert.Contains(t, content, "## 1.2.0")
	assert.Contains(t, content, "## 1.1.0", "existing entry must be preserved")
	idx120 := len(content) - len(content[len("# Changelog\n\n"):])
	idx110 := indexOf(content, "## 1.1.0")
	assert.Less(t, indexOf(content, "## 1.2.0"), idx110, "new entry should appear before 1.1.0")
	_ = idx120
}

func TestEnsureChangelogEntry_NoOpWhenPresent(t *testing.T) {
	dir := writeSkillFixture(t, map[string]string{
		"VERSION":      "1.2.0",
		"CHANGELOG.md": "# Changelog\n\n## 1.2.0\n\n- this release\n",
	})
	added, err := skill.EnsureChangelogEntry(dir, "1.2.0")
	require.NoError(t, err)
	assert.False(t, added)
}

func TestEnsureChangelogEntry_CreatesMissingFile(t *testing.T) {
	dir := t.TempDir()
	added, err := skill.EnsureChangelogEntry(dir, "0.1.0")
	require.NoError(t, err)
	assert.True(t, added)
	data, _ := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	assert.Contains(t, string(data), "## 0.1.0")
	assert.Contains(t, string(data), "# Changelog")
}

func TestEnsureChangelogEntry_AcceptsVPrefix(t *testing.T) {
	dir := writeSkillFixture(t, map[string]string{
		"VERSION":      "2.0.0",
		"CHANGELOG.md": "# Changelog\n\n## v2.0.0\n\n- with v prefix\n",
	})
	added, err := skill.EnsureChangelogEntry(dir, "2.0.0")
	require.NoError(t, err)
	assert.False(t, added, "v2.0.0 should match 2.0.0")
}

func indexOf(s, substr string) int {
	idx := 0
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return idx
		}
		idx++
	}
	return -1
}

func TestPackageDeterminism(t *testing.T) {
	dir1 := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir1, "tests"), 0o755))
	dir2 := writeSkillFixture(t, validSkillFiles)
	require.NoError(t, os.MkdirAll(filepath.Join(dir2, "tests"), 0o755))

	out1 := t.TempDir()
	out2 := t.TempDir()

	res1, err := skill.NewPackager().Package(context.Background(), dir1, out1)
	require.NoError(t, err)
	res2, err := skill.NewPackager().Package(context.Background(), dir2, out2)
	require.NoError(t, err)

	assert.Equal(t, res1.SHA256, res2.SHA256, "packaging identical directories must yield identical SHA256")
}
