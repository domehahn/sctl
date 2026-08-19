package registry_test

// Publish-path coverage for the registry.PublishingRegistry backends. These
// tests replace tests/unit/publisher/publisher_test.go, which covered the
// same backends (github/gitlab/artifactory/local) through the now-removed
// internal/publisher package before skpm publish/release were migrated onto
// internal/registry (registry.PublishingRegistry) so that SkillForge and
// generic-http registries could publish through the same adapter used for
// pull/resolve.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makePublishZIP(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.zip")
	require.NoError(t, err)
	_, err = f.WriteString("PK fake zip content")
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func publishReq(t *testing.T, name, version, sha256 string) registry.PublishRequest {
	return registry.PublishRequest{
		ArtifactPath: makePublishZIP(t),
		Manifest:     skill.SkillManifest{Name: name, Version: version},
		SHA256:       sha256,
		TagFormat:    "prefixed",
	}
}

// ── Artifactory ───────────────────────────────────────────────────────────

func TestArtifactoryRegistryPublish(t *testing.T) {
	var receivedPath, receivedAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		assert.NotEmpty(t, body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	r := registry.NewArtifactoryRegistry(srv.URL, "agent-skills", "my-token")
	result, err := r.Publish(context.Background(), publishReq(t, "my-skill", "1.0.0", "abc123"))
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.0.0", result.Version)
	assert.Equal(t, "abc123", result.SHA256)
	assert.Equal(t, "/agent-skills/my-skill/1.0.0/my-skill-1.0.0.zip", receivedPath)
	assert.Equal(t, "Bearer my-token", receivedAuth)
}

func TestArtifactoryRegistryPublishHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("access denied"))
	}))
	defer srv.Close()

	r := registry.NewArtifactoryRegistry(srv.URL, "skills", "bad-token")
	_, err := r.Publish(context.Background(), publishReq(t, "skill", "1.0.0", ""))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403")
}

func TestArtifactoryRegistryPublishMissingZIP(t *testing.T) {
	r := registry.NewArtifactoryRegistry("http://localhost", "skills", "")
	req := publishReq(t, "skill", "1.0.0", "")
	req.ArtifactPath = "/nonexistent/path.zip"
	_, err := r.Publish(context.Background(), req)
	require.Error(t, err)
}

// ── Local ─────────────────────────────────────────────────────────────────

func TestLocalRegistryPublish(t *testing.T) {
	baseDir := t.TempDir()
	r := registry.NewLocalRegistry(baseDir)
	result, err := r.Publish(context.Background(), publishReq(t, "my-skill", "1.2.3", "deadbeef"))
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.2.3", result.Version)
	assert.Equal(t, "deadbeef", result.SHA256)
	assert.Contains(t, result.DownloadURL, "my-skill-1.2.3.zip")
	assert.FileExists(t, filepath.Join(baseDir, "my-skill", "1.2.3", "my-skill-1.2.3.zip"))
}

func TestLocalRegistryPublishMissingZIP(t *testing.T) {
	r := registry.NewLocalRegistry(t.TempDir())
	req := publishReq(t, "skill", "1.0.0", "")
	req.ArtifactPath = "/no/such/file.zip"
	_, err := r.Publish(context.Background(), req)
	require.Error(t, err)
}

// ── GitLab ────────────────────────────────────────────────────────────────

func TestGitLabRegistryPublish(t *testing.T) {
	releaseCreated := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut:
			assert.Equal(t, "mytoken", r.Header.Get("PRIVATE-TOKEN"))
			body, _ := io.ReadAll(r.Body)
			assert.NotEmpty(t, body)
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
		case r.Method == http.MethodPost:
			releaseCreated = true
			var reqBody map[string]interface{}
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.Equal(t, "my-skill/v1.0.0", reqBody["tag_name"])
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"tag_name": "my-skill/v1.0.0"})
		}
	}))
	defer srv.Close()

	r, err := registry.NewGitLabRegistry(srv.URL, "platform/skills", "mytoken")
	require.NoError(t, err)

	result, err := r.Publish(context.Background(), publishReq(t, "my-skill", "1.0.0", "abc123"))
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.0.0", result.Version)
	assert.Contains(t, result.DownloadURL, "my-skill")
	assert.True(t, releaseCreated)
}

func TestGitLabRegistryPublishPlainTagFormat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
		case http.MethodPost:
			var reqBody map[string]interface{}
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.Equal(t, "v1.0.0", reqBody["tag_name"])
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"tag_name": "v1.0.0"})
		}
	}))
	defer srv.Close()

	r, err := registry.NewGitLabRegistry(srv.URL, "platform/skills", "mytoken")
	require.NoError(t, err)

	req := publishReq(t, "my-skill", "1.0.0", "abc123")
	req.TagFormat = "plain"
	_, err = r.Publish(context.Background(), req)
	require.NoError(t, err)
}
