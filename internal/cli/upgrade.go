package cli

import (
	"fmt"
	"sort"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
)

func newUpgradeCmd() *cobra.Command {
	var all bool
	var skills []string
	var lockPath string
	var install bool

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Preview and apply upgrades to outdated skills",
		Long: `Checks the registry for newer versions of locked skills, shows a preview
table, then applies the selected upgrades to the manifest and lockfile.

Without --all or --skill, upgrade prints the outdated table and exits (preview mode).
Use --all to upgrade every outdated skill, or --skill to pick specific ones.

After upgrading, run 'skpm install' to materialize the new versions on disk,
or pass --install to do it automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", manifest.DefaultFilename, err)}
			}
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			type candidate struct {
				name    string
				current string
				latest  string
				source  string
			}

			var candidates []candidate
			for _, sl := range lf.Skills {
				src := sl.Source
				if src == "" {
					src = cfg.DefaultRegistry
				}
				reg, discovery, discErr := registryDiscovery(cmd.Context(), src, cfg)
				if discErr != nil {
					continue
				}
				versions, listErr := discovery.ListVersions(cmd.Context(), registry.ParseSkillRef(sl.Name, registryDefaultNamespace(reg)))
				if listErr != nil {
					continue
				}
				latest := latestVersion(versions)
				if latest != "" && latest != sl.Version {
					candidates = append(candidates, candidate{
						name:    sl.Name,
						current: sl.Version,
						latest:  latest,
						source:  src,
					})
				}
			}

			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].name < candidates[j].name
			})

			if len(candidates) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All skills are up to date.")
				return nil
			}

			// Preview table — always printed.
			fmt.Fprintf(cmd.OutOrStdout(), "  %-28s  %-12s  %s\n", "skill", "current", "latest")
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", "──────────────────────────────────────────────────")
			for _, c := range candidates {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-28s  %-12s  %s\n", c.name, c.current, c.latest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d skill(s) can be upgraded.\n\n", len(candidates))

			// Determine which skills to upgrade.
			skillSet := make(map[string]bool, len(skills))
			for _, s := range skills {
				skillSet[s] = true
			}

			if !all && len(skillSet) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Pass --all to upgrade all, or --skill <name> to select specific skills.")
				return nil
			}

			upgraded := 0
			for _, c := range candidates {
				if !all && !skillSet[c.name] {
					continue
				}
				// Update manifest constraint to the latest version.
				for i := range mf.Skills {
					if mf.Skills[i].Name == c.name {
						mf.Skills[i].Version = c.latest
						break
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  ↑  %s  %s → %s\n", c.name, c.current, c.latest)
				upgraded++
			}

			if upgraded == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills matched the --skill filter.")
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would upgrade %d skill(s)\n", upgraded)
				return nil
			}

			if err := mf.Write(manifest.DefaultFilename); err != nil {
				return &InternalError{Message: "write manifest", Cause: err}
			}

			newLF, err := lockFromManifest(cmd.Context(), mf, cfg, lockPath)
			if err != nil {
				return err
			}
			if err := newLF.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nUpgraded %d skill(s). Run 'skpm install' to apply.\n", upgraded)

			if install {
				return newInstallCmd().RunE(cmd, nil)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Upgrade all outdated skills")
	cmd.Flags().StringArrayVar(&skills, "skill", nil, "Skill name(s) to upgrade (repeatable)")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Lockfile path (default: agent-skills.lock)")
	cmd.Flags().BoolVar(&install, "install", false, "Run skpm install after upgrading")
	return cmd
}
