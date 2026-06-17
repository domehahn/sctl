package cli

import (
	"encoding/json"
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

type freezeResult struct {
	Skill   string `json:"skill"`
	Version string `json:"version"`
}

func newFreezeCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "freeze",
		Short: "Snapshot locked versions as exact constraints in the manifest",
		Long: `Reads every skill's locked version from the lockfile and writes it as the
exact version constraint in agent-skills.yaml.

After freezing, 'skpm upgrade' and 'skpm bulk-update' will not change these
skills unless you manually edit the manifest or run 'skpm unfreeze'.`,
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

			// Build a version map from the lockfile.
			locked := make(map[string]string, len(lf.Skills))
			for _, sl := range lf.Skills {
				locked[sl.Name] = sl.Version
			}

			var results []freezeResult
			for i, ms := range mf.Skills {
				v, ok := locked[ms.Name]
				if !ok {
					continue
				}
				if mf.Skills[i].Version == v {
					continue
				}
				mf.Skills[i].Version = v
				results = append(results, freezeResult{Skill: ms.Name, Version: v})
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All manifest versions already match lockfile.")
				return nil
			}

			if globalDryRun {
				if outputFormat() == OutputJSON {
					return json.NewEncoder(cmd.OutOrStdout()).Encode(results)
				}
				for _, r := range results {
					fmt.Fprintf(cmd.OutOrStdout(), "  would freeze  %-28s  %s\n", r.Skill, r.Version)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would freeze %d skill(s).\n", len(results))
				return nil
			}

			if err := mf.Write(manifest.DefaultFilename); err != nil {
				return &InternalError{Message: "write manifest", Cause: err}
			}

			if outputFormat() == OutputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(results)
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "  frozen  %-28s  %s\n", r.Skill, r.Version)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nFroze %d skill(s) in %s.\n", len(results), manifest.DefaultFilename)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "", "Lockfile path (default: agent-skills.lock.yaml)")
	return cmd
}
