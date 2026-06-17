package cli

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newRenameSkillCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "rename <old-name> <new-name>",
		Short: "Rename a lockfile entry in-place, preserving all other fields",
		Long: `Updates the Name field of a lockfile entry from <old-name> to <new-name>
and re-keys the entry. All other fields (version, source, SHA256, namespace,
metadata, etc.) are preserved.

This is different from 'skpm move', which only changes the namespace.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName, newName := args[0], args[1]

			if strings.ContainsAny(newName, "\n\r:") || strings.TrimSpace(newName) != newName || newName == "" {
				return &UserError{Message: fmt.Sprintf("invalid skill name %q: must not contain newlines, colons, or surrounding whitespace", newName)}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			if _, exists := lf.Find(newName); exists {
				return &UserError{Message: fmt.Sprintf("skill %q already exists in lockfile", newName)}
			}

			slPtr, ok := lf.Find(oldName)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", oldName)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would rename %s → %s\n", oldName, newName)
				return nil
			}

			sl := *slPtr
			lf.Remove(oldName)
			sl.Name = newName
			lf.Upsert(sl)

			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Renamed %s → %s\n", oldName, newName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
