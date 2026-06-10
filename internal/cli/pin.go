package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

func newPinCmd() *cobra.Command {
	var manifestPath string
	var lockPath string

	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin manifest version constraints to exact locked versions",
		Long: `Replaces every version constraint in agent-skills.yaml with the exact version
recorded in agent-skills.lock. The result is a fully reproducible manifest
where every skill resolves to the same artifact regardless of registry state.

Run "skpm lock && skpm install" after pinning to verify the result.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mf, err := manifest.Read(manifestPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read manifest: %v", err)}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			locked := make(map[string]string, len(lf.Skills))
			for _, sl := range lf.Skills {
				locked[sl.Name] = sl.Version
			}

			type pinRow struct {
				Name   string `json:"name"`
				Before string `json:"before"`
				After  string `json:"after"`
				Pinned bool   `json:"pinned"`
			}

			var rows []pinRow
			for i, s := range mf.Skills {
				v, ok := locked[s.Name]
				if !ok {
					rows = append(rows, pinRow{Name: s.Name, Before: s.Version, After: s.Version})
					continue
				}
				before := s.Version
				if before == v {
					rows = append(rows, pinRow{Name: s.Name, Before: before, After: v})
					continue
				}
				if !globalDryRun {
					mf.Skills[i].Version = v
				}
				rows = append(rows, pinRow{Name: s.Name, Before: before, After: v, Pinned: true})
			}

			if !globalDryRun {
				if err := mf.Write(manifestPath); err != nil {
					return &InternalError{Message: "write manifest", Cause: err}
				}
			}

			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "pin", Data: rows})
				return nil
			}

			pinned := 0
			for _, r := range rows {
				if r.Pinned {
					fmt.Fprintf(cmd.OutOrStdout(), "pinned  %-32s  %s  →  %s\n", r.Name, r.Before, r.After)
					pinned++
				}
			}

			switch {
			case pinned == 0:
				fmt.Fprintln(cmd.OutOrStdout(), "All skills are already pinned to exact versions.")
			case globalDryRun:
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: %d skill(s) would be pinned in %s\n", pinned, manifestPath)
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "\nPinned %d skill(s) in %s\n", pinned, manifestPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&manifestPath, "manifest", manifest.DefaultFilename, "Path to agent-skills.yaml")
	cmd.Flags().StringVar(&lockPath, "lock", lockfile.DefaultFilename, "Path to agent-skills.lock")
	return cmd
}
