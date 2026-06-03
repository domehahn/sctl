package registry

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalRegistryResolveExact(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "my-skill", "1.2.3")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "my-skill-1.2.3.zip"), []byte("zip-content"), 0o644))

	r := NewLocalRegistry(base)
	art, err := r.Resolve(context.Background(), "my-skill", "1.2.3")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", art.Name)
	assert.Equal(t, "1.2.3", art.Version)
}

func TestLocalRegistryResolveLatest(t *testing.T) {
	base := t.TempDir()
	for _, v := range []string{"1.0.0", "2.1.0", "1.9.0"} {
		dir := filepath.Join(base, "skill", v)
		require.NoError(t, os.MkdirAll(dir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "skill-"+v+".zip"), []byte("zip"), 0o644))
	}

	r := NewLocalRegistry(base)
	art, err := r.Resolve(context.Background(), "skill", "")
	require.NoError(t, err)
	assert.Equal(t, "2.1.0", art.Version)
}

func TestLocalRegistryDownload(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "skill", "1.0.0")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	content := []byte("zip-bytes")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "skill-1.0.0.zip"), content, 0o644))

	r := NewLocalRegistry(base)
	art, err := r.Resolve(context.Background(), "skill", "1.0.0")
	require.NoError(t, err)

	var buf []byte
	w := &writeCollector{buf: &buf}
	require.NoError(t, r.Download(context.Background(), art, w))
	assert.Equal(t, content, buf)
}

func TestLocalRegistryMissingVersion(t *testing.T) {
	base := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(base, "skill"), 0o755))

	r := NewLocalRegistry(base)
	_, err := r.Resolve(context.Background(), "skill", "")
	assert.Error(t, err)
}

func TestArtifactoryResolveExact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := NewArtifactoryRegistry(srv.URL, "skills", "token")
	art, err := reg.Resolve(context.Background(), "my-skill", "1.0.0")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", art.Name)
	assert.Equal(t, "1.0.0", art.Version)
	assert.Contains(t, art.DownloadURL, "my-skill-1.0.0.zip")
}

func TestArtifactoryDownload(t *testing.T) {
	content := []byte("fake-zip")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}))
	defer srv.Close()

	reg := NewArtifactoryRegistry(srv.URL, "skills", "")
	art := &ResolvedArtifact{DownloadURL: srv.URL + "/skills/my-skill/1.0.0/my-skill-1.0.0.zip"}

	var buf []byte
	w := &writeCollector{buf: &buf}
	require.NoError(t, reg.Download(context.Background(), art, w))
	assert.Equal(t, content, buf)
}

func TestArtifactoryDownload404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	reg := NewArtifactoryRegistry(srv.URL, "skills", "")
	art := &ResolvedArtifact{DownloadURL: srv.URL + "/not-found.zip"}

	err := reg.Download(context.Background(), art, io.Discard)
	assert.Error(t, err)
}

type writeCollector struct {
	buf *[]byte
}

func (w *writeCollector) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}
