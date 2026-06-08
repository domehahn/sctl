package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

func newLockCmd() *cobra.Command {
	var check bool
	var update bool
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Resolve agent-skills.yaml into agent-skills.lock",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", manifest.DefaultFilename, err)}
			}
			resolved, err := lockFromManifest(cmd.Context(), mf, cfg, lockfile.DefaultFilename)
			if err != nil {
				return err
			}
			if check {
				current, err := lockfile.Read(lockfile.DefaultFilename)
				if err != nil || !sameLockedSkills(resolved, current) {
					return &UserError{Message: fmt.Sprintf("%s is outdated; run 'skpm lock'", lockfile.DefaultFilename)}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s is up to date\n", lockfile.DefaultFilename)
				return nil
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would write %s with %d skill(s)\n", lockfile.DefaultFilename, len(resolved.Skills))
				return nil
			}
			if !update {
				if current, err := lockfile.Read(lockfile.DefaultFilename); err == nil && sameLockedSkills(resolved, current) {
					fmt.Fprintf(cmd.OutOrStdout(), "%s is already up to date\n", lockfile.DefaultFilename)
					return nil
				}
			}
			if err := resolved.Write(lockfile.DefaultFilename); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Wrote %s (%d skill(s))\n", lockfile.DefaultFilename, len(resolved.Skills))
			return nil
		},
	}
	cmd.Flags().BoolVar(&check, "check", false, "Fail if agent-skills.lock is outdated")
	cmd.Flags().BoolVar(&update, "update", false, "Refresh resolved versions")
	return cmd
}
