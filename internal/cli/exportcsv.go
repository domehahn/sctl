package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newExportCSVCmd() *cobra.Command {
	var lockPath string
	var out string

	cmd := &cobra.Command{
		Use:   "export-csv",
		Short: "Export the lockfile as CSV for spreadsheet and compliance tools",
		Long: `Writes one row per skill: name, version, source, sha256, namespace, tags.
Use --out to write to a file, or omit (default: stdout).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			var w *csv.Writer
			var outFile *os.File
			if out == "" || out == "-" {
				w = csv.NewWriter(cmd.OutOrStdout())
			} else {
				f, err := os.Create(out)
				if err != nil {
					return fmt.Errorf("create %s: %w", out, err)
				}
				outFile = f
				w = csv.NewWriter(f)
			}

			if err := w.Write([]string{"name", "version", "source", "sha256", "namespace", "tags"}); err != nil {
				if outFile != nil {
					outFile.Close()
				}
				return fmt.Errorf("write CSV header: %w", err)
			}
			for _, sl := range lf.Skills {
				ns := sl.Namespace
				if ns == "" {
					ns = "default"
				}
				tags := strings.Join(skillTags(sl), ";")
				_ = w.Write([]string{sl.Name, sl.Version, sl.Source, sl.SHA256, ns, tags})
			}
			w.Flush()
			if err := w.Error(); err != nil {
				if outFile != nil {
					outFile.Close()
				}
				return fmt.Errorf("write CSV: %w", err)
			}

			if outFile != nil {
				if err := outFile.Close(); err != nil {
					return fmt.Errorf("close %s: %w", out, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote %s (%d row(s))\n", out, len(lf.Skills))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&out, "out", "-", "Output path; use '-' for stdout")
	return cmd
}
