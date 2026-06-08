package registry_test

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Local Registry ────────────────────────────────────────────────────────

func TestLocalRegistryResolveExact(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "my-skill", "1.2.3")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "my-skill-1.2.3.zip"), []byte("zip"), 0o644))

	r := registry.NewLocalRegistry(base)
	art, err := r.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "my-skill"}, Constraint: "1.2.3"})
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

	r := registry.NewLocalRegistry(base)
	art, err := r.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "skill"}})
	require.NoError(t, err)
	assert.Equal(t, "2.1.0", art.Version)
}

func TestLocalRegistryMissingVersion(t *testing.T) {
	base := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(base, "skill"), 0o755))

	r := registry.NewLocalRegistry(base)
	_, err := r.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "skill"}})
	assert.Error(t, err)
}

func TestLocalRegistryMissingArtifact(t *testing.T) {
	base := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(base, "skill", "1.0.0"), 0o755))

	r := registry.NewLocalRegistry(base)
	_, err := r.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "skill"}, Constraint: "1.0.0"})
	assert.Error(t, err)
}

func TestLocalRegistryDownload(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "skill", "1.0.0")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	content := []byte("zip-bytes")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "skill-1.0.0.zip"), content, 0o644))

	r := registry.NewLocalRegistry(base)
	art, err := r.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "skill"}, Constraint: "1.0.0"})
	require.NoError(t, err)

	var buf []byte
	w := &writeCollector{buf: &buf}
	require.NoError(t, r.Download(context.Background(), art, w))
	assert.Equal(t, content, buf)
}

// ── Artifactory ────────────────────────────────────────────────────────────

func TestArtifactoryResolveExact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := registry.NewArtifactoryRegistry(srv.URL, "skills", "token")
	art, err := reg.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "my-skill"}, Constraint: "1.0.0"})
	require.NoError(t, err)
	assert.Equal(t, "my-skill", art.Name)
	assert.Equal(t, "1.0.0", art.Version)
	assert.Contains(t, art.DownloadURL, "my-skill-1.0.0.zip")
}

func TestArtifactoryDownload(t *testing.T) {
	content := []byte("fake-zip")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	reg := registry.NewArtifactoryRegistry(srv.URL, "skills", "")
	art := &registry.ResolvedArtifact{DownloadURL: srv.URL + "/skills/my-skill/1.0.0/my-skill-1.0.0.zip"}

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

	reg := registry.NewArtifactoryRegistry(srv.URL, "skills", "")
	art := &registry.ResolvedArtifact{DownloadURL: srv.URL + "/not-found.zip"}
	err := reg.Download(context.Background(), art, io.Discard)
	assert.Error(t, err)
}

// ── GitHub ────────────────────────────────────────────────────────────────

func TestNewGitHubRegistryInvalidSlug(t *testing.T) {
	_, err := registry.NewGitHubRegistry("not-a-slug", "")
	assert.Error(t, err)
}

func TestNewGitHubRegistryNoToken(t *testing.T) {
	reg, err := registry.NewGitHubRegistry("org/repo", "")
	require.NoError(t, err)
	assert.NotNil(t, reg)
}

// ── GitLab ────────────────────────────────────────────────────────────────

func TestNewGitLabRegistryDefaultBaseURL(t *testing.T) {
	reg, err := registry.NewGitLabRegistry("", "platform/skills", "token")
	require.NoError(t, err)
	assert.NotNil(t, reg)
}

// ── RefDownload ───────────────────────────────────────────────────────────

func buildArchiveZIP(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func TestExtractSkillFromZIPWithSubPath(t *testing.T) {
	archive := buildArchiveZIP(t, map[string]string{
		"owner-repo-abc/skills/my-skill/SKILL.md":   "# My Skill",
		"owner-repo-abc/skills/my-skill/skill.yaml": "name: my-skill",
		"owner-repo-abc/other-file.txt":             "other",
	})

	destDir := t.TempDir()
	err := registry.ExtractSkillFromZIP(bytes.NewReader(archive), int64(len(archive)), "my-skill", "skills/my-skill", destDir)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
	assert.FileExists(t, filepath.Join(destDir, "skill.yaml"))
	assert.NoFileExists(t, filepath.Join(destDir, "other-file.txt"))
}

func TestExtractSkillFromZIPByName(t *testing.T) {
	archive := buildArchiveZIP(t, map[string]string{
		"repo-main/my-skill/SKILL.md":   "# My Skill",
		"repo-main/my-skill/skill.yaml": "name: my-skill",
	})

	destDir := t.TempDir()
	err := registry.ExtractSkillFromZIP(bytes.NewReader(archive), int64(len(archive)), "my-skill", "", destDir)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
}

func TestExtractSkillFromZIPNotFound(t *testing.T) {
	archive := buildArchiveZIP(t, map[string]string{
		"repo-main/other-skill/SKILL.md": "# Other",
	})
	destDir := t.TempDir()
	err := registry.ExtractSkillFromZIP(bytes.NewReader(archive), int64(len(archive)), "my-skill", "skills/my-skill", destDir)
	require.Error(t, err)
}

func TestDownloadRefGitLabHTTP(t *testing.T) {
	archive := buildArchiveZIP(t, map[string]string{
		"project-main-abc/SKILL.md":   "# Skill",
		"project-main-abc/skill.yaml": "name: my-skill",
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "sha=main")
		w.Write(archive)
	}))
	defer srv.Close()

	reg, err := registry.NewGitLabRegistry(srv.URL, "platform/skills", "token")
	require.NoError(t, err)

	destDir := t.TempDir()
	err = registry.DownloadGitLabRef(context.Background(), reg, "my-skill", "main", "", destDir)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
}

func TestDownloadRefGitLabHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	reg, err := registry.NewGitLabRegistry(srv.URL, "platform/skills", "")
	require.NoError(t, err)

	err = registry.DownloadGitLabRef(context.Background(), reg, "skill", "missing-branch", "", t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestDownloadRefUnsupportedRegistry(t *testing.T) {
	err := registry.DownloadRef(context.Background(), registry.NewLocalRegistry("/tmp"), "skill", "main", "", t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only supported for github and gitlab")
}

func TestArtifactoryResolveLatestVersions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/storage/") {
			w.Write([]byte(`{"children":[{"uri":"/1.0.0","folder":true},{"uri":"/2.1.0","folder":true}]}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := registry.NewArtifactoryRegistry(srv.URL, "skills", "")
	art, err := reg.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "my-skill"}})
	require.NoError(t, err)
	assert.Equal(t, "2.1.0", art.Version)
}

// ── helpers ───────────────────────────────────────────────────────────────

type writeCollector struct{ buf *[]byte }

func (w *writeCollector) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}
