package cli

import (
	"encoding/json"
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type repairIssue struct {
	Skill string `json:"skill"`
	Issue string `json:"issue"`
	Fixed bool   `json:"fixed,omitempty"`
}

func newRepairCmd() *cobra.Command {
	var lockPath string
	var fix bool

	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Scan the lockfile for integrity issues and optionally fix them",
		Long: `Checks every locked skill for common integrity problems:
  - missing SHA256 digest
  - missing source URL
  - duplicate name entries

Use --fix to automatically remove unfixable entries (duplicates, missing required fields).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			var issues []repairIssue
			seen := map[string]int{} // name → index of first occurrence
			var keep []lockfile.SkillLock

			for _, sl := range lf.Skills {
				var skillIssues []string
				if sl.SHA256 == "" {
					skillIssues = append(skillIssues, "missing SHA256")
				}
				if sl.Source == "" {
					skillIssues = append(skillIssues, "missing source")
				}
				if _, dup := seen[sl.Name]; dup {
					skillIssues = append(skillIssues, "duplicate name")
				} else {
					seen[sl.Name] = 1
				}

				if len(skillIssues) == 0 {
					keep = append(keep, sl)
					continue
				}

				for _, iss := range skillIssues {
					issues = append(issues, repairIssue{Skill: sl.Name, Issue: iss, Fixed: fix})
				}
				// Don't add to keep — we drop problematic entries when --fix is set.
			}

			if len(issues) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No issues found in lockfile.")
				return nil
			}

			format := outputFormat()
			if format == OutputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(issues)
			}

			for _, iss := range issues {
				status := "found"
				if fix {
					status = "fixed"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s: %s\n", status, iss.Skill, iss.Issue)
			}

			if !fix {
				fmt.Fprintf(cmd.OutOrStdout(), "\n%d issue(s) found. Run with --fix to remove problematic entries.\n", len(issues))
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would remove %d entry/entries with issues.\n", len(lf.Skills)-len(keep))
				return nil
			}

			// Rebuild lockfile with only clean entries.
			newLF := lockfile.New()
			for _, sl := range keep {
				newLF.Upsert(sl)
			}
			if err := newLF.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nRemoved %d entry/entries with issues.\n", len(lf.Skills)-len(keep))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().BoolVar(&fix, "fix", false, "Remove entries with integrity issues")
	return cmd
}
