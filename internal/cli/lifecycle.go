package cli

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"time"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/domehahn/skpm/v2/internal/skill"
	"golang.org/x/mod/semver"
)

func lockFromManifest(ctx context.Context, mf *manifest.ManifestFile, cfg *config.Config, lockPath string) (*lockfile.LockFile, error) {
	lf := lockfile.New()
	lf.ResolvedAt = time.Now().UTC().Format(time.RFC3339)
	lf.GeneratedBy = "skpm"

	for _, entry := range mf.Skills {
		src := entry.Source
		if src == "" {
			src = cfg.DefaultRegistry
		}
		if src == "" {
			return nil, &UserError{Message: fmt.Sprintf("skill %q has no source and no default_registry is configured", entry.Name)}
		}
		reg, err := registry.New(src, cfg)
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("registry for %q: %v", entry.Name, err)}
		}
		ref := registry.ParseSkillRef(entry.Name, "")
		artifact, err := reg.Resolve(ctx, registry.ResolveRequest{Ref: ref, Constraint: entry.Version})
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("resolve %s: %v", entry.Name, err)}
		}

		platforms := artifact.CompatibleWith
		installPaths := []string(nil)
		if len(platforms) > 0 {
			paths, err := installer.ResolvePaths(entry.Name, stringsToPlatforms(platforms))
			if err != nil {
				return nil, &UserError{Message: fmt.Sprintf("resolve paths for %s: %v", entry.Name, err)}
			}
			installPaths = paths
		}

		lf.Upsert(lockfile.SkillLock{
			Name:           entry.Name,
			Namespace:      artifact.Namespace,
			Version:        artifact.Version,
			Source:         src,
			RegistryType:   artifact.RegistryType,
			DownloadURL:    artifact.DownloadURL,
			Artifact:       artifact.ArtifactName,
			SHA256:         artifact.SHA256,
			PackageType:    artifact.PackageType,
			CompatibleWith: stringsToPlatforms(platforms),
			InstalledTo:    installPaths,
			Metadata:       artifact.Metadata,
		})
	}
	lf.Sort()
	return lf, nil
}

func lockfileOutdated(ctx context.Context, cfg *config.Config, manifestPath, lockPath string) (bool, error) {
	mf, err := manifest.Read(manifestPath)
	if err != nil {
		return false, err
	}
	resolved, err := lockFromManifest(ctx, mf, cfg, lockPath)
	if err != nil {
		return false, err
	}
	current, err := lockfile.Read(lockPath)
	if err != nil {
		return true, nil
	}
	return !sameLockedSkills(resolved, current), nil
}

func sameLockedSkills(a, b *lockfile.LockFile) bool {
	aa := append([]lockfile.SkillLock(nil), a.Skills...)
	bb := append([]lockfile.SkillLock(nil), b.Skills...)
	sort.Slice(aa, func(i, j int) bool { return aa[i].Name < aa[j].Name })
	sort.Slice(bb, func(i, j int) bool { return bb[i].Name < bb[j].Name })
	for i := range aa {
		aa[i].InstalledTo = nil
		bb[i].InstalledTo = nil
	}
	return reflect.DeepEqual(aa, bb)
}


func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func registryDiscovery(ctx context.Context, source string, cfg *config.Config) (registry.Registry, registry.DiscoveryRegistry, error) {
	reg, err := registry.New(source, cfg)
	if err != nil {
		return nil, nil, err
	}
	discovery, ok := reg.(registry.DiscoveryRegistry)
	if !ok {
		return nil, nil, fmt.Errorf("registry %q does not support discovery", source)
	}
	return reg, discovery, nil
}

func registryDefaultNamespace(_ registry.Registry) string {
	return "default"
}

func latestVersion(versions []registry.VersionInfo) string {
	latest := ""
	for _, v := range versions {
		if v.Yanked {
			continue
		}
		if latest == "" || semver.Compare("v"+v.Version, "v"+latest) > 0 {
			latest = v.Version
		}
	}
	return latest
}

func platformsToStrings(platforms []skill.Platform) []string {
	out := make([]string, 0, len(platforms))
	for _, p := range platforms {
		out = append(out, string(p))
	}
	return out
}

func stringsToPlatforms(platforms []string) []skill.Platform {
	out := make([]skill.Platform, 0, len(platforms))
	for _, p := range platforms {
		out = append(out, skill.Platform(p))
	}
	return out
}

func readZipFile(path, name string) ([]byte, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, rc); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, os.ErrNotExist
}

func pruneInstallPaths(workDir string, lf *lockfile.LockFile) error {
	keep := map[string]bool{}
	for _, sl := range lf.Skills {
		for _, p := range sl.InstalledTo {
			keep[filepath.Clean(p)] = true
		}
	}
	roots := []string{".claude/skills", ".github/skills", ".agents/skills", "skills"}
	for _, root := range roots {
		entries, err := os.ReadDir(filepath.Join(workDir, root))
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			rel := filepath.Clean(filepath.Join(root, e.Name()))
			if !keep[rel] {
				if err := os.RemoveAll(filepath.Join(workDir, rel)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
