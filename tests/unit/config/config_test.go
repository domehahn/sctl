package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 4, cfg.Concurrency)
	assert.NotEmpty(t, cfg.CacheDir)
}

func TestLoadMissingFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfgDir := filepath.Join(dir, "skpm")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	data, _ := yaml.Marshal(map[string]any{
		"log_level":   "debug",
		"concurrency": 8,
		"cache_dir":   "/tmp/skpm-cache",
	})
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0o644))

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 8, cfg.Concurrency)
	assert.Equal(t, "/tmp/skpm-cache", cfg.CacheDir)
}

func TestEnvOverrides(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SKPM_CACHE_DIR", "/env/cache")
	t.Setenv("SKPM_LOG_LEVEL", "warn")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "/env/cache", cfg.CacheDir)
	assert.Equal(t, "warn", cfg.LogLevel)
}

func TestEnvRegistryTokenOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("SKPM_REGISTRY_TOKEN", "secret-token")

	cfgDir := filepath.Join(dir, "skpm")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	data, _ := yaml.Marshal(map[string]any{
		"default_registry": "mygh",
		"registries": map[string]any{
			"mygh": map[string]any{"type": "github", "url": "https://github.com"},
		},
	})
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0o644))

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "secret-token", cfg.Registries["mygh"].Token)
}

func TestLoadFrom(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data, _ := yaml.Marshal(map[string]any{"log_level": "debug"})
	require.NoError(t, os.WriteFile(path, data, 0o644))

	cfg, err := config.LoadFrom(path)
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestLoadFromInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte("invalid: [yaml: {{"), 0o644))

	_, err := config.LoadFrom(path)
	assert.Error(t, err)
}
