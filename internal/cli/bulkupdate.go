package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
)

func newBulkUpdateCmd() *cobra.Command {
	var lockPath string
	var protectPath string
	var install bool

	cmd := &cobra.Command{
		Use:   "bulk-update <pattern>",
		Short: "Update all installed skills whose names match a shell glob",
		Long: `Matches installed skill names against <pattern> (supports * and ?) and runs
the upgrade flow for each matching skill that is not protected.

Protected skills (listed in ` + protectFile + `) are always skipped.

Example:
  skpm bulk-update 'lint-*'   # upgrades all skills starting with "lint-"
  skpm bulk-update '*'        # upgrades every unprotected skill`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read manifest: %v", err)}
			}

			// Match and filter skills.
			var matched []string
			for _, sl := range lf.Skills {
				if !matchesGlob(pattern, sl.Name) {
					continue
				}
				if protectPath != "" && IsProtected(protectPath, sl.Name) {
					fmt.Fprintf(cmd.OutOrStdout(), "  skip  %s (protected)\n", sl.Name)
					continue
				}
				matched = append(matched, sl.Name)
			}

			if len(matched) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No skills matched pattern %q.\n", pattern)
				return nil
			}

			upgraded := 0
			for _, name := range matched {
				src := ""
				for _, sl := range lf.Skills {
					if sl.Name == name {
						src = sl.Source
						break
					}
				}
				if src == "" {
					src = cfg.DefaultRegistry
				}

				reg, discovery, discErr := registryDiscovery(cmd.Context(), src, cfg)
				if discErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  skip  %s (registry error: %v)\n", name, discErr)
					continue
				}
				versions, listErr := discovery.ListVersions(cmd.Context(), registry.ParseSkillRef(name, registryDefaultNamespace(reg)))
				if listErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  skip  %s (list versions: %v)\n", name, listErr)
					continue
				}
				latest := latestVersion(versions)
				if latest == "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  skip  %s (no versions available)\n", name)
					continue
				}

				// Find current locked version.
				current := ""
				for _, sl := range lf.Skills {
					if sl.Name == name {
						current = sl.Version
						break
					}
				}
				if latest == current {
					fmt.Fprintf(cmd.OutOrStdout(), "  ok    %s (%s is already latest)\n", name, current)
					continue
				}

				if globalDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "  would update  %s  %s → %s\n", name, current, latest)
					upgraded++
					continue
				}

				// Update manifest.
				for i := range mf.Skills {
					if mf.Skills[i].Name == name {
						mf.Skills[i].Version = latest
						break
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  ↑  %s  %s → %s\n", name, current, latest)
				upgraded++
			}

			if upgraded == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All matched skills are already up to date.")
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would update %d skill(s).\n", upgraded)
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

			fmt.Fprintf(cmd.OutOrStdout(), "\nUpdated %d skill(s). Run 'skpm install' to apply.\n", upgraded)

			if install {
				return newInstallCmd().RunE(cmd, nil)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "", "Lockfile path (default: agent-skills.lock.yaml)")
	cmd.Flags().StringVar(&protectPath, "protect-file", protectFile, "Path to the protect store")
	cmd.Flags().BoolVar(&install, "install", false, "Run skpm install after updating")
	return cmd
}
