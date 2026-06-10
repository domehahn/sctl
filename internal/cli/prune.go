package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newPruneCmd() *cobra.Command {
	var lockPath string
	var roots []string

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove installed skill directories not present in the lockfile",
		Long: `Scans the known platform skill directories (.agents/skills, .claude/skills,
.github/skills, etc.) and removes any subdirectory that is not recorded
in the lockfile's installed_to paths.

Useful after 'skpm remove' to also clean up the installed files, or to
reconcile a directory that has drifted from the lockfile.

Use --dry-run to preview what would be removed without deleting anything.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}

			// Build the set of paths that should be kept.
			var kept map[string]bool
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				if !os.IsNotExist(err) {
					return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
				}
				kept = map[string]bool{} // treat missing lockfile as empty
			} else {
				kept = make(map[string]bool, len(lf.Skills)*3)
				for _, sl := range lf.Skills {
					for _, p := range sl.InstalledTo {
						kept[filepath.Clean(p)] = true
					}
				}
			}

			// Skill root directories to scan.
			scanRoots := roots
			if len(scanRoots) == 0 {
				scanRoots = defaultSkillRoots()
			}

			workDir, _ := os.Getwd()

			type pruneEntry struct {
				Path   string `json:"path"`
				Reason string `json:"reason"`
			}
			var toRemove []pruneEntry
			var kept_ int

			for _, root := range scanRoots {
				entries, err := os.ReadDir(filepath.Join(workDir, root))
				if err != nil {
					continue // root doesn't exist — skip silently
				}
				for _, e := range entries {
					if !e.IsDir() {
						continue
					}
					rel := filepath.Clean(filepath.Join(root, e.Name()))
					if kept[rel] {
						kept_++
						continue
					}
					// Check if it looks like a skill dir (has SKILL.md).
					if _, err := os.Stat(filepath.Join(workDir, rel, "SKILL.md")); err != nil {
						continue // not a skill dir, leave it alone
					}
					toRemove = append(toRemove, pruneEntry{Path: rel, Reason: "not in lockfile"})
				}
			}

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "prune",
					Data: map[string]interface{}{
						"to_remove": toRemove,
						"kept":      kept_,
						"dry_run":   globalDryRun,
					},
				})
				return nil
			}

			if len(toRemove) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to prune.")
				return nil
			}

			verb := "removed"
			if globalDryRun {
				verb = "would remove"
			}

			for _, entry := range toRemove {
				if !globalDryRun {
					if err := os.RemoveAll(filepath.Join(workDir, entry.Path)); err != nil {
						return &InternalError{Message: fmt.Sprintf("remove %s", entry.Path), Cause: err}
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s  %s\n", verb, entry.Path)
			}

			paths := make([]string, len(toRemove))
			for i, e := range toRemove {
				paths[i] = e.Path
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s %d director(ies): %s\n",
				strings.Title(verb), len(toRemove), strings.Join(paths, ", "))

			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().StringArrayVar(&roots, "root", nil, "Additional skill root directories to scan")
	return cmd
}

// defaultSkillRoots returns the standard platform skill directories.
func defaultSkillRoots() []string {
	return []string{
		".agents/skills",
		".claude/skills",
		".github/skills",
		".gitlab/skills",
		".cursor/skills",
		".opencode/skills",
		".kiro/skills",
		".roo/skills",
		".gemini/skills",
		".junie/skills",
		"skills",
	}
}
