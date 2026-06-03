package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDefaults(t *testing.T) {
	cfg := defaults()
	assert.Equal(t, defaultLogLevel, cfg.LogLevel)
	assert.Equal(t, defaultConcurrency, cfg.Concurrency)
	assert.NotEmpty(t, cfg.CacheDir)
}

func TestLoadMissingFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, defaultLogLevel, cfg.LogLevel)
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfgDir := filepath.Join(dir, "sctl")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	data, _ := yaml.Marshal(map[string]any{
		"log_level":   "debug",
		"concurrency": 8,
		"cache_dir":   "/tmp/sctl-cache",
	})
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0o644))

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 8, cfg.Concurrency)
	assert.Equal(t, "/tmp/sctl-cache", cfg.CacheDir)
}

func TestEnvOverrides(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SCTL_CACHE_DIR", "/env/cache")
	t.Setenv("SCTL_LOG_LEVEL", "warn")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "/env/cache", cfg.CacheDir)
	assert.Equal(t, "warn", cfg.LogLevel)
}

func TestEnvRegistryTokenOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("SCTL_REGISTRY_TOKEN", "secret-token")

	cfgDir := filepath.Join(dir, "sctl")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	data, _ := yaml.Marshal(map[string]any{
		"default_registry": "mygh",
		"registries": map[string]any{
			"mygh": map[string]any{"type": "github", "url": "https://github.com"},
		},
	})
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0o644))

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "secret-token", cfg.Registries["mygh"].Token)
}
