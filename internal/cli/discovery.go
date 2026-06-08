package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	var source string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search a skill registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			src := sourceOrDefault(source, cfg)
			_, discovery, err := registryDiscovery(cmd.Context(), src, cfg)
			if err != nil {
				return &UserError{Message: err.Error()}
			}
			results, err := discovery.Search(cmd.Context(), registry.SearchRequest{Query: args[0]})
			if err != nil {
				return &UserError{Message: err.Error()}
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "search", Data: results})
				return nil
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%-28s %-12s %s\n", r.Name, r.Version, r.Description)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", "", "Registry source")
	return cmd
}

func newInfoCmd() *cobra.Command {
	var source string
	var versions bool
	cmd := &cobra.Command{
		Use:   "info <skill>",
		Short: "Show registry information for a skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			src := sourceOrDefault(source, cfg)
			reg, discovery, err := registryDiscovery(cmd.Context(), src, cfg)
			if err != nil {
				return &UserError{Message: err.Error()}
			}
			info, err := discovery.Info(cmd.Context(), registry.ParseSkillRef(args[0], registryDefaultNamespace(reg)))
			if err != nil {
				return &UserError{Message: err.Error()}
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "info", Data: info})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", info.Name)
			if info.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", info.Description)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  latest: %s\n", info.LatestVersion)
			if versions {
				for _, v := range info.Versions {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", v.Version)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", "", "Registry source")
	cmd.Flags().BoolVar(&versions, "versions", false, "Show all versions")
	return cmd
}

func newOutdatedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outdated",
		Short: "Show locked skills with newer registry versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			mf, _ := manifest.Read(manifest.DefaultFilename)
			lf, err := readCurrentLockfile()
			if err != nil {
				return err
			}
			type row struct {
				Name             string `json:"name"`
				Current          string `json:"current"`
				LatestCompatible string `json:"latest_compatible"`
				LatestOverall    string `json:"latest_overall"`
				Constraint       string `json:"constraint"`
			}
			var rows []row
			for _, sl := range lf.Skills {
				reg, discovery, err := registryDiscovery(cmd.Context(), sl.Source, cfg)
				if err != nil {
					continue
				}
				versions, err := discovery.ListVersions(cmd.Context(), registry.ParseSkillRef(sl.Name, registryDefaultNamespace(reg)))
				if err != nil {
					continue
				}
				constraint := manifestConstraint(mf, sl.Name)
				latest := latestVersion(versions)
				if latest != "" && latest != sl.Version {
					rows = append(rows, row{Name: sl.Name, Current: sl.Version, LatestCompatible: latest, LatestOverall: latest, Constraint: constraint})
				}
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "outdated", Data: rows})
				return nil
			}
			for _, r := range rows {
				fmt.Fprintf(cmd.OutOrStdout(), "%-28s current %-12s latest %-12s constraint %s\n", r.Name, r.Current, r.LatestOverall, r.Constraint)
			}
			if len(rows) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All locked skills are current or registries do not expose version discovery.")
			}
			return nil
		},
	}
	return cmd
}

func sourceOrDefault(source string, cfg *config.Config) string {
	if source != "" {
		return source
	}
	return cfg.DefaultRegistry
}

func manifestConstraint(mf *manifest.ManifestFile, name string) string {
	if mf == nil {
		return ""
	}
	for _, s := range mf.Skills {
		if s.Name == name {
			return s.Version
		}
	}
	return ""
}

func readCurrentLockfile() (*lockfile.LockFile, error) {
	lf, err := lockfile.Read(lockfile.DefaultFilename)
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("read %s: %v", lockfile.DefaultFilename, err)}
	}
	return lf, nil
}
