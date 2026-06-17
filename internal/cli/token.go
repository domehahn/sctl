package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/spf13/cobra"
)

func newTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage registry auth tokens without interactive prompts",
		Long: `CI-friendly credential management. Unlike 'skpm login' (which prompts
interactively), 'skpm token add' reads the token from a flag or an
environment variable — suitable for scripts and CI pipelines.`,
	}
	cmd.AddCommand(newTokenListCmd())
	cmd.AddCommand(newTokenAddCmd())
	cmd.AddCommand(newTokenRevokeCmd())
	return cmd
}

func newTokenListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured registries and their token status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			if len(cfg.Registries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No registries configured.")
				return nil
			}

			names := make([]string, 0, len(cfg.Registries))
			for name := range cfg.Registries {
				names = append(names, name)
			}
			sort.Strings(names)

			type row struct {
				Registry string `json:"registry"`
				Type     string `json:"type"`
				HasToken bool   `json:"has_token"`
				TokenEnv string `json:"token_env,omitempty"`
				Default  bool   `json:"default,omitempty"`
			}
			var rows []row
			for _, name := range names {
				reg := cfg.Registries[name]
				rows = append(rows, row{
					Registry: name,
					Type:     reg.Type,
					HasToken: reg.Token != "",
					TokenEnv: reg.Auth.TokenEnv,
					Default:  name == cfg.DefaultRegistry,
				})
			}

			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "token list", Data: rows})
				return nil
			}

			for _, r := range rows {
				status := "no token"
				if r.HasToken {
					status = "token set"
				} else if r.TokenEnv != "" {
					status = "via env: " + r.TokenEnv
				}
				def := ""
				if r.Default {
					def = "  (default)"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %-28s  %-12s  %s%s\n", r.Registry, r.Type, status, def)
			}
			return nil
		},
	}
}

func newTokenAddCmd() *cobra.Command {
	var value string
	var fromEnv string
	var registryName string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Store a token for a registry without interactive prompts",
		Long: `Stores an authentication token for a registry without prompting.

Provide the token directly with --value, or name an environment variable
with --from-env (the variable is read at command time).

Examples:
  skpm token add --registry my-reg --value "$CI_TOKEN"
  skpm token add --registry my-reg --from-env CI_REGISTRY_TOKEN`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			reg := registryName
			if reg == "" {
				reg = cfg.DefaultRegistry
			}
			if reg == "" {
				return &UserError{Message: "no --registry specified and no default_registry configured"}
			}
			if _, ok := cfg.Registries[reg]; !ok {
				return &UserError{Message: fmt.Sprintf("registry %q not configured — run 'skpm registry add' first", reg)}
			}

			tok := value
			if tok == "" && fromEnv != "" {
				tok = os.Getenv(fromEnv)
				if tok == "" {
					return &UserError{Message: fmt.Sprintf("environment variable %s is empty or not set", fromEnv)}
				}
			}
			if tok == "" {
				return &UserError{Message: "provide --value <token> or --from-env <ENV_VAR_NAME>"}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would store token for registry %q\n", reg)
				return nil
			}

			if err := config.SetRegistryToken(reg, tok); err != nil {
				return &InternalError{Message: "save token", Cause: err}
			}
			path, _ := config.Path()
			fmt.Fprintf(cmd.OutOrStdout(), "Token stored for registry %q → %s\n", reg, path)
			return nil
		},
	}

	cmd.Flags().StringVar(&value, "value", "", "Token value")
	cmd.Flags().StringVar(&fromEnv, "from-env", "", "Read token from this environment variable name")
	cmd.Flags().StringVar(&registryName, "registry", "", "Registry name (default: default_registry from config)")
	return cmd
}

func newTokenRevokeCmd() *cobra.Command {
	var registryName string

	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Remove the stored token for a registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			reg := registryName
			if reg == "" {
				reg = cfg.DefaultRegistry
			}
			if reg == "" {
				return &UserError{Message: "no --registry specified and no default_registry configured"}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove token for registry %q\n", reg)
				return nil
			}

			if err := config.ClearRegistryToken(reg); err != nil {
				return &InternalError{Message: "clear token", Cause: err}
			}
			path, _ := config.Path()
			fmt.Fprintf(cmd.OutOrStdout(), "Token removed for registry %q → %s\n", reg, path)
			return nil
		},
	}

	cmd.Flags().StringVar(&registryName, "registry", "", "Registry name (default: default_registry from config)")
	return cmd
}
