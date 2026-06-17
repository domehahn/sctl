package cli

import (
	"bytes"
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newFormatLockCmd() *cobra.Command {
	var lockPath string
	var check bool

	cmd := &cobra.Command{
		Use:   "format-lock",
		Short: "Normalise and sort the lockfile to a canonical form",
		Long: `Sorts lockfile entries alphabetically by name, removes zero-value optional
fields, and rewrites the file. Idempotent: running it twice produces the same output.

With --check, exits non-zero if the file is not already in canonical form (useful in CI).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			// Capture the current on-disk bytes for comparison.
			before, readErr := os.ReadFile(lockPath)
			if readErr != nil {
				return &InternalError{Message: "read lockfile bytes", Cause: readErr}
			}

			// Sort is always applied by lf.Write via lf.Sort().
			lf.Sort()

			// Write to a buffer to compare without touching disk.
			var buf bytes.Buffer
			tmp := lockPath + ".tmp"
			if err := lf.Write(tmp); err != nil {
				return &InternalError{Message: "write temp lockfile", Cause: err}
			}
			after, err := os.ReadFile(tmp)
			_ = os.Remove(tmp)
			if err != nil {
				return &InternalError{Message: "read temp lockfile", Cause: err}
			}
			_ = buf // not needed

			if bytes.Equal(before, after) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already in canonical form\n", lockPath)
				return nil
			}

			if check {
				fmt.Fprintf(cmd.ErrOrStderr(), "%s is not in canonical form — run 'skpm format-lock' to fix\n", lockPath)
				return fmt.Errorf("lockfile not canonical")
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would rewrite %s (%d → %d bytes)\n",
					lockPath, len(before), len(after))
				return nil
			}

			if err := os.WriteFile(lockPath, after, 0o644); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Reformatted %s (%d → %d bytes)\n", lockPath, len(before), len(after))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().BoolVar(&check, "check", false, "Exit non-zero if file is not already canonical (no writes)")
	return cmd
}
