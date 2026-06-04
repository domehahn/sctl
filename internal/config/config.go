package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultConcurrency = 4
	defaultLogLevel    = "info"
)

type RegistryConfig struct {
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
	// Project is the GitLab namespace/project path (e.g. "platform/agent-skills").
	// Required for gitlab registries when using skpm add.
	// Not needed for skpm install — the lockfile source_url is used directly.
	Project string `yaml:"project,omitempty"`
}

type Config struct {
	DefaultRegistry string                    `yaml:"default_registry"`
	Registries      map[string]RegistryConfig `yaml:"registries"`
	CacheDir        string                    `yaml:"cache_dir"`
	LogLevel        string                    `yaml:"log_level"`
	Concurrency     int                       `yaml:"concurrency"`
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}
	return LoadFrom(path)
}

func LoadFrom(path string) (*Config, error) {
	cfg := defaults()

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config %s: %w", path, err)
		}
	}

	applyEnv(cfg)
	return cfg, nil
}

func defaults() *Config {
	cacheDir, _ := os.UserCacheDir()
	return &Config{
		CacheDir:    filepath.Join(cacheDir, "skpm"),
		LogLevel:    defaultLogLevel,
		Concurrency: defaultConcurrency,
		Registries:  make(map[string]RegistryConfig),
	}
}

func configPath() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "skpm", "config.yaml"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "skpm", "config.yaml"), nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("SKPM_CACHE_DIR"); v != "" {
		cfg.CacheDir = v
	}
	if v := os.Getenv("SKPM_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("SKPM_REGISTRY_TOKEN"); v != "" {
		if cfg.DefaultRegistry != "" {
			if reg, ok := cfg.Registries[cfg.DefaultRegistry]; ok {
				reg.Token = v
				cfg.Registries[cfg.DefaultRegistry] = reg
			}
		}
	}
}
