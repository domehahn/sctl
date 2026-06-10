package cli

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newImportCmd() *cobra.Command {
	var platforms []string
	var target string
	var updateLock bool

	cmd := &cobra.Command{
		Use:   "import <zipfile>",
		Short: "Install a skill from a local ZIP file (air-gapped / CI artifact)",
		Long: `Extracts a skill ZIP artifact into the project without contacting a registry.

The skill name and platforms are read from skill.yaml inside the ZIP.
Use --platform to override the platform list.

The lockfile is updated with a "local" source entry unless --no-lock is passed.

Examples:
  skpm import dist/my-skill-1.2.0.zip
  skpm import artifact.zip --platform claude-code --platform gitlab-duo`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			zipPath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			if _, err := os.Stat(zipPath); err != nil {
				return &UserError{Message: fmt.Sprintf("file not found: %s", zipPath)}
			}

			sy, err := readSkillYAMLFromZip(zipPath)
			if err != nil {
				return err
			}
			if sy.Name == "" {
				return &UserError{Message: "skill.yaml inside the ZIP must declare a name"}
			}
			if sy.Version == "" {
				return &UserError{Message: "skill.yaml inside the ZIP must declare a version"}
			}

			plats := sy.CompatibleWith
			if len(platforms) > 0 {
				plats = make([]skill.Platform, len(platforms))
				for i, p := range platforms {
					plats[i] = skill.Platform(p)
				}
			}
			if len(plats) == 0 {
				plats = []skill.Platform{skill.PlatformAll}
			}

			installPaths, err := installer.ResolvePaths(sy.Name, plats)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve paths: %v", err)}
			}

			workDir, _ := os.Getwd()
			if target != "" {
				workDir = target
			}

			sha, _ := sha256File(zipPath)

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would import %s@%s from %s\n", sy.Name, sy.Version, zipPath)
				for _, p := range installPaths {
					fmt.Fprintf(cmd.OutOrStdout(), "  → %s\n", p)
				}
				return nil
			}

			for _, p := range installPaths {
				dest := filepath.Join(workDir, p)
				if err := atomicUnzipLocal(zipPath, dest); err != nil {
					return &InternalError{Message: fmt.Sprintf("install to %s", p), Cause: err}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "installed  %s\n", p)
			}

			if updateLock {
				lf, _ := lockfile.Read(lockfile.DefaultFilename)
				if lf == nil {
					lf = lockfile.New()
				}
				lf.Upsert(lockfile.SkillLock{
					Name:           sy.Name,
					Version:        sy.Version,
					Source:         "local",
					DownloadURL:    zipPath,
					SHA256:         sha,
					CompatibleWith: plats,
					InstalledTo:    installPaths,
				})
				if err := lf.Write(lockfile.DefaultFilename); err != nil {
					return &InternalError{Message: "update lockfile", Cause: err}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "updated  %s\n", lockfile.DefaultFilename)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nImported %s@%s (%d path(s))\n", sy.Name, sy.Version, len(installPaths))
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&platforms, "platform", nil, "Override compatible platforms")
	cmd.Flags().StringVar(&target, "target", "", "Install into a target directory (default: current dir)")
	cmd.Flags().BoolVar(&updateLock, "lock", true, "Update agent-skills.lock with a local source entry")
	return cmd
}

// readSkillYAMLFromZip extracts skill.yaml from a ZIP archive and unmarshals it.
func readSkillYAMLFromZip(zipPath string) (*skill.SkillYAML, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("open ZIP %s: %v", zipPath, err)}
	}
	defer r.Close()

	prefix := commonZipPrefix(r.File)

	for _, f := range r.File {
		name := f.Name
		if prefix != "" {
			name = strings.TrimPrefix(name, prefix)
		}
		if name != "skill.yaml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, &InternalError{Message: "open skill.yaml in ZIP", Cause: err}
		}
		defer rc.Close()

		var sy skill.SkillYAML
		if err := yaml.NewDecoder(rc).Decode(&sy); err != nil {
			return nil, &UserError{Message: fmt.Sprintf("parse skill.yaml in ZIP: %v", err)}
		}
		return &sy, nil
	}
	return nil, &UserError{Message: "ZIP does not contain skill.yaml"}
}

// atomicUnzipLocal is like atomicUnzip from installer, but reads from a local path
// without cache indirection.
func atomicUnzipLocal(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	prefix := commonZipPrefix(r.File)
	r.Close()

	staging := destDir + "~skpm-import-staging"
	backup := destDir + "~skpm-import-backup"

	if err := extractZipWithPrefix(zipPath, staging, prefix); err != nil {
		os.RemoveAll(staging)
		return fmt.Errorf("extract: %w", err)
	}

	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, backup); err != nil {
			os.RemoveAll(staging)
			return fmt.Errorf("backup: %w", err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, destDir)
		return err
	}
	if err := os.Rename(staging, destDir); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, destDir)
		return err
	}
	os.RemoveAll(backup)
	return nil
}

func extractZipWithPrefix(zipPath, destDir, prefix string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rel := f.Name
		if prefix != "" {
			rel = strings.TrimPrefix(rel, prefix)
		}
		if rel == "" {
			continue
		}
		outPath := filepath.Join(destDir, filepath.FromSlash(rel))
		if !isWithinBase(destDir, outPath) {
			return fmt.Errorf("zip slip: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0o755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0o755)
		if err := extractZipEntry(f, outPath); err != nil {
			return err
		}
	}
	return nil
}
