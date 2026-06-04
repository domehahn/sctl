package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/sctl/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage skpm configuration",
		Long: `Inspect and validate the skpm configuration file.

  skpm config validate    — validate ~/.config/skpm/config.yaml
  skpm config show        — print the resolved config (tokens masked)`,
	}
	cmd.AddCommand(newConfigValidateCmd())
	cmd.AddCommand(newConfigShowCmd())
	return cmd
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

	validTypes := map[string]bool{"github": true, "gitlab": true, "artifactory": true, "local": true}
	if !validTypes[rc.Type] {
		errs = append(errs, configError{
			Field:   field("type"),
			Message: fmt.Sprintf("%q is not a valid type (github, gitlab, artifactory, local)", rc.Type),
		})
		return errs
	}

	if rc.URL == "" {
		errs = append(errs, configError{Field: field("url"), Message: "url is required"})
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
		if rc.URL != "" && strings.HasPrefix(rc.URL, "https://") {
			errs = append(errs, configError{
				Field:   field("url"),
				Message: "github url must be owner/repo (e.g. myorg/agent-skills), not a full URL",
			})
		} else if rc.URL != "" && !strings.Contains(rc.URL, "/") {
			errs = append(errs, configError{
				Field:   field("url"),
				Message: fmt.Sprintf("%q must be in owner/repo format", rc.URL),
			})
		}
	case "artifactory":
		if rc.URL != "" && !strings.Contains(rc.URL, "#") {
			errs = append(errs, configError{
				Field:   field("url"),
				Message: "artifactory url must be in format <base-url>#<repo-name> (e.g. https://artifactory.company.com/artifactory#agent-skills)",
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
		entry := map[string]interface{}{
			"type": rc.Type,
			"url":  rc.URL,
		}
		if rc.Project != "" {
			entry["project"] = rc.Project
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
