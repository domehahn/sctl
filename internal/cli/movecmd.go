package cli

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newMoveCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "move <name> <new-namespace>",
		Short: "Update a skill's namespace in the lockfile",
		Long: `Sets the namespace: field on the named skill's lockfile entry to
<new-namespace>. The skill name and all other fields are unchanged.

Useful when reorganising skills into different namespace groups without
needing to reinstall.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, newNS := args[0], args[1]

			if newNS == "" || strings.ContainsAny(newNS, "\n\r:") {
				return &UserError{Message: fmt.Sprintf("invalid namespace %q: must not be empty or contain newlines or colons", newNS)}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			slPtr, ok := lf.Find(name)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", name)}
			}
			sl := *slPtr

			oldNS := sl.Namespace
			if oldNS == "" {
				oldNS = "default"
			}
			if oldNS == newNS {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already in namespace %q.\n", name, newNS)
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would move %s: %s → %s\n", name, oldNS, newNS)
				return nil
			}

			sl.Namespace = newNS
			lf.Upsert(sl)
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Moved %s: %s → %s\n", name, oldNS, newNS)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
