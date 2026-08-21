package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	defaultConcurrency = 4
	defaultLogLevel    = "info"
)

type RegistryConfig struct {
	Name      string            `yaml:"name,omitempty"`
	Type      string            `yaml:"type"`
	URL       string            `yaml:"url,omitempty"`
	Repo      string            `yaml:"repo,omitempty"`
	Project   string            `yaml:"project,omitempty"`
	Path      string            `yaml:"path,omitempty"`
	Namespace string            `yaml:"namespace,omitempty"`
	Token     string            `yaml:"token,omitempty"`
	Auth      AuthConfig        `yaml:"auth,omitempty"`
	Headers   map[string]string `yaml:"headers,omitempty"`
	Endpoints map[string]string `yaml:"endpoints,omitempty"`
	// Capabilities can be used when a registry cannot expose a capabilities endpoint.
	Capabilities map[string]bool   `yaml:"capabilities,omitempty"`
	TLS          TLSConfig         `yaml:"tls,omitempty"`
	Metadata     map[string]string `yaml:"metadata,omitempty"`
	// Project is the GitLab namespace/project path (e.g. "platform/agent-skills").
	// Required for gitlab registries when using skpm add.
	// Not needed for skpm install — the lockfile source_url is used directly.
}

type AuthConfig struct {
	Type        string `yaml:"type,omitempty"` // none, bearer, basic, header
	Token       string `yaml:"token,omitempty"`
	TokenEnv    string `yaml:"token_env,omitempty"`
	Username    string `yaml:"username,omitempty"`
	Password    string `yaml:"password,omitempty"`
	PasswordEnv string `yaml:"password_env,omitempty"`
	HeaderName  string `yaml:"header_name,omitempty"`
}

type TLSConfig struct {
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify,omitempty"`
	CAFile             string `yaml:"ca_file,omitempty"`
}

type Config struct {
	DefaultRegistry string                    `yaml:"default_registry"`
	Registries      map[string]RegistryConfig `yaml:"registries"`
	CacheDir        string                    `yaml:"cache_dir"`
	LogLevel        string                    `yaml:"log_level"`
	Concurrency     int                       `yaml:"concurrency"`
	// TrustedSigners maps an attestation signature's key_id (e.g. from a
	// `skil key generate` output) to its base64-encoded Ed25519 public
	// key. Used by `skpm attestations --verify` (see
	// internal/attestation.Verify) to independently check a stored
	// attestation's signature without depending on skil as a library —
	// mirrors skil's own policy.TrustedSigners naming and shape.
	TrustedSigners map[string]string `yaml:"trusted_signers,omitempty"`
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

	normalizeRegistries(cfg)
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
				reg.Auth.Token = v
				cfg.Registries[cfg.DefaultRegistry] = reg
			}
		}
	}
	for name, reg := range cfg.Registries {
		if reg.Auth.TokenEnv != "" {
			if v := os.Getenv(reg.Auth.TokenEnv); v != "" {
				reg.Auth.Token = v
				reg.Token = v
			}
		}
		if reg.Auth.PasswordEnv != "" {
			if v := os.Getenv(reg.Auth.PasswordEnv); v != "" {
				reg.Auth.Password = v
			}
		}
		cfg.Registries[name] = reg
	}
}

func normalizeRegistries(cfg *Config) {
	for name, reg := range cfg.Registries {
		reg.Name = name
		if reg.Type == "" {
			reg.Type = "local"
		}
		if reg.Auth.Type == "" {
			if reg.Token != "" {
				reg.Auth.Type = "bearer"
				reg.Auth.Token = reg.Token
			} else {
				reg.Auth.Type = "none"
			}
		}
		if reg.Auth.Token == "" && reg.Token != "" {
			reg.Auth.Token = reg.Token
		}
		if reg.Token == "" && reg.Auth.Token != "" {
			reg.Token = reg.Auth.Token
		}
		if reg.Type == "github" && reg.Repo == "" {
			reg.Repo = reg.URL
		}
		if reg.Type == "local" && reg.Path == "" {
			reg.Path = reg.URL
		}
		if reg.Type == "artifactory" && reg.Repo == "" && reg.URL != "" {
			// Legacy artifactory URLs used <base-url>#<repo>.
			if hash := strings.LastIndex(reg.URL, "#"); hash >= 0 {
				reg.Repo = reg.URL[hash+1:]
				reg.URL = reg.URL[:hash]
			}
		}
		if reg.Namespace == "" {
			reg.Namespace = "default"
		}
		cfg.Registries[name] = reg
	}
}
