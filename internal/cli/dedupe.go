package cli

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newDedupeCmd() *cobra.Command {
	var lockPath string
	var fix bool

	cmd := &cobra.Command{
		Use:   "dedupe",
		Short: "Find and optionally remove duplicate skill entries in the lockfile",
		Long: `Scans the lockfile for skills that share an identical SHA-256 digest,
indicating they are the same artifact installed under different names (aliases).

Without --fix, reports duplicates only.
With --fix, removes the lower-priority duplicate (keeping the first alphabetically).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			// Group by SHA256 to find content duplicates.
			bySHA := map[string][]lockfile.SkillLock{}
			for _, sl := range lf.Skills {
				bySHA[sl.SHA256] = append(bySHA[sl.SHA256], sl)
			}

			type dupeGroup struct {
				SHA256 string
				Skills []string
			}
			var groups []dupeGroup
			for sha, skills := range bySHA {
				if len(skills) < 2 {
					continue
				}
				names := make([]string, len(skills))
				for i, s := range skills {
					names[i] = s.Name + "@" + s.Version
				}
				sort.Strings(names)
				groups = append(groups, dupeGroup{SHA256: sha[:12] + "...", Skills: names})
			}
			sort.Slice(groups, func(i, j int) bool { return groups[i].SHA256 < groups[j].SHA256 })

			if len(groups) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No duplicate skills found.")
				return nil
			}

			format := outputFormat()
			if format == OutputJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(CommandResult{
					Success: true,
					Command: "dedupe",
					Data:    groups,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Found %d group(s) of duplicate skills:\n\n", len(groups))
			toRemove := map[string]bool{}
			for _, g := range groups {
				fmt.Fprintf(cmd.OutOrStdout(), "  digest %s\n", g.SHA256)
				for i, name := range g.Skills {
					if i == 0 {
						fmt.Fprintf(cmd.OutOrStdout(), "    ✓  %s  (keep)\n", name)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "    ~  %s  (duplicate)\n", name)
						// Extract just the name part before '@'.
						n := name
						for j := 0; j < len(name); j++ {
							if name[j] == '@' {
								n = name[:j]
								break
							}
						}
						toRemove[n] = true
					}
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}

			if !fix {
				fmt.Fprintf(cmd.OutOrStdout(), "Run with --fix to remove %d duplicate(s)\n", len(toRemove))
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove %d duplicate(s)\n", len(toRemove))
				return nil
			}

			for name := range toRemove {
				lf.Remove(name)
			}
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed %d duplicate(s) from %s\n", len(toRemove), lockPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().BoolVar(&fix, "fix", false, "Remove duplicate entries from the lockfile")
	return cmd
}
