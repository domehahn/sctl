package installer

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/domehahn/skpm/internal/cache"
	"github.com/domehahn/skpm/internal/lockfile"
	"github.com/domehahn/skpm/internal/registry"
	"golang.org/x/sync/errgroup"
)

type Options struct {
	DryRun      bool
	Concurrency int
	WorkDir     string
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
		if err := atomicUnzip(zipPath, absTarget); err != nil {
			return fmt.Errorf("skill %s: install to %s: %w", sl.Name, dest, err)
		}
	}

	if fromCache {
		result.addFromCache(sl.Name)
	} else {
		result.addInstalled(sl.Name)
	}
	return nil
}

func (ins *Installer) ensureCached(ctx context.Context, sl lockfile.SkillLock) (string, bool, error) {
	if sl.SHA256 != "" && ins.cache.Has(sl.SHA256) {
		return ins.cache.Path(sl.SHA256), true, nil
	}

	reg := buildRegistryFromLock(sl)
	artifact := &registry.ResolvedArtifact{
		Name:        sl.Name,
		Version:     sl.Version,
		DownloadURL: sl.SourceURL,
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

func atomicUnzip(zipPath, destDir string) error {
	staging := destDir + "~skpm-staging-" + strconv.Itoa(rand.Int())
	backup := destDir + "~skpm-backup-" + strconv.Itoa(rand.Int())

	if err := unzip(zipPath, staging); err != nil {
		os.RemoveAll(staging)
		return fmt.Errorf("unzip to staging: %w", err)
	}

	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, backup); err != nil {
			os.RemoveAll(staging)
			return fmt.Errorf("backup existing dir: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, destDir)
		return fmt.Errorf("create parent dir: %w", err)
	}

	if err := os.Rename(staging, destDir); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, destDir)
		return fmt.Errorf("rename staging to dest: %w", err)
	}

	os.RemoveAll(backup)
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		outPath := filepath.Join(dest, filepath.FromSlash(f.Name))
		if !isWithinDir(dest, outPath) {
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
		if err := extractFile(f, outPath); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(f *zip.File, dest string) error {
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

func isWithinDir(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	if rel == ".." {
		return false
	}
	if len(rel) >= 3 && rel[:3] == "../" {
		return false
	}
	return true
}

func buildRegistryFromLock(sl lockfile.SkillLock) registry.Registry {
	if sl.Source == "local" {
		return registry.NewLocalRegistry(filepath.Dir(sl.SourceURL))
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download %s: %w", artifact.DownloadURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", artifact.DownloadURL, resp.StatusCode)
	}
	_, err = io.Copy(dest, resp.Body)
	return err
}
