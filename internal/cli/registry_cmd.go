package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "registry", Short: "Manage configured skill registries"}
	cmd.AddCommand(newRegistryListCmd())
	cmd.AddCommand(newRegistryShowCmd())
	cmd.AddCommand(newRegistryAddCmd())
	cmd.AddCommand(newRegistryRemoveCmd())
	cmd.AddCommand(newRegistryTestCmd())
	cmd.AddCommand(newRegistryCapabilitiesCmd())
	cmd.AddCommand(newRegistryLoginCmd())
	return cmd
}

func newRegistryListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured registries",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			type row struct {
				Name    string `json:"name"`
				Type    string `json:"type"`
				Default bool   `json:"default"`
			}
			var rows []row
			for name, rc := range cfg.Registries {
				rows = append(rows, row{Name: name, Type: rc.Type, Default: name == cfg.DefaultRegistry})
			}
			sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "registry list", Data: rows})
				return nil
			}
			for _, r := range rows {
				marker := " "
				if r.Default {
					marker = "*"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s %-28s %s\n", marker, r.Name, r.Type)
			}
			return nil
		},
	}
}

func newRegistryShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show one registry configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			rc, ok := cfg.Registries[args[0]]
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q is not configured", args[0])}
			}
			if rc.Token != "" {
				rc.Token = "***"
			}
			if rc.Auth.Token != "" {
				rc.Auth.Token = "***"
			}
			if rc.Auth.Password != "" {
				rc.Auth.Password = "***"
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "registry show", Data: rc})
				return nil
			}
			data, _ := yaml.Marshal(rc)
			fmt.Fprint(cmd.OutOrStdout(), string(data))
			return nil
		},
	}
}

func newRegistryAddCmd() *cobra.Command {
	var rc config.RegistryConfig
	var tokenEnv string
	var endpointResolve string
	var endpointDownload string
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add or update a registry configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := resolveConfigPath("")
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}
			cfg, err := config.LoadFrom(path)
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			if cfg.Registries == nil {
				cfg.Registries = map[string]config.RegistryConfig{}
			}
			rc.Name = args[0]
			if tokenEnv != "" {
				rc.Auth.Type = "bearer"
				rc.Auth.TokenEnv = tokenEnv
			}
			if rc.Endpoints == nil {
				rc.Endpoints = map[string]string{}
			}
			if endpointResolve != "" {
				rc.Endpoints["resolve"] = endpointResolve
			}
			if endpointDownload != "" {
				rc.Endpoints["download"] = endpointDownload
			}
			cfg.Registries[args[0]] = rc
			if cfg.DefaultRegistry == "" {
				cfg.DefaultRegistry = args[0]
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would write registry %s to %s\n", args[0], path)
				return nil
			}
			if err := writeConfig(path, cfg); err != nil {
				return &InternalError{Message: "write config", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved registry %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&rc.Type, "type", "", "Registry type: github, gitlab, artifactory, local, skillforge, generic-http")
	cmd.Flags().StringVar(&rc.URL, "url", "", "Registry base URL")
	cmd.Flags().StringVar(&rc.Repo, "repo", "", "Repository name or owner/repo")
	cmd.Flags().StringVar(&rc.Project, "project", "", "GitLab project path")
	cmd.Flags().StringVar(&rc.Path, "path", "", "Local registry path")
	cmd.Flags().StringVar(&rc.Namespace, "namespace", "default", "Default namespace")
	cmd.Flags().StringVar(&tokenEnv, "token-env", "", "Environment variable containing bearer token")
	cmd.Flags().StringToStringVar(&rc.Endpoints, "endpoint", nil, "Endpoint template map entry, e.g. resolve=/api/...")
	cmd.Flags().StringVar(&endpointResolve, "endpoint-resolve", "", "Resolve endpoint template")
	cmd.Flags().StringVar(&endpointDownload, "endpoint-download", "", "Download endpoint template")
	return cmd
}

func newRegistryRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a registry configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := resolveConfigPath("")
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}
			cfg, err := config.LoadFrom(path)
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			if _, ok := cfg.Registries[args[0]]; !ok {
				return &UserError{Message: fmt.Sprintf("registry %q is not configured", args[0])}
			}
			delete(cfg.Registries, args[0])
			if cfg.DefaultRegistry == args[0] {
				cfg.DefaultRegistry = ""
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove registry %s\n", args[0])
				return nil
			}
			if err := writeConfig(path, cfg); err != nil {
				return &InternalError{Message: "write config", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed registry %s\n", args[0])
			return nil
		},
	}
}

func newRegistryTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test <name>",
		Short: "Test registry connectivity and capabilities",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, caps, err := registryCapabilities(cmd.Context(), args[0])
			if err != nil {
				return &UserError{Message: err.Error()}
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "registry test", Data: map[string]interface{}{"name": reg.Name(), "type": reg.Type(), "capabilities": caps}})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Registry %s (%s) is reachable\n", reg.Name(), reg.Type())
			printCapabilities(cmd, caps)
			return nil
		},
	}
}

func newRegistryCapabilitiesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities <name>",
		Short: "Show registry capabilities",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, caps, err := registryCapabilities(cmd.Context(), args[0])
			if err != nil {
				return &UserError{Message: err.Error()}
			}
			data := map[string]interface{}{
				"name":               reg.Name(),
				"type":               reg.Type(),
				"resolve":            caps.Resolve,
				"download":           caps.Download,
				"search":             caps.Search,
				"info":               caps.Info,
				"publish":            caps.Publish,
				"deprecate":          caps.Deprecate,
				"yank":               caps.Yank,
				"unyank":             caps.Unyank,
				"semver_constraints": caps.SemVerConstraints,
				"checksums":          caps.Checksums,
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "registry capabilities", Data: data})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s (%s)\n", reg.Name(), reg.Type())
			printCapabilities(cmd, caps)
			return nil
		},
	}
}

func newRegistryLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login <name>",
		Short: "Show auth guidance for a registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			rc, ok := cfg.Registries[args[0]]
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q is not configured", args[0])}
			}
			if rc.Auth.TokenEnv != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Set %s to authenticate registry %s.\n", rc.Auth.TokenEnv, args[0])
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Registry %s uses auth type %q. Configure auth.token_env to avoid storing secrets.\n", args[0], rc.Auth.Type)
			return nil
		},
	}
}

func registryCapabilities(ctx context.Context, name string) (registry.Registry, *registry.RegistryCapabilities, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	reg, err := registry.New(name, cfg)
	if err != nil {
		return nil, nil, err
	}
	caps, err := reg.Capabilities(ctx)
	if err != nil {
		return nil, nil, err
	}
	return reg, caps, nil
}

func printCapabilities(cmd *cobra.Command, caps *registry.RegistryCapabilities) {
	fmt.Fprintf(cmd.OutOrStdout(), "  resolve: %v\n", caps.Resolve)
	fmt.Fprintf(cmd.OutOrStdout(), "  download: %v\n", caps.Download)
	fmt.Fprintf(cmd.OutOrStdout(), "  search: %v\n", caps.Search)
	fmt.Fprintf(cmd.OutOrStdout(), "  info: %v\n", caps.Info)
	fmt.Fprintf(cmd.OutOrStdout(), "  publish: %v\n", caps.Publish)
	fmt.Fprintf(cmd.OutOrStdout(), "  deprecate: %v\n", caps.Deprecate)
	fmt.Fprintf(cmd.OutOrStdout(), "  yank: %v\n", caps.Yank)
	fmt.Fprintf(cmd.OutOrStdout(), "  unyank: %v\n", caps.Unyank)
	fmt.Fprintf(cmd.OutOrStdout(), "  semver_constraints: %v\n", caps.SemVerConstraints)
	fmt.Fprintf(cmd.OutOrStdout(), "  checksums: %v\n", caps.Checksums)
}

func writeConfig(path string, cfg *config.Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return writeAtomicWithMode(path, string(data), 0o600)
}
