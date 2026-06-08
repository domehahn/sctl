package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	var prune bool
	cmd := &cobra.Command{
		Use:   "remove <skill>",
		Short: "Remove a skill from agent-skills.yaml and refresh the lockfile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", manifest.DefaultFilename, err)}
			}
			removed := mf.Remove(name)
			if !removed {
				return &UserError{Message: fmt.Sprintf("%s is not declared in %s", name, manifest.DefaultFilename)}
			}
			lf, _ := lockfile.Read(lockfile.DefaultFilename)
			var removedPaths []string
			if lf != nil {
				if sl, ok := lf.Find(name); ok {
					removedPaths = append(removedPaths, sl.InstalledTo...)
				}
				lf.Remove(name)
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove %s\n", name)
				return nil
			}
			if err := mf.Write(manifest.DefaultFilename); err != nil {
				return &InternalError{Message: "write manifest", Cause: err}
			}
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			refreshed, err := lockFromManifest(cmd.Context(), mf, cfg, lockfile.DefaultFilename)
			if err != nil {
				if lf == nil {
					return err
				}
				refreshed = lf
			}
			if err := refreshed.Write(lockfile.DefaultFilename); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			if prune {
				workDir, _ := os.Getwd()
				for _, p := range removedPaths {
					if err := os.RemoveAll(filepath.Join(workDir, p)); err != nil {
						return &InternalError{Message: fmt.Sprintf("prune %s", p), Cause: err}
					}
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed %s\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&prune, "prune", false, "Remove installed skill files too")
	return cmd
}
