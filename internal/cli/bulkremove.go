package cli

import (
	"fmt"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newBulkRemoveCmd() *cobra.Command {
	var lockPath string
	var lockDir string

	cmd := &cobra.Command{
		Use:   "bulk-remove <pattern>",
		Short: "Remove all lockfile entries matching a glob pattern, skipping protected skills",
		Long: `Finds every lockfile entry whose name matches <pattern> (glob syntax: * and ?)
and removes it. Protected skills (in ` + protectFile + `) are skipped.

Use --dry-run to preview without making changes.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			storePath := filepath.Join(lockDir, protectFile)
			ps, err := loadProtectStore(storePath)
			if err != nil {
				return &InternalError{Message: "load protect store", Cause: err}
			}

			var toRemove, skipped []string
			for _, sl := range lf.Skills {
				if !matchesGlob(pattern, sl.Name) {
					continue
				}
				if ps.isProtected(sl.Name) {
					skipped = append(skipped, sl.Name)
					fmt.Fprintf(cmd.OutOrStdout(), "  skip (protected)  %s\n", sl.Name)
					continue
				}
				toRemove = append(toRemove, sl.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "  remove  %s\n", sl.Name)
			}

			if len(toRemove)+len(skipped) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No skills match %q.\n", pattern)
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would remove %d, skip %d (protected).\n", len(toRemove), len(skipped))
				return nil
			}

			if len(toRemove) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\nAll matched skills are protected — nothing removed.\n")
				return nil
			}

			for _, name := range toRemove {
				lf.Remove(name)
			}
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile (no skills were removed — lockfile unchanged on disk)", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nRemoved %d skill(s), skipped %d (protected).\n", len(toRemove), len(skipped))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&lockDir, "lock-dir", ".", "Directory containing "+protectFile)
	return cmd
}
