package registry_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/domehahn/sklib/spec"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── helpers ───────────────────────────────────────────────────────────────

func genericHTTPReg(t *testing.T, srv *httptest.Server, endpoints map[string]string) *registry.GenericHTTPRegistry {
	t.Helper()
	rc := config.RegistryConfig{
		Type:      "generic-http",
		URL:       srv.URL,
		Endpoints: endpoints,
	}
	return registry.NewGenericHTTPRegistry("test", rc)
}

func skillForgeReg(t *testing.T, srv *httptest.Server) *registry.GenericHTTPRegistry {
	t.Helper()
	rc := config.RegistryConfig{URL: srv.URL}
	return registry.NewSkillForgeRegistry("sf", rc)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// ── Capabilities ──────────────────────────────────────────────────────────

func TestGenericHTTPCapabilitiesFromServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/caps", r.URL.Path)
		writeJSON(w, registry.RegistryCapabilities{Resolve: true, Download: true, Checksums: true})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"capabilities": "/caps"})
	caps, err := reg.Capabilities(context.Background())
	require.NoError(t, err)
	assert.True(t, caps.Resolve)
	assert.True(t, caps.Download)
	assert.True(t, caps.Checksums)
	assert.False(t, caps.Publish)
}

func TestGenericHTTPCapabilitiesFallbackConservative(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// No capabilities endpoint — fallback must NOT assume SemVer/Checksums.
	reg := genericHTTPReg(t, srv, map[string]string{
		"resolve":  "/resolve",
		"download": "/download",
	})
	caps, err := reg.Capabilities(context.Background())
	require.NoError(t, err)
	assert.True(t, caps.Resolve)
	assert.True(t, caps.Download)
	assert.False(t, caps.SemVerConstraints, "SemVerConstraints must default false without explicit config")
	assert.False(t, caps.Checksums, "Checksums must default false without explicit config")
}

func TestGenericHTTPCapabilitiesStaticCapsOverride(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	rc := config.RegistryConfig{
		Type:         "generic-http",
		URL:          srv.URL,
		Endpoints:    map[string]string{"resolve": "/resolve"},
		Capabilities: map[string]bool{"semver_constraints": true, "checksums": true},
	}
	reg := registry.NewGenericHTTPRegistry("test", rc)
	caps, err := reg.Capabilities(context.Background())
	require.NoError(t, err)
	assert.True(t, caps.SemVerConstraints)
	assert.True(t, caps.Checksums)
}

func TestGenericHTTPCapabilitiesServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	// 500 from capabilities endpoint → fall back to static (no crash).
	reg := genericHTTPReg(t, srv, map[string]string{
		"capabilities": "/caps",
		"resolve":      "/resolve",
	})
	caps, err := reg.Capabilities(context.Background())
	require.NoError(t, err)
	assert.True(t, caps.Resolve)
	assert.False(t, caps.Checksums)
}

// ── Resolve ───────────────────────────────────────────────────────────────

func TestGenericHTTPResolve(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/resolve", r.URL.Path)
		writeJSON(w, registry.ResolvedArtifact{
			Name:           "my-skill",
			Version:        "1.2.3",
			DownloadURL:    "http://example.com/my-skill-1.2.3.zip",
			SHA256:         "abc",
			CompatibleWith: []spec.Platform{spec.PlatformClaudeCode},
		})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"resolve": "/resolve"})
	art, err := reg.Resolve(context.Background(), registry.ResolveRequest{
		Ref:        registry.SkillRef{Namespace: "default", Name: "my-skill"},
		Constraint: "^1.0.0",
	})
	require.NoError(t, err)
	assert.Equal(t, "my-skill", art.Name)
	assert.Equal(t, "1.2.3", art.Version)
	assert.Equal(t, []spec.Platform{spec.PlatformClaudeCode}, art.CompatibleWith)
}

func TestGenericHTTPResolveNoEndpoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{})
	_, err := reg.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "x"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resolve")
}

func TestGenericHTTPResolveHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"resolve": "/resolve"})
	_, err := reg.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "x"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

// ── Download ──────────────────────────────────────────────────────────────

func TestGenericHTTPDownloadDirectURL(t *testing.T) {
	content := []byte("fake-zip-bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{})
	art := &registry.ResolvedArtifact{DownloadURL: srv.URL + "/skill.zip"}
	var buf []byte
	require.NoError(t, reg.Download(context.Background(), art, &writeCollector{buf: &buf}))
	assert.Equal(t, content, buf)
}

func TestGenericHTTPDownloadViaEndpoint(t *testing.T) {
	content := []byte("zip-from-endpoint")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"download": "/dl/{namespace}/{name}/{version}"})
	art := &registry.ResolvedArtifact{Namespace: "default", Name: "skill", Version: "1.0.0"}
	var buf []byte
	require.NoError(t, reg.Download(context.Background(), art, &writeCollector{buf: &buf}))
	assert.Equal(t, content, buf)
}

func TestGenericHTTPDownloadHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{})
	art := &registry.ResolvedArtifact{DownloadURL: srv.URL + "/skill.zip"}
	err := reg.Download(context.Background(), art, io.Discard)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403")
}

// ── Search ────────────────────────────────────────────────────────────────

func TestGenericHTTPSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "query=foo")
		writeJSON(w, map[string]any{
			"skills": []registry.SkillSearchResult{
				{Name: "foo-skill", LatestVersion: "1.0.0"},
			},
		})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"search": "/search?query={query}"})
	results, err := reg.Search(context.Background(), registry.SearchRequest{Query: "foo"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "foo-skill", results[0].Name)
}

// ── Info ──────────────────────────────────────────────────────────────────

func TestGenericHTTPInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, registry.SkillInfo{
			Namespace:     "default",
			Name:          "my-skill",
			LatestVersion: "2.0.0",
			Description:   "Does things",
		})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"info": "/skills/{namespace}/{name}"})
	info, err := reg.Info(context.Background(), registry.SkillRef{Namespace: "default", Name: "my-skill"})
	require.NoError(t, err)
	assert.Equal(t, "my-skill", info.Name)
	assert.Equal(t, "2.0.0", info.LatestVersion)
}

// ── ListVersions ──────────────────────────────────────────────────────────

func TestGenericHTTPListVersions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]any{
			"versions": []registry.VersionInfo{
				{Version: "1.0.0"},
				{Version: "2.0.0", Yanked: true},
			},
		})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{"versions": "/versions/{namespace}/{name}"})
	versions, err := reg.ListVersions(context.Background(), registry.SkillRef{Namespace: "default", Name: "skill"})
	require.NoError(t, err)
	require.Len(t, versions, 2)
	assert.Equal(t, "1.0.0", versions[0].Version)
	assert.True(t, versions[1].Yanked)
}

// ── Publish ───────────────────────────────────────────────────────────────

func TestGenericHTTPPublish(t *testing.T) {
	var receivedSHA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		receivedSHA = r.Header.Get("X-SKPM-SHA256")
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, registry.PublishResult{Name: "my-skill", Version: "1.0.0", Created: true})
	}))
	defer srv.Close()

	tmp := t.TempDir()
	artifactPath := tmp + "/skill.zip"
	require.NoError(t, writeFile(artifactPath, []byte("zip")))

	reg := genericHTTPReg(t, srv, map[string]string{"publish": "/skills/{namespace}/{name}/versions/{version}"})
	result, err := reg.Publish(context.Background(), registry.PublishRequest{
		ArtifactPath: artifactPath,
		Manifest:     skill.SkillManifest{Name: "my-skill", Version: "1.0.0"},
		SHA256:       "deadbeef",
	})
	require.NoError(t, err)
	assert.True(t, result.Created)
	assert.Equal(t, "deadbeef", receivedSHA)
}

func TestGenericHTTPPublishContentTypeByPackageType(t *testing.T) {
	cases := []struct {
		packageType string
		want        string
	}{
		{"", "application/zip"},
		{"zip", "application/zip"},
		{"tgz", "application/gzip"},
		{"tar.gz", "application/gzip"},
	}

	for _, c := range cases {
		var receivedContentType string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedContentType = r.Header.Get("Content-Type")
			w.WriteHeader(http.StatusCreated)
			writeJSON(w, registry.PublishResult{Name: "my-skill", Version: "1.0.0", Created: true})
		}))

		tmp := t.TempDir()
		artifactPath := tmp + "/skill.zip"
		require.NoError(t, writeFile(artifactPath, []byte("zip")))

		reg := genericHTTPReg(t, srv, map[string]string{"publish": "/skills/{namespace}/{name}/versions/{version}"})
		_, err := reg.Publish(context.Background(), registry.PublishRequest{
			ArtifactPath: artifactPath,
			Manifest:     skill.SkillManifest{Name: "my-skill", Version: "1.0.0"},
			PackageType:  c.packageType,
		})
		require.NoError(t, err)
		assert.Equal(t, c.want, receivedContentType, "package_type=%q", c.packageType)
		srv.Close()
	}
}

func TestGenericHTTPPublishNoEndpoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{})
	_, err := reg.Publish(context.Background(), registry.PublishRequest{
		ArtifactPath: "/nonexistent.zip",
		Manifest:     skill.SkillManifest{Name: "x", Version: "1.0.0"},
	})
	require.Error(t, err)
}

// ── Governance ────────────────────────────────────────────────────────────

func TestGenericHTTPDeprecate(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{
		"deprecate": "/skills/{namespace}/{name}/versions/{version}/deprecate",
	})
	err := reg.Deprecate(context.Background(), registry.SkillVersionRef{Namespace: "default", Name: "skill", Version: "1.0.0"}, "old")
	require.NoError(t, err)
	assert.Contains(t, path, "deprecate")
}

func TestGenericHTTPYank(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{
		"yank": "/skills/{namespace}/{name}/versions/{version}/yank",
	})
	err := reg.Yank(context.Background(), registry.SkillVersionRef{Namespace: "default", Name: "skill", Version: "1.0.0"}, "broken")
	require.NoError(t, err)
	assert.Contains(t, path, "yank")
}

func TestGenericHTTPUnyank(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{
		"unyank": "/skills/{namespace}/{name}/versions/{version}/unyank",
	})
	err := reg.Unyank(context.Background(), registry.SkillVersionRef{Namespace: "default", Name: "skill", Version: "1.0.0"})
	require.NoError(t, err)
	assert.Contains(t, path, "unyank")
}

func TestGenericHTTPAttest(t *testing.T) {
	var path, method string
	var received registry.AttestationRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, method = r.URL.Path, r.Method
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, registry.AttestationRecord{ID: 7, Type: received.Type, Digest: received.Digest, CreatedBy: "alice"})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{
		"attest": "/skills/{namespace}/{name}/versions/{version}/attestations",
	})
	rec, err := reg.Attest(context.Background(), registry.SkillVersionRef{Namespace: "default", Name: "skill", Version: "1.0.0"}, registry.AttestationRequest{
		Type:      "scan",
		Digest:    "deadbeef",
		Predicate: json.RawMessage(`{"verdict":"clear"}`),
	})
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, method)
	assert.Contains(t, path, "attestations")
	assert.Equal(t, "scan", received.Type)
	assert.Equal(t, "deadbeef", received.Digest)
	assert.Equal(t, int64(7), rec.ID)
	assert.Equal(t, "alice", rec.CreatedBy)
}

func TestGenericHTTPAttestNoEndpoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{})
	_, err := reg.Attest(context.Background(), registry.SkillVersionRef{Namespace: "default", Name: "skill", Version: "1.0.0"}, registry.AttestationRequest{Type: "scan", Digest: "x"})
	require.Error(t, err)
}

func TestGenericHTTPListAttestations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		writeJSON(w, map[string]any{"attestations": []registry.AttestationRecord{
			{ID: 1, Type: "scan", Digest: "deadbeef", CreatedBy: "alice"},
		}})
	}))
	defer srv.Close()

	reg := genericHTTPReg(t, srv, map[string]string{
		"attestations": "/skills/{namespace}/{name}/versions/{version}/attestations",
	})
	records, err := reg.ListAttestations(context.Background(), registry.SkillVersionRef{Namespace: "default", Name: "skill", Version: "1.0.0"})
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "deadbeef", records[0].Digest)
}

// ── SkillForge defaults ───────────────────────────────────────────────────

func TestSkillForgeDefaultEndpoints(t *testing.T) {
	called := map[string]bool{}
	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called[r.URL.Path] = true
		switch {
		case r.URL.Path == "/api/v1/capabilities":
			writeJSON(w, registry.RegistryCapabilities{Resolve: true, Download: true, Publish: true, SemVerConstraints: true, Checksums: true})
		case r.URL.Path == "/api/v1/skills" || r.URL.Path == "/api/v1/skills/":
			writeJSON(w, map[string]any{"skills": []registry.SkillSearchResult{{Name: "sf-skill"}}})
		case r.URL.Path == "/api/v1/skills/default/sf-skill":
			writeJSON(w, registry.SkillInfo{Name: "sf-skill", LatestVersion: "1.0.0"})
		case r.URL.Path == "/api/v1/skills/default/sf-skill/versions":
			writeJSON(w, map[string]any{"versions": []registry.VersionInfo{{Version: "1.0.0"}}})
		case r.URL.Path == "/api/v1/skills/default/sf-skill/resolve":
			writeJSON(w, registry.ResolvedArtifact{Name: "sf-skill", Version: "1.0.0", DownloadURL: srvURL + "/dl"})
		case r.URL.Path == "/dl":
			w.Write([]byte("zip-bytes"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	srvURL = srv.URL
	defer srv.Close()

	reg := skillForgeReg(t, srv)

	// capabilities — fetched live from server
	caps, err := reg.Capabilities(context.Background())
	require.NoError(t, err)
	assert.True(t, caps.SemVerConstraints)
	assert.True(t, caps.Checksums)

	// resolve
	art, err := reg.Resolve(context.Background(), registry.ResolveRequest{
		Ref:        registry.SkillRef{Namespace: "default", Name: "sf-skill"},
		Constraint: "^1.0.0",
	})
	require.NoError(t, err)
	assert.Equal(t, "sf-skill", art.Name)

	// download
	var buf []byte
	require.NoError(t, reg.Download(context.Background(), art, &writeCollector{buf: &buf}))
	assert.Equal(t, []byte("zip-bytes"), buf)

	// search
	results, err := reg.Search(context.Background(), registry.SearchRequest{Query: "sf"})
	require.NoError(t, err)
	require.Len(t, results, 1)

	// info
	info, err := reg.Info(context.Background(), registry.SkillRef{Namespace: "default", Name: "sf-skill"})
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", info.LatestVersion)

	// versions
	versions, err := reg.ListVersions(context.Background(), registry.SkillRef{Namespace: "default", Name: "sf-skill"})
	require.NoError(t, err)
	require.Len(t, versions, 1)
}

func TestSkillForgeEndpointOverride(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/custom/resolve" {
			writeJSON(w, registry.ResolvedArtifact{Name: "x", Version: "1.0.0"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	rc := config.RegistryConfig{
		URL: srv.URL,
		Endpoints: map[string]string{
			"resolve": "/custom/resolve",
		},
	}
	reg := registry.NewSkillForgeRegistry("sf", rc)
	art, err := reg.Resolve(context.Background(), registry.ResolveRequest{Ref: registry.SkillRef{Name: "x"}})
	require.NoError(t, err)
	assert.Equal(t, "x", art.Name)
}

// ── helpers ───────────────────────────────────────────────────────────────

func writeFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
