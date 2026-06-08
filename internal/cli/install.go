package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/internal/cache"
	"github.com/domehahn/skpm/internal/config"
	"github.com/domehahn/skpm/internal/installer"
	"github.com/domehahn/skpm/internal/lockfile"
	"github.com/domehahn/skpm/internal/manifest"
	"github.com/domehahn/skpm/internal/registry"
	"github.com/domehahn/skpm/internal/skill"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var lockPath string
	var frozenLockfile bool
	var check bool
	var prune bool
	var platform string
	var target string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install skills from agent-skills.lock (or resolve from agent-skills.yaml)",
		Long: `Installs all pinned skills.

If agent-skills.lock exists, installs exactly what is pinned (deterministic).
If agent-skills.lock is missing but agent-skills.yaml exists, resolves all
versions from the registry, generates a new agent-skills.lock, and installs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			if frozenLockfile {
				outdated, err := lockfileOutdated(cmd.Context(), cfg, manifest.DefaultFilename, lockPath)
				if err != nil {
					return &UserError{Message: fmt.Sprintf("check lockfile: %v", err)}
				}
				if outdated {
					return &UserError{Message: fmt.Sprintf("%s is inconsistent with %s", lockPath, manifest.DefaultFilename)}
				}
			}

			lf, err := resolveLockfile(cmd, lockPath, cfg)
			if err != nil {
				return err
			}
			if platform != "" {
				filterLockfileForPlatform(lf, platform)
			}
			fillMissingInstallPaths(lf)
			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to install.")
				return nil
			}
			if check {
				missing := missingInstallPaths(lf, target)
				if len(missing) > 0 {
					return &UserError{Message: fmt.Sprintf("installation incomplete:\n  %s", stringsJoin(missing, "\n  "))}
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Installation is complete.")
				return nil
			}

			c := cache.New(cfg.CacheDir)
			ins := installer.New(c)

			workDir, _ := os.Getwd()
			if target != "" {
				workDir = target
			}
			opts := installer.Options{
				DryRun:      globalDryRun,
				Concurrency: globalConcurrency,
				WorkDir:     workDir,
			}

			log.Debug().Str("lockfile", lockPath).Int("skills", len(lf.Skills)).Msg("installing")

			result, err := ins.Install(cmd.Context(), lf, opts)
			if err != nil {
				return &InternalError{Message: "install failed", Cause: err}
			}
			if prune && !globalDryRun {
				if err := pruneInstallPaths(workDir, lf); err != nil {
					return &InternalError{Message: "prune installed skills", Cause: err}
				}
			}

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "install",
					Data: map[string]interface{}{
						"installed":  result.Installed,
						"from_cache": result.FromCache,
						"skipped":    result.Skipped,
						"dry_run":    globalDryRun,
					},
				})
				return nil
			}

			total := len(result.Installed) + len(result.FromCache)
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would install %d skill(s)\n", len(result.Skipped))
				for _, name := range result.Skipped {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", name)
				}
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Installed %d skill(s)", total)
			if len(result.FromCache) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), " (%d from cache)", len(result.FromCache))
			}
			fmt.Fprintln(cmd.OutOrStdout())
			for _, name := range append(result.Installed, result.FromCache...) {
				sl, _ := lf.Find(name)
				if sl != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s@%s\n", name, sl.Version)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().BoolVar(&frozenLockfile, "frozen-lockfile", false, "Fail if manifest and lockfile are inconsistent")
	cmd.Flags().BoolVar(&check, "check", false, "Check whether locked skills are installed without writing files")
	cmd.Flags().BoolVar(&prune, "prune", false, "Remove installed skills no longer present in the lockfile")
	cmd.Flags().StringVar(&platform, "platform", "", "Install only skills compatible with a platform")
	cmd.Flags().StringVar(&target, "target", "", "Install into a target directory")
	return cmd
}

func filterLockfileForPlatform(lf *lockfile.LockFile, platform string) {
	out := lf.Skills[:0]
	for _, sl := range lf.Skills {
		if len(sl.CompatibleWith) == 0 {
			out = append(out, sl)
			continue
		}
		for _, p := range sl.CompatibleWith {
			if p == platform || p == "all" {
				out = append(out, sl)
				break
			}
		}
	}
	lf.Skills = out
}

func missingInstallPaths(lf *lockfile.LockFile, target string) []string {
	workDir, _ := os.Getwd()
	if target != "" {
		workDir = target
	}
	var missing []string
	for _, sl := range lf.Skills {
		for _, p := range sl.InstalledTo {
			if _, err := os.Stat(filepath.Join(workDir, p, "SKILL.md")); err != nil {
				missing = append(missing, fmt.Sprintf("%s: %s", sl.Name, p))
			}
		}
	}
	return missing
}

func fillMissingInstallPaths(lf *lockfile.LockFile) {
	for i := range lf.Skills {
		if len(lf.Skills[i].InstalledTo) > 0 {
			continue
		}
		platforms := stringsToPlatforms(lf.Skills[i].CompatibleWith)
		if len(platforms) == 0 {
			platforms = []skill.Platform{skill.PlatformAll}
		}
		paths, err := installer.ResolvePaths(lf.Skills[i].Name, platforms)
		if err == nil {
			lf.Skills[i].InstalledTo = paths
		}
	}
}

// resolveLockfile returns a ready-to-install LockFile.
// If the lockfile exists it is read directly.
// If it is missing but agent-skills.yaml exists, versions are resolved from
// the registry, a new lockfile is written, and that lockfile is returned.
func resolveLockfile(cmd *cobra.Command, lockPath string, cfg *config.Config) (*lockfile.LockFile, error) {
	if _, err := os.Stat(lockPath); err == nil {
		lf, err := lockfile.Read(lockPath)
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
		}
		return lf, nil
	}

	// lockfile missing — try manifest
	mf, err := manifest.Read(manifest.DefaultFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &UserError{Message: fmt.Sprintf(
				"neither %s nor %s found\nRun 'skpm init' to get started.",
				lockPath, manifest.DefaultFilename,
			)}
		}
		return nil, &UserError{Message: fmt.Sprintf("read manifest: %v", err)}
	}

	if len(mf.Skills) == 0 {
		lf := lockfile.New()
		return lf, nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "agent-skills.lock not found — resolving from %s\n\n", manifest.DefaultFilename)

	lf, err := resolveManifest(cmd, mf, cfg)
	if err != nil {
		return nil, err
	}

	if err := lf.Write(lockPath); err != nil {
		return nil, &InternalError{Message: "write lockfile", Cause: err}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "✓ Generated %s\n\n", lockPath)
	return lf, nil
}

// resolveManifest contacts the registry for each skill entry and builds a LockFile.
func resolveManifest(cmd *cobra.Command, mf *manifest.ManifestFile, cfg *config.Config) (*lockfile.LockFile, error) {
	lf := lockfile.New()

	for _, entry := range mf.Skills {
		src := entry.Source
		if src == "" {
			src = cfg.DefaultRegistry
		}
		if src == "" {
			return nil, &UserError{Message: fmt.Sprintf(
				"skill %q has no source and no default_registry is configured", entry.Name,
			)}
		}

		reg, err := registry.New(src, cfg)
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("registry for %q: %v", entry.Name, err)}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "  Resolving %s", entry.Name)
		if entry.Version != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "@%s", entry.Version)
		}
		fmt.Fprintln(cmd.OutOrStdout(), " ...")

		artifact, err := reg.Resolve(cmd.Context(), entry.Name, entry.Version)
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("resolve %s: %v", entry.Name, err)}
		}

		lf.Upsert(lockfile.SkillLock{
			Name:      entry.Name,
			Version:   artifact.Version,
			Source:    src,
			SourceURL: artifact.DownloadURL,
			SHA256:    artifact.SHA256,
		})
	}
	return lf, nil
}
