package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage skpm configuration",
		Long: `Inspect and validate the skpm configuration file.

  skpm config validate         — validate ~/.config/skpm/config.yaml
  skpm config show             — print the resolved config (tokens masked)
  skpm config get <key>        — print a single top-level config value
  skpm config set <key> value  — set a single top-level config value`,
	}
	cmd.AddCommand(newConfigValidateCmd())
	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigExportCmd())
	cmd.AddCommand(newConfigImportCmd())
	return cmd
}

var configScalarKeys = []string{"default_registry", "cache_dir", "log_level", "concurrency"}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "get <key>",
		Short:     "Print a top-level config value",
		ValidArgs: configScalarKeys,
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := resolveConfigPath("")
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}
			cfg, err := config.LoadFrom(path)
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			val := configGetScalar(cfg, args[0])
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "config get", Data: map[string]string{args[0]: val}})
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), val)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "set <key> <value>",
		Short:     "Set a top-level config value",
		ValidArgs: configScalarKeys,
		Args:      cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			for _, k := range configScalarKeys {
				if k == key {
					goto valid
				}
			}
			return &UserError{Message: fmt.Sprintf("unknown config key %q; valid keys: %s", key, stringsJoin(configScalarKeys, ", "))}
		valid:
			path, err := resolveConfigPath("")
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}
			cfg, err := config.LoadFrom(path)
			if err != nil && !os.IsNotExist(err) {
				return &InternalError{Message: "load config", Cause: err}
			}
			if cfg == nil {
				cfg = &config.Config{Registries: map[string]config.RegistryConfig{}}
			}
			if err := configSetScalar(cfg, key, value); err != nil {
				return err
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would set %s = %s in %s\n", key, value, path)
				return nil
			}
			if err := writeConfig(path, cfg); err != nil {
				return &InternalError{Message: "write config", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Set %s = %s\n", key, value)
			return nil
		},
	}
}

func configGetScalar(cfg *config.Config, key string) string {
	switch key {
	case "default_registry":
		return cfg.DefaultRegistry
	case "cache_dir":
		return cfg.CacheDir
	case "log_level":
		return cfg.LogLevel
	case "concurrency":
		return fmt.Sprintf("%d", cfg.Concurrency)
	default:
		return ""
	}
}

func configSetScalar(cfg *config.Config, key, value string) error {
	switch key {
	case "default_registry":
		cfg.DefaultRegistry = value
	case "cache_dir":
		cfg.CacheDir = value
	case "log_level":
		cfg.LogLevel = value
	case "concurrency":
		n := 0
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 1 {
			return &UserError{Message: fmt.Sprintf("concurrency must be a positive integer, got %q", value)}
		}
		cfg.Concurrency = n
	default:
		return &UserError{Message: fmt.Sprintf("unknown config key %q", key)}
	}
	return nil
}

// ── skpm config validate ──────────────────────────────────────────────────

type configError struct {
	Field   string
	Message string
}

func newConfigValidateCmd() *cobra.Command {
	var cfgPath string
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the skpm config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			path, err := resolveConfigPath(cfgPath)
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf(
					"config not found: %s\nRun 'skpm init config' to create one.", path,
				)}
			}

			cfg, loadErr := config.LoadFrom(path)
			errs := validateConfig(cfg, path, loadErr)

			if format == OutputJSON {
				errStrs := make([]string, len(errs))
				for i, e := range errs {
					errStrs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
				}
				PrintResult(format, CommandResult{
					Success: len(errs) == 0,
					Command: "config validate",
					Data:    map[string]interface{}{"path": path},
					Errors:  errStrs,
				})
				if len(errs) > 0 {
					os.Exit(1)
				}
				return nil
			}

			if len(errs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "✓ %s is valid\n", path)
				return nil
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "✗ %s has %d error(s)\n\n", path, len(errs))
			for _, e := range errs {
				fmt.Fprintf(cmd.ErrOrStderr(), "  [%s] %s\n", e.Field, e.Message)
			}
			os.Exit(1)
			return nil
		},
	}
	cmd.Flags().StringVar(&cfgPath, "path", "", "Path to config file (default: ~/.config/skpm/config.yaml)")
	return cmd
}

func validateConfig(cfg *config.Config, path string, loadErr error) []configError {
	var errs []configError

	if loadErr != nil {
		return []configError{{Field: "file", Message: fmt.Sprintf("parse error: %v", loadErr)}}
	}

	if cfg.DefaultRegistry != "" {
		if _, ok := cfg.Registries[cfg.DefaultRegistry]; !ok {
			errs = append(errs, configError{
				Field:   "default_registry",
				Message: fmt.Sprintf("%q is not defined in registries", cfg.DefaultRegistry),
			})
		}
	}

	for name, rc := range cfg.Registries {
		errs = append(errs, validateRegistryConfig(name, rc)...)
	}

	return errs
}

func validateRegistryConfig(name string, rc config.RegistryConfig) []configError {
	var errs []configError
	field := func(f string) string { return fmt.Sprintf("registries.%s.%s", name, f) }

	validTypes := map[string]bool{"github": true, "gitlab": true, "artifactory": true, "local": true, "skillforge": true, "generic-http": true}
	if !validTypes[rc.Type] {
		errs = append(errs, configError{
			Field:   field("type"),
			Message: fmt.Sprintf("%q is not a valid type (github, gitlab, artifactory, local, skillforge, generic-http)", rc.Type),
		})
		return errs
	}

	if rc.URL == "" && rc.Path == "" && rc.Repo == "" {
		errs = append(errs, configError{Field: field("url"), Message: "url, path, or repo is required"})
	}

	switch rc.Type {
	case "gitlab":
		if rc.Project == "" {
			errs = append(errs, configError{
				Field:   field("project"),
				Message: "project is required for gitlab registries (format: namespace/project)",
			})
		} else if !strings.Contains(rc.Project, "/") {
			errs = append(errs, configError{
				Field:   field("project"),
				Message: fmt.Sprintf("%q must be in namespace/project format", rc.Project),
			})
		}
	case "github":
		repo := rc.Repo
		if repo == "" {
			repo = rc.URL
		}
		if repo != "" && strings.HasPrefix(repo, "https://") {
			fieldName := "repo"
			if rc.Repo == "" && rc.URL != "" {
				fieldName = "url"
			}
			errs = append(errs, configError{
				Field:   field(fieldName),
				Message: "github repo must be owner/repo (e.g. myorg/agent-skills), not a full URL",
			})
		} else if repo != "" && !strings.Contains(repo, "/") {
			fieldName := "repo"
			if rc.Repo == "" && rc.URL != "" {
				fieldName = "url"
			}
			errs = append(errs, configError{
				Field:   field(fieldName),
				Message: fmt.Sprintf("%q must be in owner/repo format", repo),
			})
		}
	case "artifactory":
		if rc.Repo == "" && rc.URL != "" && strings.Contains(rc.URL, "#") {
			return errs
		}
		if rc.Repo == "" {
			fieldName := "repo"
			if rc.URL != "" {
				fieldName = "url"
			}
			errs = append(errs, configError{
				Field:   field(fieldName),
				Message: "repo is required for artifactory registries (or use legacy url format <base-url>#<repo>)",
			})
		}
	case "generic-http":
		if rc.Endpoints["resolve"] == "" && rc.Endpoints["download"] == "" {
			errs = append(errs, configError{
				Field:   field("endpoints"),
				Message: "generic-http should define at least endpoints.resolve or endpoints.download",
			})
		}
	}

	return errs
}

// ── skpm config show ──────────────────────────────────────────────────────

func newConfigShowCmd() *cobra.Command {
	var cfgPath string
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Print the resolved config (tokens masked)",
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			path, err := resolveConfigPath(cfgPath)
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf(
					"config not found: %s\nRun 'skpm init config' to create one.", path,
				)}
			}

			cfg, err := config.LoadFrom(path)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("parse config: %v", err)}
			}

			masked := maskTokens(cfg)

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "config show",
					Data:    masked,
				})
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "# Resolved config from %s\n\n", path)
			data, _ := yaml.Marshal(masked)
			fmt.Fprint(cmd.OutOrStdout(), string(data))
			return nil
		},
	}
	cmd.Flags().StringVar(&cfgPath, "path", "", "Path to config file (default: ~/.config/skpm/config.yaml)")
	return cmd
}

// maskTokens returns a copy of the config with all tokens replaced by "***".
func maskTokens(cfg *config.Config) map[string]interface{} {
	registries := make(map[string]interface{}, len(cfg.Registries))
	for name, rc := range cfg.Registries {
		token := ""
		if rc.Token != "" {
			token = "***"
		}
		auth := rc.Auth
		if auth.Token != "" {
			auth.Token = "***"
		}
		if auth.Password != "" {
			auth.Password = "***"
		}
		entry := map[string]interface{}{
			"type":      rc.Type,
			"url":       rc.URL,
			"repo":      rc.Repo,
			"path":      rc.Path,
			"namespace": rc.Namespace,
			"auth":      auth,
		}
		if rc.Project != "" {
			entry["project"] = rc.Project
		}
		if len(rc.Headers) > 0 {
			entry["headers"] = rc.Headers
		}
		if len(rc.Endpoints) > 0 {
			entry["endpoints"] = rc.Endpoints
		}
		if len(rc.Capabilities) > 0 {
			entry["capabilities"] = rc.Capabilities
		}
		entry["token"] = token
		registries[name] = entry
	}
	return map[string]interface{}{
		"default_registry": cfg.DefaultRegistry,
		"cache_dir":        cfg.CacheDir,
		"log_level":        cfg.LogLevel,
		"concurrency":      cfg.Concurrency,
		"registries":       registries,
	}
}

func resolveConfigPath(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	return configFilePath()
}
