package cli_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// keyID reproduces skil's internal/signing.KeyID / skpm's
// internal/attestation.KeyID, kept independent here so this test does not
// depend on either package's internals — only their documented output
// shape (a JSON object with a "signature" field).
func keyID(pub ed25519.PublicKey) string {
	sum := sha256.Sum256(pub)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// signPredicate signs attestation using the same canonical-JSON-then-Ed25519
// scheme as skil's internal/signing and skpm's internal/attestation, so
// this test proves the CLI wiring, not just the already-covered Verify
// function.
func signPredicate(t *testing.T, attestation map[string]any, priv ed25519.PrivateKey, kid string) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(attestation)
	require.NoError(t, err)
	var generic any
	require.NoError(t, json.Unmarshal(raw, &generic))
	canonical, err := json.Marshal(generic)
	require.NoError(t, err)

	attestation["signature"] = map[string]string{
		"provider":  "builtin.ed25519",
		"algorithm": "Ed25519",
		"key_id":    kid,
		"value":     base64.StdEncoding.EncodeToString(ed25519.Sign(priv, canonical)),
	}
	out, err := json.Marshal(attestation)
	require.NoError(t, err)
	return out
}

func writeSkpmConfig(t *testing.T, registryURL, trustedKeyID, trustedKeyB64 string) {
	t.Helper()
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)
	skpmDir := filepath.Join(cfgDir, "skpm")
	require.NoError(t, os.MkdirAll(skpmDir, 0o755))

	cfg := fmt.Sprintf(`default_registry: test
registries:
  test:
    type: generic-http
    url: %s
    endpoints:
      attestations: /skills/{namespace}/{name}/versions/{version}/attestations
trusted_signers:
  %s: %s
`, registryURL, trustedKeyID, trustedKeyB64)
	require.NoError(t, os.WriteFile(filepath.Join(skpmDir, "config.yaml"), []byte(cfg), 0o600))
}

func attestationsServer(t *testing.T, predicate json.RawMessage) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"attestations": []map[string]any{
				{"type": "scan", "digest": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "predicate": predicate, "created_by": "skil"},
			},
		})
	}))
}

func TestAttestationsCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "attestations <skill>@<version>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("verify"))
			break
		}
	}
	assert.True(t, found, "attestations command should be registered on root")
}

func TestAttestationsVerifyAcceptsGenuineSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	kid := keyID(pub)
	predicate := signPredicate(t, map[string]any{
		"version": 1,
		"subject": map[string]any{"name": "demo-skill", "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}, priv, kid)

	srv := attestationsServer(t, predicate)
	defer srv.Close()
	writeSkpmConfig(t, srv.URL, kid, base64.StdEncoding.EncodeToString(pub))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"attestations", "demo-skill@1.0.0", "--verify"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "signature=verified")
	assert.Contains(t, out, kid)
}

func TestAttestationsVerifyRejectsUntrustedKey(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	kid := keyID(pub)
	predicate := signPredicate(t, map[string]any{
		"version": 1,
		"subject": map[string]any{"name": "demo-skill", "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}, priv, kid)

	srv := attestationsServer(t, predicate)
	defer srv.Close()
	// Config trusts a different key than the one that actually signed it.
	otherPub, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	writeSkpmConfig(t, srv.URL, keyID(otherPub), base64.StdEncoding.EncodeToString(otherPub))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"attestations", "demo-skill@1.0.0", "--verify"})
	root.SetOut(&buf)
	require.ErrorContains(t, root.Execute(), "ATTESTATION_VERIFICATION_FAILED")

	out := buf.String()
	assert.Contains(t, out, "NOT-VERIFIED")
	assert.Contains(t, out, "untrusted")
}

func TestAttestationsWithoutVerifyDoesNotCheckSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	kid := keyID(pub)
	predicate := signPredicate(t, map[string]any{
		"version": 1,
		"subject": map[string]any{"name": "demo-skill", "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}, priv, kid)

	srv := attestationsServer(t, predicate)
	defer srv.Close()
	writeSkpmConfig(t, srv.URL, kid, base64.StdEncoding.EncodeToString(pub))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"attestations", "demo-skill@1.0.0"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.NotContains(t, out, "signature=")
}
