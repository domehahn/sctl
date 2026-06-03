//go:build integration

package integration

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testServer struct {
	*httptest.Server
	skills map[string][]byte
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	ts := &testServer{
		skills: make(map[string][]byte),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/skills/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path
		data, ok := ts.skills[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	ts.Server = httptest.NewServer(mux)
	t.Cleanup(ts.Server.Close)
	return ts
}

// addSkill registers a pre-built ZIP for the given skill name and version.
// Returns the URL path and SHA256 of the ZIP.
func (ts *testServer) addSkill(t *testing.T, name, version string, files map[string]string) (urlPath, sha256sum string) {
	t.Helper()
	zipData := buildZIP(t, files)
	h := sha256.New()
	h.Write(zipData)
	sha256sum = hex.EncodeToString(h.Sum(nil))

	path := fmt.Sprintf("/skills/%s/%s/%s-%s.zip", name, version, name, version)
	ts.skills[path] = zipData
	return path, sha256sum
}

func (ts *testServer) skillURL(path string) string {
	return ts.Server.URL + path
}

func buildZIP(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = io.WriteString(w, content)
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func writeFixtureSkill(t *testing.T) string {
	t.Helper()
	src := filepath.Join("..", "..", "testdata", "skills", "example-skill")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("testdata fixture not found at %s", src)
	}
	return src
}
