package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newLockDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock-diff <before.lock> <after.lock>",
		Short: "Show what changed between two lockfiles",
		Long: `Compares two agent-skills.lock files and reports:
  + added    — skills present in <after> but not <before>
  ~ updated  — skills present in both but with a different version
  - removed  — skills present in <before> but not <after>

Useful for reviewing lockfile changes in pull requests without reading raw YAML.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			beforePath, afterPath := args[0], args[1]

			before, err := readLockOrEmpty(beforePath)
			if err != nil {
				return err
			}
			after, err := readLockOrEmpty(afterPath)
			if err != nil {
				return err
			}

			type change struct {
				kind      string // added, removed, updated
				name      string
				beforeVer string
				afterVer  string
				beforeSHA string
				afterSHA  string
			}

			beforeMap := make(map[string]lockfile.SkillLock, len(before.Skills))
			for _, sl := range before.Skills {
				beforeMap[sl.Name] = sl
			}
			afterMap := make(map[string]lockfile.SkillLock, len(after.Skills))
			for _, sl := range after.Skills {
				afterMap[sl.Name] = sl
			}

			var changes []change

			for _, sl := range after.Skills {
				if bsl, ok := beforeMap[sl.Name]; ok {
					if bsl.Version != sl.Version || bsl.SHA256 != sl.SHA256 {
						changes = append(changes, change{
							kind:      "updated",
							name:      sl.Name,
							beforeVer: bsl.Version,
							afterVer:  sl.Version,
							beforeSHA: bsl.SHA256,
							afterSHA:  sl.SHA256,
						})
					}
				} else {
					changes = append(changes, change{
						kind:     "added",
						name:     sl.Name,
						afterVer: sl.Version,
						afterSHA: sl.SHA256,
					})
				}
			}
			for _, sl := range before.Skills {
				if _, ok := afterMap[sl.Name]; !ok {
					changes = append(changes, change{
						kind:      "removed",
						name:      sl.Name,
						beforeVer: sl.Version,
						beforeSHA: sl.SHA256,
					})
				}
			}

			format := outputFormat()
			if format == OutputJSON {
				type jsonChange struct {
					Kind      string `json:"kind"`
					Name      string `json:"name"`
					Before    string `json:"before,omitempty"`
					After     string `json:"after,omitempty"`
					BeforeSHA string `json:"before_sha,omitempty"`
					AfterSHA  string `json:"after_sha,omitempty"`
				}
				var out []jsonChange
				for _, c := range changes {
					out = append(out, jsonChange{
						Kind:      c.kind,
						Name:      c.name,
						Before:    c.beforeVer,
						After:     c.afterVer,
						BeforeSHA: abbrev(c.beforeSHA),
						AfterSHA:  abbrev(c.afterSHA),
					})
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(CommandResult{
					Success: true,
					Command: "lock-diff",
					Data:    out,
				})
			}

			if len(changes) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No changes between lockfiles.")
				return nil
			}

			icons := map[string]string{"added": "+", "removed": "-", "updated": "~"}
			for _, c := range changes {
				icon := icons[c.kind]
				switch c.kind {
				case "added":
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  %-28s  %s  [%s]\n", icon, c.name, c.afterVer, abbrev(c.afterSHA))
				case "removed":
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  %-28s  %s  [%s]\n", icon, c.name, c.beforeVer, abbrev(c.beforeSHA))
				case "updated":
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  %-28s  %s → %s\n", icon, c.name, c.beforeVer, c.afterVer)
				}
			}

			added, removed, updated := 0, 0, 0
			for _, c := range changes {
				switch c.kind {
				case "added":
					added++
				case "removed":
					removed++
				case "updated":
					updated++
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n+%d added  ~%d updated  -%d removed\n", added, updated, removed)
			return nil
		},
	}
	return cmd
}

func readLockOrEmpty(path string) (*lockfile.LockFile, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return lockfile.New(), nil
	}
	lf, err := lockfile.Read(path)
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("read %s: %v", path, err)}
	}
	return lf, nil
}

func abbrev(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}
