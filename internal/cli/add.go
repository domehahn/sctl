package cli

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/sctl/internal/cache"
	"github.com/domehahn/sctl/internal/config"
	"github.com/domehahn/sctl/internal/installer"
	"github.com/domehahn/sctl/internal/lockfile"
	"github.com/domehahn/sctl/internal/manifest"
	"github.com/domehahn/sctl/internal/registry"
	"github.com/domehahn/sctl/internal/skill"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newAddCmd() *cobra.Command {
	var source string

	cmd := &cobra.Command{
		Use:   "add <skill[@version]>",
		Short: "Add a skill and write agent-skills.lock",
		Long: `Resolves the skill version from the configured registry, downloads and
verifies the artifact, installs it to all compatible platform paths, and
writes or updates agent-skills.lock.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			name, version := parseSkillArg(args[0])

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			src := source
			if src == "" {
				src = cfg.DefaultRegistry
			}
			if src == "" {
				return &UserError{Message: "no registry specified — use --source or set default_registry in config"}
			}

			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}

			log.Debug().Str("skill", name).Str("version", version).Str("source", src).Msg("resolving")

			artifact, err := reg.Resolve(cmd.Context(), name, version)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve %s@%s: %v", name, version, err)}
			}

			c := cache.New(cfg.CacheDir)
			zipPath, actualSHA, err := downloadAndVerify(cmd.Context(), reg, artifact, c)
			if err != nil {
				return &InternalError{Message: "download", Cause: err}
			}

			compatibleWith, err := readCompatibleWith(zipPath)
			if err != nil {
				log.Warn().Err(err).Msg("could not read compatible_with from skill.yaml — defaulting to all platforms")
				compatibleWith = []skill.Platform{skill.PlatformAll}
			}

			installPaths, err := installer.ResolvePaths(name, compatibleWith)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve platform paths: %v", err)}
			}

			workDir, _ := os.Getwd()
			for _, dest := range installPaths {
				absTarget := filepath.Join(workDir, dest)
				if err := atomicUnzipPublic(zipPath, absTarget); err != nil {
					return &InternalError{Message: fmt.Sprintf("install to %s", dest), Cause: err}
				}
			}

			// update lockfile
			lf, _ := lockfile.Read(lockfile.DefaultFilename)
			if lf == nil {
				lf = lockfile.New()
			}
			lf.Upsert(lockfile.SkillLock{
				Name:        name,
				Version:     artifact.Version,
				Source:      src,
				SourceURL:   artifact.DownloadURL,
				SHA256:      actualSHA,
				InstalledTo: installPaths,
			})
			if err := lf.Write(lockfile.DefaultFilename); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}

			// update manifest (source of truth for re-generating the lockfile)
			mf, _ := manifest.Read(manifest.DefaultFilename)
			if mf == nil {
				mf = manifest.New()
			}
			mf.Upsert(manifest.SkillEntry{
				Name:    name,
				Version: artifact.Version,
				Source:  src,
			})
			if err := mf.Write(manifest.DefaultFilename); err != nil {
				return &InternalError{Message: "write manifest", Cause: err}
			}

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "add",
					Data: map[string]interface{}{
						"name":         name,
						"version":      artifact.Version,
						"sha256":       actualSHA,
						"installed_to": installPaths,
					},
				})
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Added %s@%s\n", name, artifact.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  SHA256:       %s\n", actualSHA)
			fmt.Fprintf(cmd.OutOrStdout(), "  Installed to:\n")
			for _, p := range installPaths {
				fmt.Fprintf(cmd.OutOrStdout(), "    - %s\n", p)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  Lockfile:     %s\n", lockfile.DefaultFilename)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (named config key, github, gitlab, artifactory, local, or path)")
	return cmd
}

func parseSkillArg(arg string) (name, version string) {
	parts := strings.SplitN(arg, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return arg, ""
}

func downloadAndVerify(ctx context.Context, reg registry.Registry, artifact *registry.ResolvedArtifact, c *cache.Cache) (string, string, error) {
	if artifact.SHA256 != "" && c.Has(artifact.SHA256) {
		return c.Path(artifact.SHA256), artifact.SHA256, nil
	}

	tmp := filepath.Join(os.TempDir(), "skm-download-*.zip")
	f, err := os.CreateTemp("", "skm-download-*.zip")
	if err != nil {
		return "", "", fmt.Errorf("create temp: %w", err)
	}
	defer func() {
		f.Close()
		if artifact.SHA256 == "" {
			os.Remove(tmp)
		}
	}()

	h := sha256.New()
	mw := io.MultiWriter(f, h)

	if err := reg.Download(ctx, artifact, mw); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", "", fmt.Errorf("download: %w", err)
	}
	f.Close()

	actualSHA := hex.EncodeToString(h.Sum(nil))
	if artifact.SHA256 != "" && actualSHA != artifact.SHA256 {
		os.Remove(f.Name())
		return "", "", fmt.Errorf("SHA256 mismatch: expected %s, got %s", artifact.SHA256, actualSHA)
	}

	if err := c.Put(actualSHA, mustOpen(f.Name())); err != nil {
		_ = err
	}

	return f.Name(), actualSHA, nil
}

func mustOpen(path string) io.Reader {
	f, _ := os.Open(path)
	return f
}

func readCompatibleWith(zipPath string) ([]skill.Platform, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "skill.yaml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			var sy skill.SkillYAML
			if err := yaml.NewDecoder(rc).Decode(&sy); err != nil {
				return nil, err
			}
			return sy.CompatibleWith, nil
		}
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			var m struct {
				CompatibleWith []skill.Platform `json:"compatible_with"`
			}
			if err := json.NewDecoder(rc).Decode(&m); err != nil {
				return nil, err
			}
			return m.CompatibleWith, nil
		}
	}
	return nil, fmt.Errorf("skill.yaml not found in ZIP")
}

// atomicUnzipPublic delegates to the installer package's internal function
// by re-using the same logic via the public Install path.
func atomicUnzipPublic(zipPath, destDir string) error {
	stagingDir := destDir + "~skm-stage"
	backupDir := destDir + "~skm-bak"

	if err := unzipDir(zipPath, stagingDir); err != nil {
		os.RemoveAll(stagingDir)
		return err
	}
	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, backupDir); err != nil {
			os.RemoveAll(stagingDir)
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		os.RemoveAll(stagingDir)
		os.Rename(backupDir, destDir)
		return err
	}
	if err := os.Rename(stagingDir, destDir); err != nil {
		os.RemoveAll(stagingDir)
		os.Rename(backupDir, destDir)
		return err
	}
	os.RemoveAll(backupDir)
	return nil
}

func unzipDir(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		outPath := filepath.Join(dest, filepath.FromSlash(f.Name))
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0o755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0o755)
		dst, err := os.Create(outPath)
		if err != nil {
			return err
		}
		src, err := f.Open()
		if err != nil {
			dst.Close()
			return err
		}
		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
