package cli

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
)

func newCloneCmd() *cobra.Command {
	var source string
	var dir string

	cmd := &cobra.Command{
		Use:   "clone <skill>[@version]",
		Short: "Download a skill from the registry into a local directory for editing",
		Long: `Downloads and extracts a skill's source ZIP into a local directory.

Useful for forking a skill or inspecting its contents before installing.
The version defaults to the latest available.

Examples:
  skpm clone my-skill              # latest version → ./my-skill/
  skpm clone my-skill@1.2.0        # specific version
  skpm clone my-skill --dir my-fork`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, err := parseSkillRef(args[0])
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			src := sourceOrDefault(source, cfg)
			if src == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}

			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Resolving %s", name)
			if version != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "@%s", version)
			}
			fmt.Fprintln(cmd.OutOrStdout(), " …")

			artifact, err := reg.Resolve(cmd.Context(), registry.ResolveRequest{
				Ref:        registry.ParseSkillRef(name, "default"),
				Constraint: version,
			})
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve: %v", err)}
			}

			destDir := dir
			if destDir == "" {
				destDir = artifact.Name
			}

			if _, err := os.Stat(destDir); err == nil {
				return &UserError{Message: fmt.Sprintf("destination %q already exists", destDir)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would clone %s@%s → %s\n", artifact.Name, artifact.Version, destDir)
				return nil
			}

			// Download into a temp file, then extract.
			tmp, err := os.CreateTemp("", "skpm-clone-*.zip")
			if err != nil {
				return &InternalError{Message: "create temp file", Cause: err}
			}
			defer os.Remove(tmp.Name())

			fmt.Fprintf(cmd.OutOrStdout(), "Downloading %s@%s …\n", artifact.Name, artifact.Version)
			if err := reg.Download(cmd.Context(), artifact, tmp); err != nil {
				tmp.Close()
				return &InternalError{Message: "download", Cause: err}
			}
			if err := tmp.Close(); err != nil {
				return &InternalError{Message: "flush download", Cause: err}
			}

			// Verify SHA256 if available.
			if artifact.SHA256 != "" {
				actual, err := sha256File(tmp.Name())
				if err != nil {
					return &InternalError{Message: "hash check", Cause: err}
				}
				if actual != artifact.SHA256 {
					return &InternalError{Message: fmt.Sprintf("SHA256 mismatch: expected %s, got %s", artifact.SHA256, actual)}
				}
			}

			if err := extractZipToDir(tmp.Name(), destDir); err != nil {
				os.RemoveAll(destDir)
				return &InternalError{Message: "extract", Cause: err}
			}

			// Cache the zip so subsequent installs are faster.
			if artifact.SHA256 != "" {
				cacheDir := cfg.CacheDir
				if cacheDir == "" {
					cacheDir, _ = defaultCacheDir()
				}
				c := cache.New(cacheDir)
				if !c.Has(artifact.SHA256) {
					_ = os.Rename(tmp.Name(), c.Path(artifact.SHA256))
				}
			}

			abs, _ := filepath.Abs(destDir)
			fmt.Fprintf(cmd.OutOrStdout(), "\nCloned %s@%s → %s\n", artifact.Name, artifact.Version, abs)
			fmt.Fprintf(cmd.OutOrStdout(), "To publish your fork: cd %s && skpm publish\n", destDir)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (uses default_registry if not set)")
	cmd.Flags().StringVar(&dir, "dir", "", "Target directory name (defaults to skill name)")
	return cmd
}

// parseSkillRef splits "name" or "name@version" into (name, version).
// Unlike parseSkillAtVersion, a missing @ is valid (means latest).
func parseSkillRef(raw string) (name, version string, err error) {
	at := lastAt(raw)
	if at < 0 {
		return raw, "", nil
	}
	n, v := raw[:at], raw[at+1:]
	if n == "" {
		return "", "", &UserError{Message: fmt.Sprintf("invalid skill reference %q", raw)}
	}
	return n, v, nil
}

func lastAt(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '@' {
			return i
		}
	}
	return -1
}

// extractZipToDir extracts a ZIP archive to destDir, stripping a single top-level
// directory if all entries share one (common convention for skill ZIPs).
func extractZipToDir(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Detect a common top-level prefix to strip.
	prefix := commonZipPrefix(r.File)

	for _, f := range r.File {
		rel := f.Name
		if prefix != "" && len(rel) > len(prefix) {
			rel = rel[len(prefix):]
		}
		if rel == "" {
			continue
		}

		outPath := filepath.Join(destDir, filepath.FromSlash(rel))
		if !isWithinBase(destDir, outPath) {
			return fmt.Errorf("zip slip detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		if err := extractZipEntry(f, outPath); err != nil {
			return err
		}
	}
	return nil
}

func commonZipPrefix(files []*zip.File) string {
	if len(files) == 0 {
		return ""
	}
	first := files[0].Name
	slash := len(first)
	for i := range first {
		if first[i] == '/' {
			slash = i + 1
			break
		}
	}
	prefix := first[:slash]
	for _, f := range files[1:] {
		if len(f.Name) < len(prefix) || f.Name[:len(prefix)] != prefix {
			return ""
		}
	}
	return prefix
}

func extractZipEntry(f *zip.File, dest string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func isWithinBase(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != ".." && (len(rel) < 3 || rel[:3] != "../")
}

func defaultCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "skpm"), nil
}
