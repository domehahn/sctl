package cli

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
)

func newDeprecateCmd() *cobra.Command {
	var source string
	var reason string

	cmd := &cobra.Command{
		Use:   "deprecate <skill>@<version>",
		Short: "Mark a skill version as deprecated in the registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, err := parseSkillAtVersion(args[0])
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			src := sourceOrDefault(source, cfg)
			if src == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}

			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}
			gov, ok := reg.(registry.GovernanceRegistry)
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q does not support governance operations", src)}
			}

			ref := registry.SkillVersionRef{
				Namespace: "default",
				Name:      name,
				Version:   version,
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would deprecate %s@%s", name, version)
				if reason != "" {
					fmt.Fprintf(cmd.OutOrStdout(), " (reason: %s)", reason)
				}
				fmt.Fprintln(cmd.OutOrStdout())
				return nil
			}

			if err := gov.Deprecate(cmd.Context(), ref, reason); err != nil {
				return &InternalError{Message: "deprecate", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deprecated %s@%s\n", name, version)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (uses default_registry if not set)")
	cmd.Flags().StringVar(&reason, "reason", "", "Deprecation message shown to users")
	return cmd
}

func newYankCmd() *cobra.Command {
	var source string
	var reason string

	cmd := &cobra.Command{
		Use:   "yank <skill>@<version>",
		Short: "Yank a skill version from the registry (prevents new installs)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, err := parseSkillAtVersion(args[0])
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			src := sourceOrDefault(source, cfg)
			if src == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}

			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}
			gov, ok := reg.(registry.GovernanceRegistry)
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q does not support governance operations", src)}
			}

			ref := registry.SkillVersionRef{
				Namespace: "default",
				Name:      name,
				Version:   version,
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would yank %s@%s", name, version)
				if reason != "" {
					fmt.Fprintf(cmd.OutOrStdout(), " (reason: %s)", reason)
				}
				fmt.Fprintln(cmd.OutOrStdout())
				return nil
			}

			if err := gov.Yank(cmd.Context(), ref, reason); err != nil {
				return &InternalError{Message: "yank", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Yanked %s@%s\n", name, version)
			fmt.Fprintln(cmd.OutOrStdout(), "Note: existing installations are unaffected. Use 'skpm unyank' to reverse.")
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (uses default_registry if not set)")
	cmd.Flags().StringVar(&reason, "reason", "", "Yank reason shown to users")
	return cmd
}

func newUnyankCmd() *cobra.Command {
	var source string

	cmd := &cobra.Command{
		Use:   "unyank <skill>@<version>",
		Short: "Reverse a yank, making a skill version installable again",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, err := parseSkillAtVersion(args[0])
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			src := sourceOrDefault(source, cfg)
			if src == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}

			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}
			gov, ok := reg.(registry.GovernanceRegistry)
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q does not support governance operations", src)}
			}

			ref := registry.SkillVersionRef{
				Namespace: "default",
				Name:      name,
				Version:   version,
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would unyank %s@%s\n", name, version)
				return nil
			}

			if err := gov.Unyank(cmd.Context(), ref); err != nil {
				return &InternalError{Message: "unyank", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Unyanked %s@%s — version is now installable again\n", name, version)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (uses default_registry if not set)")
	return cmd
}

// parseSkillAtVersion splits "skill@version" into (name, version).
func parseSkillAtVersion(raw string) (name, version string, err error) {
	at := strings.LastIndex(raw, "@")
	if at < 1 || at == len(raw)-1 {
		return "", "", &UserError{Message: fmt.Sprintf("argument must be <skill>@<version>, got %q", raw)}
	}
	return raw[:at], raw[at+1:], nil
}
