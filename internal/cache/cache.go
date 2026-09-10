package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Cache stores downloaded artifacts keyed by their SHA256 hash or key.
// All access is verified via digest before reuse and purged on corruption.
type Cache struct {
	dir string
	mu  sync.Mutex
}

func New(dir string) *Cache {
	return &Cache{dir: dir}
}

func normalizeDigestKey(key string) string {
	key = strings.TrimPrefix(key, "sha256:")
	return strings.ToLower(key)
}

func (c *Cache) Has(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	normKey := normalizeDigestKey(key)
	if normKey == "" {
		return false
	}
	filePath := c.path(normKey)
	if _, err := os.Stat(filePath); err != nil {
		return false
	}

	if len(normKey) == 64 {
		actual, err := fileSHA256(filePath)
		if err != nil || actual != normKey {
			_ = os.Remove(filePath)
			return false
		}
	}
	return true
}

func (c *Cache) Get(key string) (io.ReadCloser, error) {
	if !c.Has(key) {
		return nil, fmt.Errorf("cache miss or corrupted entry: %s", key)
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	filePath := c.path(normalizeDigestKey(key))
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cache get %s: %w", key, err)
	}
	return f, nil
}

// Put streams r into the cache under the given key.
// The write is atomic: data lands in a .part file first, then renamed.
func (c *Cache) Put(key string, r io.Reader) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	normKey := normalizeDigestKey(key)
	if normKey == "" {
		return fmt.Errorf("cache put: empty key")
	}

	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return fmt.Errorf("cache mkdir: %w", err)
	}
	tmp := filepath.Join(c.dir, "part-"+normKey+".part")
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("cache create part: %w", err)
	}

	h := sha256.New()
	mw := io.MultiWriter(f, h)

	if _, err := io.Copy(mw, r); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("cache write: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("cache close: %w", err)
	}

	actualSHA := hex.EncodeToString(h.Sum(nil))
	if len(normKey) == 64 && actualSHA != normKey {
		os.Remove(tmp)
		return fmt.Errorf("cache digest mismatch: expected %s, got %s", normKey, actualSHA)
	}

	finalPath := c.path(normKey)
	if err := os.Rename(tmp, finalPath); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("cache rename: %w", err)
	}
	return nil
}

// Path returns the absolute filesystem path for a cached artifact.
func (c *Cache) Path(key string) string {
	return c.path(normalizeDigestKey(key))
}

func (c *Cache) path(key string) string {
	return filepath.Join(c.dir, key+".zip")
}

func fileSHA256(path string) (string, error) {
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
