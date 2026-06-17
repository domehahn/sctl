package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

const defaultVendorDir = "vendor/skills"

func newVendorCmd() *cobra.Command {
	var lockPath string
	var vendorDir string
	var clean bool

	cmd := &cobra.Command{
		Use:   "vendor",
		Short: "Materialize installed skills into vendor/skills/ for hermetic builds",
		Long: `Copies every installed skill directory into a local vendor tree so that
builds can run without registry access (containers, air-gapped CI, monorepos
where the install path may not be stable across runs).

The vendor tree mirrors the lockfile: each skill gets its own subdirectory
named <skill-name> (or <namespace>/<skill-name> for namespaced skills).

Use --clean to remove the vendor dir before copying, ensuring a fresh state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			if vendorDir == "" {
				vendorDir = defaultVendorDir
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would vendor %d skill(s) into %s\n", len(lf.Skills), vendorDir)
				return nil
			}

			if clean {
				if err := os.RemoveAll(vendorDir); err != nil {
					return &InternalError{Message: "clean vendor dir", Cause: err}
				}
			}

			if err := os.MkdirAll(vendorDir, 0o755); err != nil {
				return &InternalError{Message: "create vendor dir", Cause: err}
			}

			vendored := 0
			for _, sl := range lf.Skills {
				dest := sl.Name
				if sl.Namespace != "" && sl.Namespace != "default" {
					dest = filepath.Join(sl.Namespace, sl.Name)
				}
				destPath := filepath.Join(vendorDir, dest)

				copied := false
				for _, installPath := range sl.InstalledTo {
					info, err := os.Stat(installPath)
					if err != nil || !info.IsDir() {
						continue
					}
					if err := os.MkdirAll(destPath, 0o755); err != nil {
						return &InternalError{Message: fmt.Sprintf("mkdir %s", destPath), Cause: err}
					}
					if err := copyDir(installPath, destPath); err != nil {
						return &InternalError{Message: fmt.Sprintf("copy %s", installPath), Cause: err}
					}
					copied = true
					break
				}
				if !copied {
					fmt.Fprintf(cmd.OutOrStdout(), "  !  %s@%s — no install path found, skipping\n", sl.Name, sl.Version)
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  ✓  %s@%s → %s\n", sl.Name, sl.Version, destPath)
				vendored++
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nVendored %d of %d skill(s) into %s\n", vendored, len(lf.Skills), vendorDir)
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Lockfile path (default: agent-skills.lock)")
	cmd.Flags().StringVar(&vendorDir, "dir", "", "Vendor directory (default: vendor/skills)")
	cmd.Flags().BoolVar(&clean, "clean", false, "Remove vendor dir before copying for a fresh state")
	return cmd
}

// vendorGitignoreLine returns the line to add to .gitignore for the vendor dir.
func vendorGitignoreLine(dir string) string {
	return strings.TrimPrefix(dir, "./") + "/"
}

var _ = vendorGitignoreLine
