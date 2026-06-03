package cache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Cache stores downloaded artifacts keyed by their SHA256 hash.
// Multiple concurrent processes writing the same key are safe because writes
// are idempotent — both write identical bytes to the same destination path.
type Cache struct {
	dir string
}

func New(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) Has(sha256 string) bool {
	_, err := os.Stat(c.path(sha256))
	return err == nil
}

func (c *Cache) Get(sha256 string) (io.ReadCloser, error) {
	f, err := os.Open(c.path(sha256))
	if err != nil {
		return nil, fmt.Errorf("cache get %s: %w", sha256, err)
	}
	return f, nil
}

// Put streams r into the cache under the given sha256 key.
// The write is atomic: data lands in a .part file first, then renamed.
func (c *Cache) Put(sha256 string, r io.Reader) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return fmt.Errorf("cache mkdir: %w", err)
	}
	tmp := c.path(sha256) + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("cache create part: %w", err)
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("cache write: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("cache close: %w", err)
	}
	if err := os.Rename(tmp, c.path(sha256)); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("cache rename: %w", err)
	}
	return nil
}

// Path returns the absolute filesystem path for a cached artifact.
func (c *Cache) Path(sha256 string) string {
	return c.path(sha256)
}

func (c *Cache) path(sha256 string) string {
	return filepath.Join(c.dir, sha256+".zip")
}
