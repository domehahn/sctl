package cli

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newSetMetaCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "set-meta <name> <key> <value>",
		Short: "Set an arbitrary key-value pair in a skill's lockfile metadata",
		Long: `Sets sl.Metadata[<key>] = <value> on the named lockfile entry.
Use 'skpm tag add/remove' for the specialised tags workflow.

Examples:
  skpm set-meta auth-skill owner "platform-team"
  skpm set-meta auth-skill review-date "2026-07-01"`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, key, value := args[0], args[1], args[2]

			if key == "" || strings.ContainsAny(key, "\n\r:") || strings.TrimSpace(key) != key {
				return &UserError{Message: fmt.Sprintf("invalid metadata key %q: must not be empty, contain newlines, colons, or surrounding whitespace", key)}
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

			if sl.Metadata != nil && sl.Metadata[key] == value {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s is already %q\n", name, key, value)
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would set %s.metadata[%s] = %q\n", name, key, value)
				return nil
			}

			if sl.Metadata == nil {
				sl.Metadata = map[string]string{}
			}
			sl.Metadata[key] = value
			lf.Upsert(sl)

			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Set %s.metadata[%s] = %q\n", name, key, value)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
