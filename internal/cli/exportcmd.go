package cli

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	var lockPath string
	var outFile string
	var includeSnapshots bool

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Pack installed skills into a portable .tar.gz archive",
		Long: `Reads the lockfile and bundles every installed skill directory into a
.tar.gz archive that can be extracted in an air-gapped environment or shared
without registry access.

The archive includes:
  - agent-skills.lock  (the lockfile itself)
  - each skill's InstalledTo directory tree`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			if outFile == "" {
				outFile = "skills-bundle.tar.gz"
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would write %s (%d skill(s))\n", outFile, len(lf.Skills))
				return nil
			}

			f, err := os.Create(outFile)
			if err != nil {
				return &InternalError{Message: "create archive", Cause: err}
			}
			defer f.Close()

			gw := gzip.NewWriter(f)
			tw := tar.NewWriter(gw)

			// Include the lockfile itself.
			if err := addFileToTar(tw, lockPath, "agent-skills.lock"); err != nil {
				return &InternalError{Message: "add lockfile to archive", Cause: err}
			}

			skillCount := 0
			for _, sl := range lf.Skills {
				for _, installPath := range sl.InstalledTo {
					info, statErr := os.Stat(installPath)
					if statErr != nil || !info.IsDir() {
						continue
					}
					if err := addDirToTar(tw, installPath, "skills/"+filepath.Base(installPath)); err != nil {
						return &InternalError{Message: fmt.Sprintf("archive skill dir %s", installPath), Cause: err}
					}
					skillCount++
					break // one install path per skill is enough
				}
			}

			if includeSnapshots {
				if _, err := os.Stat(snapshotDir); err == nil {
					if err := addDirToTar(tw, snapshotDir, snapshotDir); err != nil {
						return &InternalError{Message: "archive snapshots", Cause: err}
					}
				}
			}

			if err := tw.Close(); err != nil {
				return &InternalError{Message: "finalize tar", Cause: err}
			}
			if err := gw.Close(); err != nil {
				return &InternalError{Message: "finalize gzip", Cause: err}
			}

			fi, _ := f.Stat()
			fmt.Fprintf(cmd.OutOrStdout(), "Exported %d skill(s) → %s (%s)\n", skillCount, outFile, humanBytes(fi.Size()))
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Lockfile path (default: agent-skills.lock)")
	cmd.Flags().StringVarP(&outFile, "out", "f", "", "Output archive path (default: skills-bundle.tar.gz)")
	cmd.Flags().BoolVar(&includeSnapshots, "include-snapshots", false, "Also include .skpm-snapshots/ in the archive")
	return cmd
}

func addFileToTar(tw *tar.Writer, srcPath, tarPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	hdr := &tar.Header{
		Name:    tarPath,
		Mode:    int64(info.Mode()),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Format:  tar.FormatGNU,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err = io.Copy(tw, f)
	return err
}

func addDirToTar(tw *tar.Writer, srcDir, tarBase string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		tarPath := filepath.Join(tarBase, rel)
		tarPath = strings.ReplaceAll(tarPath, string(filepath.Separator), "/")

		if info.IsDir() {
			return tw.WriteHeader(&tar.Header{
				Typeflag: tar.TypeDir,
				Name:     tarPath + "/",
				Mode:     int64(info.Mode()),
				ModTime:  info.ModTime(),
				Format:   tar.FormatGNU,
			})
		}

		hdr := &tar.Header{
			Name:    tarPath,
			Mode:    int64(info.Mode()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Format:  tar.FormatGNU,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
}

// bundleTimestamp returns a sortable timestamp suffix for archive names.
func bundleTimestamp() string {
	return time.Now().UTC().Format("20060102T150405Z")
}

var _ = bundleTimestamp // suppress unused warning if not called externally
