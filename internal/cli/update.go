package cli

import (
	"fmt"

	"github.com/domehahn/skpm/internal/config"
	"github.com/domehahn/skpm/internal/lockfile"
	"github.com/domehahn/skpm/internal/manifest"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var latest bool
	var install bool
	cmd := &cobra.Command{
		Use:   "update [skill]",
		Short: "Update one skill or all skills in the lockfile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", manifest.DefaultFilename, err)}
			}
			target := ""
			if len(args) == 1 {
				target = args[0]
			}
			if latest {
				for i := range mf.Skills {
					if target != "" && mf.Skills[i].Name != target {
						continue
					}
					src := mf.Skills[i].Source
					if src == "" {
						src = cfg.DefaultRegistry
					}
					discovery, err := registryDiscovery(cmd.Context(), src, cfg)
					if err != nil {
						return &UserError{Message: err.Error()}
					}
					versions, err := discovery.ListVersions(cmd.Context(), mf.Skills[i].Name)
					if err != nil {
						return &UserError{Message: err.Error()}
					}
					if v := latestVersion(versions); v != "" {
						mf.Skills[i].Version = v
					}
				}
				if !globalDryRun {
					if err := mf.Write(manifest.DefaultFilename); err != nil {
						return &InternalError{Message: "write manifest", Cause: err}
					}
				}
			}
			if target != "" {
				found := false
				for _, s := range mf.Skills {
					found = found || s.Name == target
				}
				if !found {
					return &UserError{Message: fmt.Sprintf("%s is not declared in %s", target, manifest.DefaultFilename)}
				}
			}
			lf, err := lockFromManifest(cmd.Context(), mf, cfg, lockfile.DefaultFilename)
			if err != nil {
				return err
			}
			if !globalDryRun {
				if err := lf.Write(lockfile.DefaultFilename); err != nil {
					return &InternalError{Message: "write lockfile", Cause: err}
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", lockfile.DefaultFilename)
			if install {
				return newInstallCmd().RunE(cmd, nil)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&latest, "latest", false, "Move to latest available version and update manifest constraints")
	cmd.Flags().BoolVar(&install, "install", false, "Install after updating the lockfile")
	return cmd
}
