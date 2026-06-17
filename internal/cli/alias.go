package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newAliasCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "alias <source> <alias>",
		Short: "Add a second lockfile entry pointing to the same artifact as an installed skill",
		Long: `Creates an alias: a new lockfile entry with the given <alias> name that shares
the same version, SHA-256, source, and install paths as <source>.

Useful when two targets need to reference the same skill under different names
without downloading it twice.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source, alias := args[0], args[1]

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			sl, ok := lf.Find(source)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in %s", source, lockPath)}
			}

			if _, exists := lf.Find(alias); exists {
				return &UserError{Message: fmt.Sprintf("alias %q already exists in %s", alias, lockPath)}
			}

			entry := *sl
			entry.Name = alias

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would add alias %q → %s@%s\n", alias, source, sl.Version)
				return nil
			}

			lf.Upsert(entry)
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Added alias %q → %s@%s\n", alias, source, sl.Version)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
