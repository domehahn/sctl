package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newRollbackCmd() *cobra.Command {
	var toVersion string
	var lockPath string
	var list bool

	cmd := &cobra.Command{
		Use:   "rollback <skill-name>",
		Short: "Roll back a single skill to a previous version from a snapshot",
		Long: `Scans .skpm-snapshots/ for earlier versions of a skill and restores
the matching lockfile entry. Only that skill's entry is changed —
other skills are unaffected.

After rollback, run 'skpm install' to align installed files.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			skillName := args[0]
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}

			// Load current lockfile.
			currentLF, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}
			current, currentFound := currentLF.Find(skillName)
			currentVersion := ""
			if currentFound && current != nil {
				currentVersion = current.Version
			}

			// Scan snapshots for prior entries of this skill.
			type candidate struct {
				snapName string
				entry    lockfile.SkillLock
			}
			var candidates []candidate

			snapEntries, readErr := os.ReadDir(snapshotDir)
			if readErr == nil {
				for _, e := range snapEntries {
					if e.IsDir() || !strings.HasSuffix(e.Name(), ".lock") {
						continue
					}
					name := strings.TrimSuffix(e.Name(), ".lock")
					snapLF, snapErr := lockfile.Read(snapshotPath(name))
					if snapErr != nil {
						continue
					}
					sl, found := snapLF.Find(skillName)
					if !found || sl == nil {
						continue
					}
					if sl.Version == currentVersion {
						continue // same version — not a rollback target
					}
					candidates = append(candidates, candidate{snapName: name, entry: *sl})
				}
			}

			// Deduplicate by version — keep the most recent snapshot for each version.
			seen := map[string]candidate{}
			for _, c := range candidates {
				if _, exists := seen[c.entry.Version]; !exists {
					seen[c.entry.Version] = c
				}
			}
			var unique []candidate
			for _, c := range seen {
				unique = append(unique, c)
			}
			sort.Slice(unique, func(i, j int) bool {
				return unique[i].entry.Version > unique[j].entry.Version
			})

			if list || toVersion == "" {
				if len(unique) == 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "No previous versions of %q found in snapshots.\n", skillName)
					if len(candidates) == 0 {
						fmt.Fprintf(cmd.OutOrStdout(), "Run 'skpm snapshot save' before upgrading to enable rollback.\n")
					}
					return nil
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Available rollback targets for %q (current: %s):\n\n", skillName, currentVersion)
				for _, c := range unique {
					fmt.Fprintf(cmd.OutOrStdout(), "  %-16s  (snapshot: %s)\n", c.entry.Version, c.snapName)
				}
				if !list {
					fmt.Fprintf(cmd.OutOrStdout(), "\nRun: skpm rollback %s --to <version>\n", skillName)
				}
				return nil
			}

			// Find the candidate for the requested version.
			target, ok := seen[toVersion]
			if !ok {
				return &UserError{Message: fmt.Sprintf("version %q not found in any snapshot for %q", toVersion, skillName)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would roll back %s  %s → %s  (snapshot: %s)\n",
					skillName, currentVersion, toVersion, target.snapName)
				return nil
			}

			// Auto-save current state before overwriting.
			backupName := "pre-rollback-" + skillName + "-" + time.Now().UTC().Format("2006-01-02T15-04-05")
			if currentData, readErr := os.ReadFile(lockPath); readErr == nil {
				if mkErr := os.MkdirAll(snapshotDir, 0o755); mkErr == nil {
					_ = os.WriteFile(snapshotPath(backupName), currentData, 0o644)
					fmt.Fprintf(cmd.OutOrStdout(), "Saved pre-rollback snapshot: %s\n", backupName)
				}
			}

			// Apply rollback: replace only this skill's entry.
			currentLF.Upsert(target.entry)
			if err := currentLF.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Rolled back %s  %s → %s\n", skillName, currentVersion, toVersion)
			fmt.Fprintln(cmd.OutOrStdout(), "Run 'skpm install' to reinstall the rolled-back version.")
			return nil
		},
	}

	cmd.Flags().StringVar(&toVersion, "to", "", "Target version to roll back to")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().BoolVar(&list, "list", false, "List available rollback targets without making changes")
	return cmd
}
