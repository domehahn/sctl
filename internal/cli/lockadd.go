package cli

import (
	"fmt"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newLockAddCmd() *cobra.Command {
	var lockPath string
	var installedTo []string
	var signature string

	cmd := &cobra.Command{
		Use:   "lock-add <name> <version> <source> <sha256>",
		Short: "Manually insert a skill entry into the lockfile",
		Long: `Adds a fully-specified skill record to the lockfile without contacting a
registry. Designed for air-gapped environments where registry resolution is
unavailable but the artifact details (name, version, source URL, SHA256) are
known in advance.

Errors if the skill already exists unless --force is used (via --dry-run first).`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, source, sha256 := args[0], args[1], args[2], args[3]

			if len(sha256) != 64 {
				return &UserError{Message: "sha256 must be a 64-character hex string"}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			if _, exists := lf.Find(name); exists {
				return &UserError{Message: fmt.Sprintf("skill %q already exists in lockfile; remove it first with 'skpm remove'", name)}
			}

			entry := lockfile.SkillLock{
				Name:        name,
				Version:     version,
				Source:      source,
				SHA256:      sha256,
				InstalledTo: installedTo,
				Signature:   signature,
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would add %s@%s (source: %s)\n", name, version, source)
				return nil
			}

			lf.Upsert(entry)
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Added %s@%s to %s\n", name, version, lockPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringArrayVar(&installedTo, "installed-to", nil, "Installed path(s) (repeatable)")
	cmd.Flags().StringVar(&signature, "signature", "", "Optional HMAC signature")
	return cmd
}
