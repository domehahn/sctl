package publisher_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/publisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeZIP(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.zip")
	require.NoError(t, err)
	_, err = f.WriteString("PK fake zip content")
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

// ── Artifactory ───────────────────────────────────────────────────────────

func TestArtifactoryPublish(t *testing.T) {
	var receivedPath, receivedAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		assert.NotEmpty(t, body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	p := publisher.NewArtifactoryPublisher(srv.URL, "agent-skills", "my-token")
	result, err := p.Publish(context.Background(), "my-skill", "1.0.0", makeZIP(t), "abc123")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.0.0", result.Version)
	assert.Equal(t, "abc123", result.SHA256)
	assert.Equal(t, "/agent-skills/my-skill/1.0.0/my-skill-1.0.0.zip", receivedPath)
	assert.Equal(t, "Bearer my-token", receivedAuth)
}

func TestArtifactoryPublishHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("access denied"))
	}))
	defer srv.Close()

	p := publisher.NewArtifactoryPublisher(srv.URL, "skills", "bad-token")
	_, err := p.Publish(context.Background(), "skill", "1.0.0", makeZIP(t), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403")
}

func TestArtifactoryPublishMissingZIP(t *testing.T) {
	p := publisher.NewArtifactoryPublisher("http://localhost", "skills", "")
	_, err := p.Publish(context.Background(), "skill", "1.0.0", "/nonexistent/path.zip", "")
	require.Error(t, err)
}

// ── Local ─────────────────────────────────────────────────────────────────

func TestLocalPublish(t *testing.T) {
	baseDir := t.TempDir()
	p := publisher.NewLocalPublisher(baseDir)
	result, err := p.Publish(context.Background(), "my-skill", "1.2.3", makeZIP(t), "deadbeef")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.2.3", result.Version)
	assert.Equal(t, "deadbeef", result.SHA256)
	assert.Contains(t, result.DownloadURL, "my-skill-1.2.3.zip")
	assert.FileExists(t, filepath.Join(baseDir, "my-skill", "1.2.3", "my-skill-1.2.3.zip"))
}

func TestLocalPublishMissingZIP(t *testing.T) {
	p := publisher.NewLocalPublisher(t.TempDir())
	_, err := p.Publish(context.Background(), "skill", "1.0.0", "/no/such/file.zip", "")
	require.Error(t, err)
}

func TestCopyFile(t *testing.T) {
	src := makeZIP(t)
	dst := filepath.Join(t.TempDir(), "out.zip")
	require.NoError(t, publisher.CopyFile(src, dst))
	assert.FileExists(t, dst)
	_, err := os.Stat(dst + ".tmp")
	assert.True(t, os.IsNotExist(err))
}

// ── GitHub ────────────────────────────────────────────────────────────────

func TestNewGitHubPublisherInvalidSlug(t *testing.T) {
	_, err := publisher.NewGitHubPublisher("not-a-slug", "", "prefixed")
	assert.Error(t, err)
}

func TestNewGitHubPublisherValid(t *testing.T) {
	p, err := publisher.NewGitHubPublisher("org/repo", "token", "prefixed")
	require.NoError(t, err)
	assert.NotNil(t, p)
}

// ── GitLab ────────────────────────────────────────────────────────────────

func TestNewGitLabPublisherInvalidBaseURL(t *testing.T) {
	_, err := publisher.NewGitLabPublisher("://bad-url", "platform/skills", "token", "prefixed")
	assert.Error(t, err)
}

func TestGitLabPublisherUploadGenericPackage(t *testing.T) {
	var receivedBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		assert.Equal(t, "mytoken", r.Header.Get("PRIVATE-TOKEN"))
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	p, err := publisher.NewGitLabPublisher(srv.URL, "platform/skills", "mytoken", "prefixed")
	require.NoError(t, err)
	url, err := p.UploadGenericPackage(context.Background(), "my-skill", "1.0.0", "my-skill-1.0.0.zip", makeZIP(t))
	require.NoError(t, err)
	assert.Contains(t, url, "my-skill")
	assert.NotEmpty(t, receivedBody)
}

func TestGitLabPublisherCreateRelease(t *testing.T) {
	releaseCreated := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
		case http.MethodPost:
			releaseCreated = true
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			assert.Equal(t, "my-skill/v1.0.0", body["tag_name"])
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"tag_name": "my-skill/v1.0.0"})
		}
	}))
	defer srv.Close()

	p, err := publisher.NewGitLabPublisher(srv.URL, "platform/skills", "token", "prefixed")
	require.NoError(t, err)
	err = p.CreateOrUpdateRelease(context.Background(),
		"my-skill/v1.0.0", "my-skill", "1.0.0", "abc123",
		"my-skill-1.0.0.zip", srv.URL+"/pkg.zip",
	)
	require.NoError(t, err)
	assert.True(t, releaseCreated)
}

// bufWriter collects bytes for assertions.
type bufWriter struct{ bytes.Buffer }
