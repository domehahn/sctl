package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newLoginCmd() *cobra.Command {
	var token string

	cmd := &cobra.Command{
		Use:   "login [registry]",
		Short: "Save an auth token for a registry",
		Long: `Saves an authentication token for a registry into the skpm config file.

If no registry name is given, the configured default_registry is used.
The token can be supplied via --token or will be prompted interactively.

Examples:
  skpm login                          # prompts for token, uses default_registry
  skpm login my-registry --token PAT  # stores token for "my-registry"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			name := registryArg(args, cfg)
			if name == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}

			tok := token
			if tok == "" {
				tok, err = promptToken(cmd, name)
				if err != nil {
					return err
				}
			}
			if strings.TrimSpace(tok) == "" {
				return &UserError{Message: "token must not be empty"}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would store token for registry %q\n", name)
				return nil
			}

			if err := config.SetRegistryToken(name, tok); err != nil {
				return &InternalError{Message: "save token", Cause: err}
			}

			path, _ := config.Path()
			fmt.Fprintf(cmd.OutOrStdout(), "Token saved for registry %q in %s\n", name, path)
			return nil
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "Auth token (prompted if omitted)")
	return cmd
}

func newLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout [registry]",
		Short: "Remove the stored auth token for a registry",
		Long: `Removes the authentication token for a registry from the skpm config file.

If no registry name is given, the configured default_registry is used.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			name := registryArg(args, cfg)
			if name == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove token for registry %q\n", name)
				return nil
			}

			if err := config.ClearRegistryToken(name); err != nil {
				return &InternalError{Message: "clear token", Cause: err}
			}

			path, _ := config.Path()
			fmt.Fprintf(cmd.OutOrStdout(), "Token removed for registry %q in %s\n", name, path)
			return nil
		},
	}
	return cmd
}

func registryArg(args []string, cfg *config.Config) string {
	if len(args) == 1 {
		return args[0]
	}
	return cfg.DefaultRegistry
}

func promptToken(cmd *cobra.Command, registryName string) (string, error) {
	fmt.Fprintf(cmd.OutOrStdout(), "Token for registry %q: ", registryName)

	if term.IsTerminal(int(syscall.Stdin)) {
		raw, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(cmd.OutOrStdout())
		if err != nil {
			return "", &InternalError{Message: "read token", Cause: err}
		}
		return string(raw), nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	if err := scanner.Err(); err != nil {
		return "", &InternalError{Message: "read token", Cause: err}
	}
	return "", &UserError{Message: "no input received"}
}
