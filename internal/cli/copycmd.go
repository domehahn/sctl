package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newCopyCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "copy <source> <new-name>",
		Short: "Deep-copy a lockfile entry to a new independent name",
		Long: `Copies all fields of the <source> skill's lockfile entry to a new entry
named <new-name>. Unlike 'skpm alias', the new entry is fully independent:
it can be upgraded, protected, or removed without affecting the original.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source, newName := args[0], args[1]

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			src, ok := lf.Find(source)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", source)}
			}
			if _, exists := lf.Find(newName); exists {
				return &UserError{Message: fmt.Sprintf("skill %q already exists in lockfile", newName)}
			}

			copied := lockfile.SkillLock{
				Name:           newName,
				Version:        src.Version,
				Source:         src.Source,
				SHA256:         src.SHA256,
				InstalledTo:    append([]string(nil), src.InstalledTo...),
				Signature:      src.Signature,
				Provenance:     src.Provenance,
				CompatibleWith: append(src.CompatibleWith[:0:0], src.CompatibleWith...),
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would copy %q → %q (%s)\n", source, newName, src.Version)
				return nil
			}

			lf.Upsert(copied)
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Copied %q → %q (%s)\n", source, newName, src.Version)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
