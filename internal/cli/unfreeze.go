package cli

import (
	"encoding/json"
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

type unfreezeResult struct {
	Skill   string `json:"skill"`
	Cleared string `json:"cleared_version"`
}

func newUnfreezeCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "unfreeze",
		Short: "Remove exact version pins from the manifest, re-enabling upgrades",
		Long: `Reads the manifest and clears the version constraint for every skill whose
locked version matches the pinned manifest version exactly, allowing 'skpm upgrade'
and 'skpm bulk-update' to move to newer versions again.

Skills not present in the lockfile are left unchanged.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			locked := make(map[string]string, len(lf.Skills))
			for _, sl := range lf.Skills {
				locked[sl.Name] = sl.Version
			}

			var results []unfreezeResult
			for i, ms := range mf.Skills {
				v, ok := locked[ms.Name]
				if !ok || ms.Version == "" {
					continue
				}
				// Only unfreeze if the manifest version exactly matches the locked version
				// (i.e. it was previously frozen by 'skpm freeze').
				if ms.Version != v {
					continue
				}
				results = append(results, unfreezeResult{Skill: ms.Name, Cleared: ms.Version})
				mf.Skills[i].Version = ""
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No frozen skills found in manifest.")
				return nil
			}

			if globalDryRun {
				if outputFormat() == OutputJSON {
					return json.NewEncoder(cmd.OutOrStdout()).Encode(results)
				}
				for _, r := range results {
					fmt.Fprintf(cmd.OutOrStdout(), "  would unfreeze  %-28s  (was %s)\n", r.Skill, r.Cleared)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would unfreeze %d skill(s).\n", len(results))
				return nil
			}

			if err := mf.Write(manifest.DefaultFilename); err != nil {
				return &InternalError{Message: "write manifest", Cause: err}
			}

			if outputFormat() == OutputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(results)
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "  unfrozen  %-28s  (was %s)\n", r.Skill, r.Cleared)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nUnfroze %d skill(s) in %s.\n", len(results), manifest.DefaultFilename)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "", "Lockfile path (default: agent-skills.lock.yaml)")
	return cmd
}
