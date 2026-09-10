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

	"github.com/domehahn/skpm/v2/internal/archive"
	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newAddCmd() *cobra.Command {
	var (
		source       string
		ref          string
		skillSubPath string
		doInstall    bool
		noInstall    bool
		doLock       bool
		noLock       bool
	)

	cmd := &cobra.Command{
		Use:   "add <skill[@version]>",
		Short: "Add a skill and write agent-skills.lock",
		Long: `Resolves the skill version from the configured registry, downloads and
verifies the artifact, installs it to all compatible platform paths, and
writes or updates agent-skills.lock.

Local paths are supported directly without a registry:
  skpm add ./my-skill
  skpm add ../shared-skills/sdlc-manager

For skills without a release tag, use --ref to download from a branch or commit:
  skpm add my-skill --source myregistry --ref main
  skpm add my-skill --source myregistry --ref feature/new-checks --path skills/my-skill`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			// Local path: skpm add ./my-skill  or  skpm add /abs/path/to/skill
			if isLocalPath(args[0]) {
				return addFromLocalPath(cmd, args[0], format, addOptions{Install: doInstall && !noInstall, Lock: doLock && !noLock})
			}

			// Ref-based download: no release tag needed
			if ref != "" {
				return addFromRef(cmd, args[0], ref, skillSubPath, source, format, addOptions{Install: doInstall && !noInstall, Lock: doLock && !noLock})
			}

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

			artifact, err := reg.Resolve(cmd.Context(), registry.ResolveRequest{
				Ref:        registry.ParseSkillRef(name, ""),
				Constraint: version,
			})
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve %s@%s: %v", name, version, err)}
			}

			c := cache.New(cfg.CacheDir)
			zipPath, actualSHA, err := downloadAndVerify(cmd.Context(), reg, artifact, c)
			if err != nil {
				if strings.Contains(err.Error(), "SHA256 mismatch") {
					return &IntegrityError{Message: fmt.Sprintf("integrity failure during download: %v", err)}
				}
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

			shouldInstall := doInstall && !noInstall
			shouldLock := doLock && !noLock
			if shouldInstall {
				workDir, _ := os.Getwd()
				for _, dest := range installPaths {
					absTarget := filepath.Join(workDir, dest)
					if err := archive.AtomicExtract(zipPath, absTarget, ""); err != nil {
						return &InternalError{Message: fmt.Sprintf("install to %s", dest), Cause: err}
					}
				}
			}

			// update lockfile
			if shouldLock {
				lf, _ := lockfile.Read(lockfile.DefaultFilename)
				if lf == nil {
					lf = lockfile.New()
				}
				lf.Upsert(lockfile.SkillLock{
					Name:           name,
					Namespace:      artifact.Namespace,
					Version:        artifact.Version,
					Source:         src,
					RegistryType:   artifact.RegistryType,
					DownloadURL:    artifact.DownloadURL,
					Artifact:       artifact.Artifact,
					SHA256:         actualSHA,
					PackageType:    artifact.PackageType,
					CompatibleWith: compatibleWith,
					InstalledTo:    installPaths,
					Metadata:       artifact.Metadata,
				})
				if err := lf.Write(lockfile.DefaultFilename); err != nil {
					return &InternalError{Message: "write lockfile", Cause: err}
				}
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

			_ = runProjectHook(cmd, "post_add")

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
	cmd.Flags().StringVar(&ref, "ref", "", "Git branch, tag, or commit SHA to download from (no release needed)")
	cmd.Flags().StringVar(&skillSubPath, "path", "", "Path of the skill within the repository (e.g. skills/my-skill)")
	cmd.Flags().BoolVar(&doInstall, "install", true, "Install after adding")
	cmd.Flags().BoolVar(&noInstall, "no-install", false, "Do not install after adding")
	cmd.Flags().BoolVar(&doLock, "lock", true, "Update agent-skills.lock after adding")
	cmd.Flags().BoolVar(&noLock, "no-lock", false, "Do not update agent-skills.lock after adding")
	return cmd
}

type addOptions struct {
	Install bool
	Lock    bool
}

// addFromRef downloads a skill from a specific git ref without requiring a release.
func addFromRef(cmd *cobra.Command, skillName, ref, skillSubPath, source string, format OutputFormat, opts addOptions) error {
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

	log.Debug().Str("skill", skillName).Str("ref", ref).Str("path", skillSubPath).Msg("downloading ref")
	fmt.Fprintf(cmd.OutOrStdout(), "Downloading %s@%s from %s\n", skillName, ref, src)

	tmpDir, err := os.MkdirTemp("", "skpm-ref-*")
	if err != nil {
		return &InternalError{Message: "create temp dir", Cause: err}
	}
	defer os.RemoveAll(tmpDir)

	skillDir := filepath.Join(tmpDir, skillName)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return &InternalError{Message: "create skill dir", Cause: err}
	}

	if err := registry.DownloadRef(cmd.Context(), reg, skillName, ref, skillSubPath, skillDir); err != nil {
		return &UserError{Message: fmt.Sprintf("download ref: %v", err)}
	}

	// Validate what we got.
	v := skill.NewValidator()
	res, err := v.Validate(cmd.Context(), skillDir)
	if err != nil {
		return &InternalError{Message: "validate", Cause: err}
	}
	if !res.Valid {
		msgs := make([]string, len(res.Errors))
		for i, e := range res.Errors {
			msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
		}
		return &UserError{Message: fmt.Sprintf("downloaded skill is invalid:\n  %s", strings.Join(msgs, "\n  "))}
	}

	// Install from the temp directory using the existing local path logic.
	return addFromLocalPath(cmd, skillDir, format, opts)
}

func isLocalPath(arg string) bool {
	return strings.HasPrefix(arg, "./") ||
		strings.HasPrefix(arg, "../") ||
		strings.HasPrefix(arg, "/") ||
		arg == "." || arg == ".."
}

func addFromLocalPath(cmd *cobra.Command, path string, format OutputFormat, opts addOptions) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return &UserError{Message: fmt.Sprintf("resolve path %q: %v", path, err)}
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return &UserError{Message: fmt.Sprintf("path not found: %s", absPath)}
	}

	// validate first
	v := skill.NewValidator()
	res, err := v.Validate(cmd.Context(), absPath)
	if err != nil {
		return &InternalError{Message: "validate", Cause: err}
	}
	if !res.Valid {
		msgs := make([]string, len(res.Errors))
		for i, e := range res.Errors {
			msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
		}
		return &UserError{Message: fmt.Sprintf("skill at %s is invalid:\n  %s", absPath, strings.Join(msgs, "\n  "))}
	}

	sy, err := readLocalSkillYAML(absPath)
	if err != nil {
		return &UserError{Message: fmt.Sprintf("read skill.yaml: %v", err)}
	}

	installPaths, err := installer.ResolvePaths(sy.Name, sy.CompatibleWith)
	if err != nil {
		return &UserError{Message: fmt.Sprintf("resolve platform paths: %v", err)}
	}

	if opts.Install {
		workDir, _ := os.Getwd()
		for _, dest := range installPaths {
			if err := copyDir(absPath, filepath.Join(workDir, dest)); err != nil {
				return &InternalError{Message: fmt.Sprintf("install to %s", dest), Cause: err}
			}
		}
	}

	sourceURL := "file://" + absPath

	if opts.Lock {
		lf, _ := lockfile.Read(lockfile.DefaultFilename)
		if lf == nil {
			lf = lockfile.New()
		}
		lf.Upsert(lockfile.SkillLock{
			Name:           sy.Name,
			Namespace:      "default",
			Version:        sy.Version,
			Source:         "local",
			RegistryType:   "local",
			DownloadURL:    sourceURL,
			PackageType:    "directory",
			CompatibleWith: sy.CompatibleWith,
			InstalledTo:    installPaths,
		})
		if err := lf.Write(lockfile.DefaultFilename); err != nil {
			return &InternalError{Message: "write lockfile", Cause: err}
		}
	}

	mf, _ := manifest.Read(manifest.DefaultFilename)
	if mf == nil {
		mf = manifest.New()
	}
	mf.Upsert(manifest.SkillEntry{
		Name:    sy.Name,
		Version: sy.Version,
		Source:  sourceURL,
	})
	if err := mf.Write(manifest.DefaultFilename); err != nil {
		return &InternalError{Message: "write manifest", Cause: err}
	}

	if format == OutputJSON {
		PrintResult(format, CommandResult{
			Success: true,
			Command: "add",
			Data: map[string]interface{}{
				"name":         sy.Name,
				"version":      sy.Version,
				"source":       "local",
				"path":         absPath,
				"installed_to": installPaths,
			},
		})
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Added %s@%s (local)\n", sy.Name, sy.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "  Path:         %s\n", absPath)
	fmt.Fprintf(cmd.OutOrStdout(), "  Installed to:\n")
	for _, p := range installPaths {
		fmt.Fprintf(cmd.OutOrStdout(), "    - %s\n", p)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "  Lockfile:     %s\n", lockfile.DefaultFilename)
	fmt.Fprintf(cmd.OutOrStdout(), "\n  Note: local source — not reproducible on other machines.\n")
	return nil
}

func readLocalSkillYAML(dir string) (*skill.SkillYAML, error) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		return nil, err
	}
	var sy skill.SkillYAML
	if err := yaml.Unmarshal(data, &sy); err != nil {
		return nil, err
	}
	return &sy, nil
}

// copyDir copies a directory tree from src to dst atomically via a staging dir.
func copyDir(src, dst string) error {
	staging := dst + "~skpm-staging"
	backup := dst + "~skpm-backup"

	if err := copyDirContents(src, staging); err != nil {
		os.RemoveAll(staging)
		return err
	}
	if _, err := os.Stat(dst); err == nil {
		if err := os.Rename(dst, backup); err != nil {
			os.RemoveAll(staging)
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, dst)
		return err
	}
	if err := os.Rename(staging, dst); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, dst)
		return err
	}
	os.RemoveAll(backup)
	return nil
}

func copyDirContents(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
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

	tmp := filepath.Join(os.TempDir(), "skpm-download-*.zip")
	f, err := os.CreateTemp("", "skpm-download-*.zip")
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
