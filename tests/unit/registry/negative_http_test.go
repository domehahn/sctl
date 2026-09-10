package registry_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericHTTPRegistry_ErrorStatusCodes(t *testing.T) {
	statusCodes := []int{
		http.StatusUnauthorized,        // 401
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusConflict,            // 409
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
	}

	for _, code := range statusCodes {
		t.Run(fmt.Sprintf("HTTP_%d", code), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
				_, _ = w.Write([]byte("error response"))
			}))
			defer server.Close()

			rc := config.RegistryConfig{
				URL:  server.URL,
				Type: "skillforge",
			}
			reg := registry.NewSkillForgeRegistry("test-sf", rc)
			ref := registry.SkillRef{Namespace: "default", Name: "my-skill"}
			vref := registry.SkillVersionRef{Namespace: "default", Name: "my-skill", Version: "1.0.0"}

			expectedSubstring := fmt.Sprintf("HTTP %d", code)

			// Resolve
			_, err := reg.Resolve(context.Background(), registry.ResolveRequest{Ref: ref, Constraint: "1.0.0"})
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)

			// Download
			var buf bytes.Buffer
			err = reg.Download(context.Background(), &registry.ResolvedArtifact{
				Namespace:   "default",
				Name:        "my-skill",
				Version:     "1.0.0",
				DownloadURL: server.URL + "/download",
			}, &buf)
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)

			// Info
			_, err = reg.Info(context.Background(), ref)
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)

			// ListVersions
			_, err = reg.ListVersions(context.Background(), ref)
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)

			// Deprecate
			err = reg.Deprecate(context.Background(), vref, "obsolete")
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)

			// Yank
			err = reg.Yank(context.Background(), vref, "security bug")
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)

			// Unyank
			err = reg.Unyank(context.Background(), vref)
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedSubstring)
		})
	}
}

func TestGenericHTTPRegistry_CorruptAndPartialDownload(t *testing.T) {
	t.Run("partial response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1000")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("partial"))
		}))
		defer server.Close()

		rc := config.RegistryConfig{URL: server.URL, Type: "skillforge"}
		reg := registry.NewSkillForgeRegistry("test-sf", rc)

		var buf bytes.Buffer
		_ = reg.Download(context.Background(), &registry.ResolvedArtifact{
			Namespace:   "default",
			Name:        "my-skill",
			Version:     "1.0.0",
			DownloadURL: server.URL + "/download",
		}, &buf)
	})

	t.Run("timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		rc := config.RegistryConfig{URL: server.URL, Type: "skillforge"}
		reg := registry.NewSkillForgeRegistry("test-sf", rc)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := reg.Resolve(ctx, registry.ResolveRequest{Ref: registry.SkillRef{Namespace: "default", Name: "my-skill"}})
		require.Error(t, err)
	})
}
