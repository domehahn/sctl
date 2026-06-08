package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── isLocalPath ───────────────────────────────────────────────────────────

func TestIsLocalPath(t *testing.T) {
	tests := []struct {
		arg  string
		want bool
	}{
		{"./my-skill", true},
		{"../my-skill", true},
		{"/abs/path", true},
		{".", true},
		{"..", true},
		{"my-skill", false},
		{"my-skill@1.0.0", false},
		{"github:org/repo", false},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.want, cli.IsLocalPath(tc.arg), "arg=%q", tc.arg)
	}
}

// ── parseSkillArg ─────────────────────────────────────────────────────────

func TestParseSkillArg(t *testing.T) {
	tests := []struct {
		arg         string
		wantName    string
		wantVersion string
	}{
		{"my-skill", "my-skill", ""},
		{"my-skill@1.0.0", "my-skill", "1.0.0"},
		{"my-skill@1.0.0-beta.1", "my-skill", "1.0.0-beta.1"},
	}
	for _, tc := range tests {
		name, version := cli.ParseSkillArg(tc.arg)
		assert.Equal(t, tc.wantName, name, "arg=%q", tc.arg)
		assert.Equal(t, tc.wantVersion, version, "arg=%q", tc.arg)
	}
}

// ── isValidSkillName ──────────────────────────────────────────────────────

func TestIsValidSkillName(t *testing.T) {
	assert.True(t, cli.IsValidSkillName("my-skill"))
	assert.True(t, cli.IsValidSkillName("skill123"))
	assert.True(t, cli.IsValidSkillName("a"))
	assert.False(t, cli.IsValidSkillName(""))
	assert.False(t, cli.IsValidSkillName("My-Skill"))
	assert.False(t, cli.IsValidSkillName("my_skill"))
	assert.False(t, cli.IsValidSkillName("my skill"))
}

// ── buildConfigYAML ───────────────────────────────────────────────────────

func TestBuildConfigYAMLGitLab(t *testing.T) {
	yaml := cli.BuildConfigYAML("mygl", "gitlab", "https://gitlab.com", "platform/skills", "secret", "~/.cache/skpm")
	assert.Contains(t, yaml, "default_registry: mygl")
	assert.Contains(t, yaml, "type: gitlab")
	assert.Contains(t, yaml, "project: platform/skills")
	assert.Contains(t, yaml, "token: secret")
	assert.Contains(t, yaml, "cache_dir: ~/.cache/skpm")
}

func TestBuildConfigYAMLNoToken(t *testing.T) {
	yaml := cli.BuildConfigYAML("mygh", "github", "myorg/repo", "", "", "~/.cache/skpm")
	assert.Contains(t, yaml, "token: \"\"")
	assert.NotContains(t, yaml, "project:")
}

// ── fileExists ────────────────────────────────────────────────────────────

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	assert.False(t, cli.FileExists(path))
	require.NoError(t, os.WriteFile(path, []byte("x"), 0o644))
	assert.True(t, cli.FileExists(path))
}

// ── writeAtomic ───────────────────────────────────────────────────────────

func TestWriteAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.yaml")
	require.NoError(t, cli.WriteAtomic(path, "content: true"))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "content: true", string(data))

	_, err = os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err))
}

// ── updateGitignore ───────────────────────────────────────────────────────

func TestUpdateGitignoreCreatesNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	updated, err := cli.UpdateGitignore(path)
	require.NoError(t, err)
	assert.True(t, updated)
	assert.FileExists(t, path)

	data, _ := os.ReadFile(path)
	assert.Contains(t, string(data), "skpm — installed skill directories")
}

func TestUpdateGitignoreAppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	require.NoError(t, os.WriteFile(path, []byte("node_modules/\n"), 0o644))

	updated, err := cli.UpdateGitignore(path)
	require.NoError(t, err)
	assert.True(t, updated)

	data, _ := os.ReadFile(path)
	assert.Contains(t, string(data), "node_modules/")
	assert.Contains(t, string(data), "skpm — installed skill directories")
}

func TestUpdateGitignoreIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	_, err := cli.UpdateGitignore(path)
	require.NoError(t, err)

	updated, err := cli.UpdateGitignore(path)
	require.NoError(t, err)
	assert.False(t, updated)
}

// ── validateConfig ────────────────────────────────────────────────────────

func TestValidateConfigValid(t *testing.T) {
	cfg := &config.Config{
		DefaultRegistry: "mygl",
		Registries: map[string]config.RegistryConfig{
			"mygl": {Type: "gitlab", URL: "https://gitlab.com", Project: "platform/skills"},
		},
	}
	errs := cli.ValidateConfig(cfg, "config.yaml", nil)
	assert.Empty(t, errs)
}

func TestValidateConfigLoadError(t *testing.T) {
	errs := cli.ValidateConfig(nil, "config.yaml", assert.AnError)
	require.Len(t, errs, 1)
	assert.Equal(t, "file", errs[0].Field)
}

func TestValidateConfigMissingDefaultRegistry(t *testing.T) {
	cfg := &config.Config{
		DefaultRegistry: "missing",
		Registries:      map[string]config.RegistryConfig{},
	}
	errs := cli.ValidateConfig(cfg, "config.yaml", nil)
	require.Len(t, errs, 1)
	assert.Equal(t, "default_registry", errs[0].Field)
}

func TestValidateRegistryConfigUnknownType(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "unknown", URL: "http://x"})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "type")
}

func TestValidateRegistryConfigGitLabMissingProject(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "gitlab", URL: "https://gitlab.com"})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "project")
}

func TestValidateRegistryConfigGitLabBadProjectFormat(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "gitlab", URL: "https://gitlab.com", Project: "noslash"})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "project")
}

func TestValidateRegistryConfigGitHubFullURL(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "github", URL: "https://github.com/org/repo"})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "url")
}

func TestValidateRegistryConfigGitHubNoSlash(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "github", URL: "noslash"})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "url")
}

func TestValidateRegistryConfigArtifactoryNoHash(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "artifactory", URL: "https://art.com/no-hash"})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "url")
}

func TestValidateRegistryConfigGitHubValid(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "github", URL: "org/repo"})
	assert.Empty(t, errs)
}

func TestValidateRegistryConfigLocalNoURL(t *testing.T) {
	errs := cli.ValidateRegistryConfig("r", config.RegistryConfig{Type: "local", URL: ""})
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Field, "url")
}

// ── maskTokens ────────────────────────────────────────────────────────────

func TestMaskTokens(t *testing.T) {
	cfg := &config.Config{
		DefaultRegistry: "mygl",
		CacheDir:        "~/.cache/skpm",
		LogLevel:        "info",
		Concurrency:     4,
		Registries: map[string]config.RegistryConfig{
			"mygl":     {Type: "gitlab", URL: "https://gitlab.com", Token: "secret"},
			"no-token": {Type: "github", URL: "org/repo", Token: ""},
		},
	}
	masked := cli.MaskTokens(cfg)

	regs := masked["registries"].(map[string]interface{})
	assert.Equal(t, "***", regs["mygl"].(map[string]interface{})["token"])
	assert.Equal(t, "", regs["no-token"].(map[string]interface{})["token"])
	assert.Equal(t, "mygl", masked["default_registry"])
}

// ── copyDir ───────────────────────────────────────────────────────────────

func TestCopyDir(t *testing.T) {
	src := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(src, "file.txt"), []byte("hello"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(src, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(src, "sub", "nested.txt"), []byte("nested"), 0o644))

	dst := filepath.Join(t.TempDir(), "dest")
	require.NoError(t, cli.CopyDir(src, dst))

	assert.FileExists(t, filepath.Join(dst, "file.txt"))
	assert.FileExists(t, filepath.Join(dst, "sub", "nested.txt"))
}

func TestCopyDirReplacesExisting(t *testing.T) {
	src := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(src, "new.txt"), []byte("new"), 0o644))

	dst := filepath.Join(t.TempDir(), "dest")
	require.NoError(t, os.MkdirAll(dst, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "old.txt"), []byte("old"), 0o644))

	require.NoError(t, cli.CopyDir(src, dst))

	assert.FileExists(t, filepath.Join(dst, "new.txt"))
	assert.NoFileExists(t, filepath.Join(dst, "old.txt"))
}

// ── configFilePath ────────────────────────────────────────────────────────

func TestConfigFilePath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	path, err := cli.ConfigFilePath()
	require.NoError(t, err)
	assert.Equal(t, "/custom/config/skpm/config.yaml", path)
}

func TestConfigFilePathDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	path, err := cli.ConfigFilePath()
	require.NoError(t, err)
	assert.Contains(t, path, "skpm/config.yaml")
}

// ── formatGitTag ─────────────────────────────────────────────────────────

func TestFormatGitTag(t *testing.T) {
	assert.Equal(t, "my-skill/v1.5.0", cli.FormatGitTag("my-skill", "1.5.0", "prefixed"))
	assert.Equal(t, "my-skill/v1.5.0", cli.FormatGitTag("my-skill", "1.5.0", ""))
	assert.Equal(t, "v1.5.0", cli.FormatGitTag("my-skill", "1.5.0", "plain"))
}
