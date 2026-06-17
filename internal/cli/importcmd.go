package cli

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newImportSkillsCmd() *cobra.Command {
	var destDir string
	var force bool

	cmd := &cobra.Command{
		Use:   "unpack <tarball>",
		Short: "Unpack a skill bundle tarball into the skill install directory",
		Long: `Reads a .tar.gz bundle produced by 'skcr bundle' and extracts each skill
directory into --dir (defaults to .agents/skills/).

Existing files are skipped unless --force is passed.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tarball := args[0]

			if destDir == "" {
				destDir = ".agents/skills"
			}

			f, err := os.Open(tarball)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("open tarball: %v", err)}
			}
			defer f.Close()

			gz, err := gzip.NewReader(f)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("open gzip stream: %v", err)}
			}
			defer gz.Close()

			tr := tar.NewReader(gz)
			created, skipped := 0, 0

			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return &InternalError{Message: "read tarball", Cause: err}
				}
				if hdr.Typeflag == tar.TypeDir {
					continue
				}

				// Security: reject paths that escape destDir.
				clean := filepath.Join(destDir, filepath.FromSlash(hdr.Name))
				if !strings.HasPrefix(filepath.Clean(clean)+string(filepath.Separator), filepath.Clean(destDir)+string(filepath.Separator)) {
					return &UserError{Message: fmt.Sprintf("tarball entry %q would escape destination directory", hdr.Name)}
				}

				dest := clean

				if _, err := os.Stat(dest); err == nil && !force {
					fmt.Fprintf(cmd.OutOrStdout(), "  skip  %s (already exists)\n", hdr.Name)
					skipped++
					continue
				}

				if globalDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "  would create  %s\n", hdr.Name)
					created++
					continue
				}

				if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
					return &InternalError{Message: "mkdir", Cause: err}
				}
				out, err := os.Create(dest)
				if err != nil {
					return &InternalError{Message: "create file", Cause: err}
				}
				if _, err := io.Copy(out, tr); err != nil {
					out.Close()
					return &InternalError{Message: "write file", Cause: err}
				}
				out.Close()
				fmt.Fprintf(cmd.OutOrStdout(), "  +  %s\n", hdr.Name)
				created++
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: would create %d file(s), skip %d.\n", created, skipped)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nImported %d file(s) into %s, skipped %d.\n", created, destDir, skipped)
			return nil
		},
	}

	cmd.Flags().StringVar(&destDir, "dir", "", "Destination directory (default: .agents/skills)")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	return cmd
}
