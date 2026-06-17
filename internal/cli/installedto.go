package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newInstalledToCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "installed-to <path>",
		Short: "List skills installed at a given path (reverse InstalledTo lookup)",
		Long: `Finds every lockfile entry whose InstalledTo list contains <path> and
prints those skills. Useful in monorepos to audit what is deployed to a
specific directory.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			searchPath := args[0]

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			var found []lockfile.SkillLock
			for _, sl := range lf.Skills {
				for _, p := range sl.InstalledTo {
					if strings.EqualFold(filepath.Clean(p), filepath.Clean(searchPath)) {
						found = append(found, sl)
						break
					}
				}
			}

			if len(found) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No skills installed to %q.\n", searchPath)
				return nil
			}

			if strings.ToLower(globalOutput) == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(found)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%d skill(s) installed to %s:\n\n", len(found), searchPath)
			for _, sl := range found {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-32s  %s\n", sl.Name, sl.Version)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
