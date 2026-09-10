package cache_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSHA = "abc123def456"

func TestHasMiss(t *testing.T) {
	c := cache.New(t.TempDir())
	assert.False(t, c.Has(testSHA))
}

func TestPutAndGet(t *testing.T) {
	c := cache.New(t.TempDir())

	require.NoError(t, c.Put(testSHA, strings.NewReader("hello cache")))
	assert.True(t, c.Has(testSHA))

	rc, err := c.Get(testSHA)
	require.NoError(t, err)
	defer rc.Close()

	data, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, "hello cache", string(data))
}

func TestPutIdempotent(t *testing.T) {
	c := cache.New(t.TempDir())

	require.NoError(t, c.Put(testSHA, strings.NewReader("data")))
	require.NoError(t, c.Put(testSHA, strings.NewReader("data")))
	assert.True(t, c.Has(testSHA))
}

func TestPath(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir)
	p := c.Path(testSHA)
	assert.Contains(t, p, testSHA)
	assert.Contains(t, p, dir)
}

func TestGetMissing(t *testing.T) {
	c := cache.New(t.TempDir())
	_, err := c.Get("nonexistent")
	assert.Error(t, err)
}

func TestPutReadError(t *testing.T) {
	c := cache.New(t.TempDir())
	err := c.Put(testSHA, &errorReader{})
	assert.Error(t, err)
	assert.False(t, c.Has(testSHA))
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestCorruptedCacheAutoPurge(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir)
	validSHA := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" // sha256("hello")

	require.NoError(t, c.Put(validSHA, strings.NewReader("hello")))
	assert.True(t, c.Has(validSHA))

	// Corrupt the cached file on disk
	require.NoError(t, os.WriteFile(c.Path(validSHA), []byte("corrupted data"), 0o644))

	// Has should detect mismatch and auto-purge
	assert.False(t, c.Has(validSHA))
	_, err := os.Stat(c.Path(validSHA))
	assert.True(t, os.IsNotExist(err), "corrupted file must be purged")
}
