package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configExportData is the portable representation written by `skpm config export`.
type configExportData struct {
	DefaultRegistry string                      `yaml:"default_registry,omitempty"`
	CacheDir        string                      `yaml:"cache_dir,omitempty"`
	LogLevel        string                      `yaml:"log_level,omitempty"`
	Concurrency     int                         `yaml:"concurrency,omitempty"`
	Registries      map[string]exportedRegistry `yaml:"registries,omitempty"`
}

type exportedRegistry struct {
	URL                string `yaml:"url,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify,omitempty"`
	// Token is omitted by default; included when --show-tokens is set.
	Token string `yaml:"token,omitempty"`
	// TokenEnv is always exported so the recipient knows which env var to set.
	TokenEnv string `yaml:"token_env,omitempty"`
}

func newConfigExportCmd() *cobra.Command {
	var outPath string
	var showTokens bool

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export skpm config to a portable YAML file",
		Long: `Dumps the current skpm configuration to a YAML file suitable for sharing
or committing to a repository (tokens masked by default).

Use --show-tokens to include actual token values — treat the output as a secret.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := resolveConfigPath("")
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}
			cfg, err := config.LoadFrom(cfgPath)
			if err != nil && !os.IsNotExist(err) {
				return &InternalError{Message: "load config", Cause: err}
			}
			if cfg == nil {
				cfg = &config.Config{Registries: map[string]config.RegistryConfig{}}
			}

			out := configExportData{
				DefaultRegistry: cfg.DefaultRegistry,
				CacheDir:        cfg.CacheDir,
				LogLevel:        cfg.LogLevel,
				Concurrency:     cfg.Concurrency,
				Registries:      map[string]exportedRegistry{},
			}
			for name, rc := range cfg.Registries {
				er := exportedRegistry{
					URL:                rc.URL,
					InsecureSkipVerify: rc.TLS.InsecureSkipVerify,
					TokenEnv:           rc.Auth.TokenEnv,
				}
				if showTokens {
					er.Token = rc.Auth.Token
					if er.Token == "" {
						er.Token = rc.Token
					}
				}
				out.Registries[name] = er
			}

			data, err := yaml.Marshal(out)
			if err != nil {
				return &InternalError{Message: "marshal config", Cause: err}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would write %d byte(s) to %s\n", len(data), outPath)
				return nil
			}

			dest := outPath
			if dest == "" {
				fmt.Fprint(cmd.OutOrStdout(), string(data))
				return nil
			}

			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return &InternalError{Message: "create output directory", Cause: err}
			}
			if err := os.WriteFile(dest, data, 0o600); err != nil {
				return &InternalError{Message: "write export file", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Config exported to %s\n", dest)
			return nil
		},
	}

	cmd.Flags().StringVar(&outPath, "out", "", "Output file path (stdout if omitted)")
	cmd.Flags().BoolVar(&showTokens, "show-tokens", false, "Include token values in the export (sensitive)")
	return cmd
}

func newConfigImportCmd() *cobra.Command {
	var inPath string
	var merge bool

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import skpm config from a portable YAML file",
		Long: `Restores skpm configuration from a file produced by 'skpm config export'.

By default, the import replaces the existing config. Use --merge to layer the
imported registries on top of the existing config without removing existing entries.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if inPath == "" {
				return &UserError{Message: "provide --in <file>"}
			}

			raw, err := os.ReadFile(inPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("cannot read import file: %v", err)}
			}

			var imported configExportData
			if err := yaml.Unmarshal(raw, &imported); err != nil {
				return &UserError{Message: fmt.Sprintf("invalid import file: %v", err)}
			}

			cfgPath, err := resolveConfigPath("")
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}

			var cfg *config.Config
			if merge {
				cfg, _ = config.LoadFrom(cfgPath)
			}
			if cfg == nil {
				cfg = &config.Config{Registries: map[string]config.RegistryConfig{}}
			}

			if imported.DefaultRegistry != "" {
				cfg.DefaultRegistry = imported.DefaultRegistry
			}
			if imported.CacheDir != "" {
				cfg.CacheDir = imported.CacheDir
			}
			if imported.LogLevel != "" {
				cfg.LogLevel = imported.LogLevel
			}
			if imported.Concurrency > 0 {
				cfg.Concurrency = imported.Concurrency
			}
			for name, er := range imported.Registries {
				rc := cfg.Registries[name]
				if er.URL != "" {
					rc.URL = er.URL
				}
				rc.TLS.InsecureSkipVerify = er.InsecureSkipVerify
				if er.TokenEnv != "" {
					rc.Auth.TokenEnv = er.TokenEnv
				}
				if er.Token != "" {
					rc.Auth.Token = er.Token
					rc.Token = er.Token
				}
				cfg.Registries[name] = rc
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would import %d registry/registries into %s\n",
					len(imported.Registries), cfgPath)
				return nil
			}

			if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
				return &InternalError{Message: "create config directory", Cause: err}
			}
			if err := writeConfig(cfgPath, cfg); err != nil {
				return &InternalError{Message: "write config", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Config imported from %s → %s (%d registry/registries)\n",
				inPath, cfgPath, len(imported.Registries))
			return nil
		},
	}

	cmd.Flags().StringVar(&inPath, "in", "", "Input file path (required)")
	cmd.Flags().BoolVar(&merge, "merge", false, "Merge with existing config instead of replacing")
	return cmd
}
