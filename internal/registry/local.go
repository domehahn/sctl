package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

type LocalRegistry struct {
	baseDir string
	name    string
}

// NewLocalRegistry creates a registry backed by the local filesystem.
// baseDir is the root directory containing skill subdirectories.
// Expected layout: baseDir/<skill-name>/<version>/<skill-name>-<version>.zip
func NewLocalRegistry(baseDir string) *LocalRegistry {
	return &LocalRegistry{baseDir: baseDir, name: "local"}
}

func (r *LocalRegistry) Type() string { return "local" }

func (r *LocalRegistry) Name() string { return r.name }

func (r *LocalRegistry) WithName(name string) *LocalRegistry {
	r.name = name
	return r
}

func (r *LocalRegistry) Capabilities(context.Context) (*RegistryCapabilities, error) {
	return &RegistryCapabilities{Resolve: true, Download: true, Search: true, Info: true, SemVerConstraints: true, Checksums: true}, nil
}

func (r *LocalRegistry) Resolve(ctx context.Context, req ResolveRequest) (*ResolvedArtifact, error) {
	name := req.Ref.Name
	version := req.Constraint
	skillDir := filepath.Join(r.baseDir, name)

	if version == "" || version == "latest" || strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") || strings.Contains(version, " ") {
		versions, err := r.ListVersions(ctx, req.Ref)
		if err != nil {
			return nil, err
		}
		version, err = SelectVersion(versions, req.Constraint, ConstraintOptions{
			IncludePrerelease: req.IncludePrerelease,
			AllowDeprecated:   req.AllowDeprecated,
			AllowYanked:       req.AllowYanked,
		})
		if err != nil {
			return nil, err
		}
	}

	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	path := filepath.Join(skillDir, version, assetName)

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("local: artifact not found at %s", path)
	}

	return &ResolvedArtifact{
		Namespace:    req.Ref.Namespace,
		Name:         name,
		Version:      version,
		Registry:     r.name,
		RegistryType: r.Type(),
		DownloadURL:  "file://" + path,
		ArtifactName: assetName,
		PackageType:  "zip",
	}, nil
}

func (r *LocalRegistry) ListVersions(_ context.Context, ref SkillRef) ([]VersionInfo, error) {
	skillDir := filepath.Join(r.baseDir, ref.Name)
	entries, err := os.ReadDir(skillDir)
	if err != nil {
		return nil, fmt.Errorf("local: read skill dir %s: %w", skillDir, err)
	}
	versions := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() && semver.IsValid("v"+e.Name()) {
			versions = append(versions, e.Name())
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare("v"+versions[i], "v"+versions[j]) > 0
	})
	out := make([]VersionInfo, 0, len(versions))
	for _, v := range versions {
		out = append(out, VersionInfo{Version: v})
	}
	return out, nil
}

func (r *LocalRegistry) Search(ctx context.Context, req SearchRequest) ([]SkillSearchResult, error) {
	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("local: read registry %s: %w", r.baseDir, err)
	}
	query := strings.ToLower(req.Query)
	var results []SkillSearchResult
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if query != "" && !strings.Contains(strings.ToLower(name), query) {
			continue
		}
		versions, _ := r.ListVersions(ctx, SkillRef{Namespace: req.Namespace, Name: name})
		result := SkillSearchResult{Namespace: req.Namespace, Name: name, Source: r.name}
		if len(versions) > 0 {
			result.Version = versions[0].Version
			result.LatestVersion = versions[0].Version
		}
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Name < results[j].Name })
	return results, nil
}

func (r *LocalRegistry) Info(ctx context.Context, ref SkillRef) (*SkillInfo, error) {
	versions, err := r.ListVersions(ctx, ref)
	if err != nil {
		return nil, err
	}
	info := &SkillInfo{Namespace: ref.Namespace, Name: ref.Name, Versions: versions, Source: r.name}
	if len(versions) > 0 {
		info.LatestVersion = versions[0].Version
	}
	return info, nil
}

func (r *LocalRegistry) resolveLatestVersion(skillDir, name string) (string, error) {
	entries, err := os.ReadDir(skillDir)
	if err != nil {
		return "", fmt.Errorf("local: read skill dir %s: %w", skillDir, err)
	}

	versions := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() && semver.IsValid("v"+e.Name()) {
			versions = append(versions, e.Name())
		}
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("local: no versions found for %s in %s", name, skillDir)
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare("v"+versions[i], "v"+versions[j]) > 0
	})
	return versions[0], nil
}

func (r *LocalRegistry) Download(_ context.Context, artifact *ResolvedArtifact, dest io.Writer) error {
	path := artifact.DownloadURL
	if len(path) > 7 && path[:7] == "file://" {
		path = path[7:]
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("local: open %s: %w", path, err)
	}
	defer f.Close()
	_, err = io.Copy(dest, f)
	return err
}
