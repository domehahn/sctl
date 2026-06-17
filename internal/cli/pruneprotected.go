package cli

import (
	"fmt"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newPruneProtectedCmd() *cobra.Command {
	var lockPath string
	var lockDir string

	cmd := &cobra.Command{
		Use:   "prune-protected",
		Short: "Remove stale entries from .skpm-protect for skills no longer in the lockfile",
		Long: `Reads ` + protectFile + ` and removes any skill name that is not present
in the lockfile. Keeps the protect file in sync after 'skpm remove' operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			installed := map[string]struct{}{}
			for _, sl := range lf.Skills {
				installed[sl.Name] = struct{}{}
			}

			storePath := filepath.Join(lockDir, protectFile)
			ps, err := loadProtectStore(storePath)
			if err != nil {
				return &InternalError{Message: "load protect store", Cause: err}
			}

			if len(ps.Protected) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No protected skills.")
				return nil
			}

			var kept, pruned []string
			for _, name := range ps.Protected {
				if _, ok := installed[name]; ok {
					kept = append(kept, name)
				} else {
					pruned = append(pruned, name)
					fmt.Fprintf(cmd.OutOrStdout(), "  prune  %s (not in lockfile)\n", name)
				}
			}

			if len(pruned) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No stale protected entries found.")
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would prune %d entry/entries.\n", len(pruned))
				return nil
			}

			ps.Protected = kept
			if err := ps.save(storePath); err != nil {
				return &InternalError{Message: "write protect store", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nPruned %d stale entry/entries from %s.\n", len(pruned), protectFile)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&lockDir, "lock-dir", ".", "Directory containing "+protectFile)
	return cmd
}
