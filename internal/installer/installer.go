package installer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/domehahn/skpm/v2/internal/archive"
	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/httpclient"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/registry"
	"golang.org/x/sync/errgroup"
)

// sharedHTTPClient is used for lockfile-URL downloads (httpRegistry) — see
// internal/httpclient's package doc for why this isn't http.DefaultClient.
var sharedHTTPClient = httpclient.New()

type Options struct {
	DryRun              bool
	Concurrency         int
	WorkDir             string
	LockfileDigest      string
	AdmissionDecisionID string
}

type Result struct {
	Installed []string
	FromCache []string
	Skipped   []string
	mu        sync.Mutex
}

func (r *Result) addInstalled(name string) {
	r.mu.Lock()
	r.Installed = append(r.Installed, name)
	r.mu.Unlock()
}

func (r *Result) addFromCache(name string) {
	r.mu.Lock()
	r.FromCache = append(r.FromCache, name)
	r.mu.Unlock()
}

func (r *Result) addSkipped(name string) {
	r.mu.Lock()
	r.Skipped = append(r.Skipped, name)
	r.mu.Unlock()
}

type Installer struct {
	cache *cache.Cache
}

func New(c *cache.Cache) *Installer {
	return &Installer{cache: c}
}

// Install reads a lockfile and installs all skills into the working directory.
func (ins *Installer) Install(ctx context.Context, lf *lockfile.LockFile, opts Options) (*Result, error) {
	if opts.WorkDir == "" {
		opts.WorkDir = "."
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = 4
	}

	result := &Result{}
	g, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, opts.Concurrency)

	for _, sl := range lf.Skills {
		sl := sl
		g.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()
			return ins.installOne(ctx, sl, opts, result)
		})
	}

	if err := g.Wait(); err != nil {
		return result, err
	}
	return result, nil
}

func (ins *Installer) installOne(ctx context.Context, sl lockfile.SkillLock, opts Options, result *Result) error {
	if opts.DryRun {
		result.addSkipped(sl.Name)
		return nil
	}

	zipPath, fromCache, err := ins.ensureCached(ctx, sl)
	if err != nil {
		return fmt.Errorf("skill %s: %w", sl.Name, err)
	}

	for _, dest := range sl.InstalledTo {
		absTarget := filepath.Join(opts.WorkDir, dest)
		if err := archive.AtomicExtract(zipPath, absTarget, ""); err != nil {
			return fmt.Errorf("skill %s: install to %s: %w", sl.Name, dest, err)
		}
		artDigest, _ := computeArtifactDigest(absTarget)
		pkgDigest := sl.SHA256
		if pkgDigest != "" && !strings.HasPrefix(pkgDigest, "sha256:") {
			pkgDigest = "sha256:" + pkgDigest
		}
		pkgRef := fmt.Sprintf("%s@%s", sl.Name, sl.Version)
		if sl.Namespace != "" && sl.Namespace != "default" {
			pkgRef = fmt.Sprintf("%s/%s@%s", sl.Namespace, sl.Name, sl.Version)
		}
		identity := MaterializedIdentity{
			SchemaVersion:       "1.0",
			Package:             pkgRef,
			PackageDigest:       pkgDigest,
			CompiledDigest:      artDigest,
			Registry:            sl.Source,
			InstalledPath:       absTarget,
			LockfileDigest:      opts.LockfileDigest,
			AdmissionDecisionID: opts.AdmissionDecisionID,
			Name:                sl.Name,
			Version:             sl.Version,
			ArtifactDigest:      artDigest,
			MaterializedPath:    absTarget,
			Verified:            true,
		}
		if idData, err := json.MarshalIndent(identity, "", "  "); err == nil {
			_ = os.WriteFile(filepath.Join(absTarget, ".skpm-installed.json"), idData, 0o644)
		}
	}

	if fromCache {
		result.addFromCache(sl.Name)
	} else {
		result.addInstalled(sl.Name)
	}
	return nil
}

type MaterializedIdentity struct {
	SchemaVersion       string `json:"schema_version"`
	Package             string `json:"package"`
	PackageDigest       string `json:"package_digest"`
	CompiledDigest      string `json:"compiled_digest"`
	Registry            string `json:"registry"`
	InstalledPath       string `json:"installed_path"`
	LockfileDigest      string `json:"lockfile_digest,omitempty"`
	AdmissionDecisionID string `json:"admission_decision_id,omitempty"`
	// Backwards compatibility fields
	Name             string `json:"name,omitempty"`
	Version          string `json:"version,omitempty"`
	ArtifactDigest   string `json:"artifact_digest,omitempty"`
	MaterializedPath string `json:"materialized_path,omitempty"`
	Verified         bool   `json:"verified"`
}

func computeArtifactDigest(dir string) (string, error) {
	var paths []string
	if err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		rel = filepath.ToSlash(rel)
		if rel == ".skpm-installed.json" {
			return nil
		}
		paths = append(paths, rel)
		return nil
	}); err != nil {
		return "", err
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, rel := range paths {
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			return "", err
		}
		_, _ = h.Write([]byte(rel))
		_, _ = h.Write(data)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

func (ins *Installer) ensureCached(ctx context.Context, sl lockfile.SkillLock) (string, bool, error) {
	if sl.SHA256 != "" && ins.cache.Has(sl.SHA256) {
		return ins.cache.Path(sl.SHA256), true, nil
	}

	reg := buildRegistryFromLock(sl)
	artifact := &registry.ResolvedArtifact{
		Name:        sl.Name,
		Version:     sl.Version,
		DownloadURL: sl.DownloadURL,
		SHA256:      sl.SHA256,
	}

	tmpKey := sl.SHA256
	if tmpKey == "" {
		tmpKey = "dl-" + sl.Name
	}
	tmp := ins.cache.Path(tmpKey) + ".part"
	if err := os.MkdirAll(filepath.Dir(tmp), 0o755); err != nil {
		return "", false, fmt.Errorf("create cache dir: %w", err)
	}

	f, err := os.Create(tmp)
	if err != nil {
		return "", false, fmt.Errorf("create download tmp: %w", err)
	}

	h := sha256.New()
	mw := io.MultiWriter(f, h)

	if err := reg.Download(ctx, artifact, mw); err != nil {
		f.Close()
		os.Remove(tmp)
		return "", false, fmt.Errorf("download: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return "", false, err
	}

	actualSHA := hex.EncodeToString(h.Sum(nil))
	if sl.SHA256 != "" && actualSHA != sl.SHA256 {
		os.Remove(tmp)
		return "", false, fmt.Errorf("SHA256 mismatch: expected %s, got %s", sl.SHA256, actualSHA)
	}

	finalPath := ins.cache.Path(actualSHA)
	if err := os.Rename(tmp, finalPath); err != nil {
		os.Remove(tmp)
		return "", false, fmt.Errorf("cache rename: %w", err)
	}
	return finalPath, false, nil
}

func buildRegistryFromLock(sl lockfile.SkillLock) registry.Registry {
	if sl.Source == "local" {
		return registry.NewLocalRegistry(filepath.Dir(sl.DownloadURL))
	}
	return &httpRegistry{}
}

// httpRegistry downloads directly from an artifact's DownloadURL via HTTP(S).
type httpRegistry struct{}

func (u *httpRegistry) Type() string { return "http" }

func (u *httpRegistry) Name() string { return "lockfile-url" }

func (u *httpRegistry) Capabilities(context.Context) (*registry.RegistryCapabilities, error) {
	return &registry.RegistryCapabilities{Download: true}, nil
}

func (u *httpRegistry) Resolve(_ context.Context, _ registry.ResolveRequest) (*registry.ResolvedArtifact, error) {
	return nil, fmt.Errorf("httpRegistry.Resolve not supported — URL is taken from lockfile")
}

func (u *httpRegistry) Download(ctx context.Context, artifact *registry.ResolvedArtifact, dest io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artifact.DownloadURL, nil)
	if err != nil {
		return err
	}
	resp, err := sharedHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("download %s: %w", artifact.DownloadURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", artifact.DownloadURL, resp.StatusCode)
	}
	_, err = httpclient.CopyLimited(dest, resp.Body)
	return err
}
