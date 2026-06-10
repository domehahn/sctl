package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

const snapshotDir = ".skpm-snapshots"

func newSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save and restore lockfile snapshots",
	}
	cmd.AddCommand(newSnapshotSaveCmd())
	cmd.AddCommand(newSnapshotRestoreCmd())
	cmd.AddCommand(newSnapshotListCmd())
	cmd.AddCommand(newSnapshotDeleteCmd())
	return cmd
}

func newSnapshotSaveCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "save [name]",
		Short: "Save the current lockfile as a named snapshot",
		Long: `Copies the current agent-skills.lock to .skpm-snapshots/<name>.lock.

If no name is given, a timestamp is used (e.g. 2026-06-10T14-05-00).
Use 'skpm snapshot restore <name>' to roll back.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}
			if name == "" {
				name = time.Now().UTC().Format("2006-01-02T15-04-05")
			}
			if strings.ContainsAny(name, "/\\") {
				return &UserError{Message: "snapshot name must not contain path separators"}
			}

			data, err := os.ReadFile(lockfile.DefaultFilename)
			if err != nil {
				if os.IsNotExist(err) {
					return &UserError{Message: fmt.Sprintf("%s not found — nothing to snapshot", lockfile.DefaultFilename)}
				}
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			dest := snapshotPath(name)

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would save snapshot %q → %s\n", name, dest)
				return nil
			}

			if err := os.MkdirAll(snapshotDir, 0o755); err != nil {
				return &InternalError{Message: "create snapshot dir", Cause: err}
			}
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return &InternalError{Message: "write snapshot", Cause: err}
			}

			lf, _ := lockfile.Read(lockfile.DefaultFilename)
			skillCount := 0
			if lf != nil {
				skillCount = len(lf.Skills)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saved snapshot %q (%d skill(s)) → %s\n", name, skillCount, dest)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Snapshot name (defaults to current timestamp)")
	return cmd
}

func newSnapshotRestoreCmd() *cobra.Command {
	var backup bool

	cmd := &cobra.Command{
		Use:   "restore <name>",
		Short: "Restore a previously saved lockfile snapshot",
		Long: `Copies .skpm-snapshots/<name>.lock back to agent-skills.lock.

By default the current lockfile is backed up as a snapshot named "pre-restore"
before overwriting. Pass --no-backup to skip.

After restoring, run 'skpm install' to align installed files with the snapshot.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			src := snapshotPath(name)

			if _, err := os.Stat(src); err != nil {
				if os.IsNotExist(err) {
					return &UserError{Message: fmt.Sprintf("snapshot %q not found (run 'skpm snapshot list' to see available snapshots)", name)}
				}
				return &InternalError{Message: "stat snapshot", Cause: err}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would restore snapshot %q → %s\n", name, lockfile.DefaultFilename)
				return nil
			}

			// Backup current lockfile before overwriting.
			if backup {
				if _, err := os.Stat(lockfile.DefaultFilename); err == nil {
					backupName := "pre-restore-" + time.Now().UTC().Format("2006-01-02T15-04-05")
					current, err := os.ReadFile(lockfile.DefaultFilename)
					if err == nil {
						if mkErr := os.MkdirAll(snapshotDir, 0o755); mkErr == nil {
							_ = os.WriteFile(snapshotPath(backupName), current, 0o644)
							fmt.Fprintf(cmd.OutOrStdout(), "Backed up current lockfile as snapshot %q\n", backupName)
						}
					}
				}
			}

			data, err := os.ReadFile(src)
			if err != nil {
				return &InternalError{Message: "read snapshot", Cause: err}
			}
			if err := os.WriteFile(lockfile.DefaultFilename, data, 0o644); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			lf, _ := lockfile.Read(lockfile.DefaultFilename)
			skillCount := 0
			if lf != nil {
				skillCount = len(lf.Skills)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Restored snapshot %q (%d skill(s)) → %s\n", name, skillCount, lockfile.DefaultFilename)
			fmt.Fprintln(cmd.OutOrStdout(), "Run 'skpm install' to align installed files with the restored lockfile.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&backup, "backup", true, "Save current lockfile as 'pre-restore' snapshot before overwriting")
	return cmd
}

func newSnapshotListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := os.ReadDir(snapshotDir)
			if os.IsNotExist(err) || len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No snapshots found.")
				return nil
			}
			if err != nil {
				return &InternalError{Message: "read snapshot dir", Cause: err}
			}

			type row struct {
				Name   string `json:"name"`
				Skills int    `json:"skills"`
				Size   int64  `json:"size_bytes"`
			}
			var rows []row
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".lock") {
					continue
				}
				name := strings.TrimSuffix(e.Name(), ".lock")
				info, _ := e.Info()
				size := int64(0)
				if info != nil {
					size = info.Size()
				}
				lf, _ := lockfile.Read(snapshotPath(name))
				skills := 0
				if lf != nil {
					skills = len(lf.Skills)
				}
				rows = append(rows, row{Name: name, Skills: skills, Size: size})
			}
			sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })

			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "snapshot list", Data: rows})
				return nil
			}
			for _, r := range rows {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-40s  %d skill(s)\n", r.Name, r.Skills)
			}
			return nil
		},
	}
}

func newSnapshotDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			path := snapshotPath(name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("snapshot %q not found", name)}
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would delete snapshot %q\n", name)
				return nil
			}
			if err := os.Remove(path); err != nil {
				return &InternalError{Message: "delete snapshot", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted snapshot %q\n", name)
			return nil
		},
	}
}

func snapshotPath(name string) string {
	return filepath.Join(snapshotDir, name+".lock")
}
