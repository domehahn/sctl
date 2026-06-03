package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/sctl/internal/cache"
	"github.com/domehahn/sctl/internal/config"
	"github.com/domehahn/sctl/internal/installer"
	"github.com/domehahn/sctl/internal/lockfile"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install skills from agent-skills.lock",
		Long: `Reads agent-skills.lock in the current directory (or --lock path),
downloads all skills, verifies SHA256 checksums, and installs them
atomically into the declared platform paths.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			if _, err := os.Stat(lockPath); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("lockfile not found: %s\nRun 'sctl add <skill>' to create one.", lockPath)}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}
			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to install — lockfile is empty.")
				return nil
			}

			c := cache.New(cfg.CacheDir)
			ins := installer.New(c)

			workDir, _ := os.Getwd()
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
				skill, _ := lf.Find(name)
				if skill != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s@%s\n", name, skill.Version)
				}
			}

			_ = filepath.Join
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	return cmd
}
