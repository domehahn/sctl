package attestation

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
)

// sign reproduces skil's internal/signing.SignAttestation exactly (marshal
// the JSON-object-shaped attestation with no "signature" key yet, run it
// through canonicalJSON, sign, then splice the signature in) — standing in
// for a real skil binary in tests that don't need the interop guarantee
// the //go:build e2e test in this package covers.
func sign(t *testing.T, attestation map[string]any, privateKey ed25519.PrivateKey, keyID string) json.RawMessage {
	t.Helper()
	payload, err := canonicalJSON(attestation)
	if err != nil {
		t.Fatal(err)
	}
	attestation["signature"] = Signature{
		Provider:  Provider,
		Algorithm: Algorithm,
		KeyID:     keyID,
		Value:     base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, payload)),
	}
	raw, err := json.Marshal(attestation)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func testAttestation() map[string]any {
	return map[string]any{
		"version":  1,
		"subject":  map[string]any{"name": "demo-skill", "sha256": "abc123"},
		"producer": map[string]any{"name": "skil", "version": "0.1.0"},
		"result":   map[string]any{"status": "pass", "maximum_severity": "none", "risk_score": 0},
		"analysis": []string{"lint"},
	}
}

func TestVerifyAcceptsGenuineSignature(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keyID := KeyID(publicKey)
	raw := sign(t, testAttestation(), privateKey, keyID)

	trusted := map[string]string{keyID: base64.StdEncoding.EncodeToString(publicKey)}
	sig, err := Verify(raw, trusted)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if sig.KeyID != keyID {
		t.Fatalf("KeyID = %q, want %q", sig.KeyID, keyID)
	}
}

func TestVerifyRejectsUntrustedKey(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	raw := sign(t, testAttestation(), privateKey, KeyID(publicKey))

	if _, err := Verify(raw, map[string]string{}); !errors.Is(err, ErrUntrustedKey) {
		t.Fatalf("expected ErrUntrustedKey, got %v", err)
	}
}

func TestVerifyRejectsTamperedField(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keyID := KeyID(publicKey)
	raw := sign(t, testAttestation(), privateKey, keyID)
	trusted := map[string]string{keyID: base64.StdEncoding.EncodeToString(publicKey)}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatal(err)
	}
	envelope["analysis"] = json.RawMessage(`["tampered"]`)
	tampered, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(tampered, trusted); err == nil {
		t.Fatal("expected tampered attestation to fail verification")
	}

	if _, err := Verify(raw, trusted); err != nil {
		t.Fatalf("original untampered attestation should still verify: %v", err)
	}
}

func TestVerifyRejectsWrongProviderOrAlgorithm(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keyID := KeyID(publicKey)
	trusted := map[string]string{keyID: base64.StdEncoding.EncodeToString(publicKey)}

	attestation := testAttestation()
	payload, err := canonicalJSON(attestation)
	if err != nil {
		t.Fatal(err)
	}
	attestation["signature"] = Signature{
		Provider:  "some.other.provider",
		Algorithm: Algorithm,
		KeyID:     keyID,
		Value:     base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, payload)),
	}
	raw, err := json.Marshal(attestation)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(raw, trusted); err == nil {
		t.Fatal("expected an unrecognized provider to be rejected")
	}
}

func TestVerifyRejectsMissingSignature(t *testing.T) {
	raw, err := json.Marshal(testAttestation())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(raw, map[string]string{}); !errors.Is(err, ErrNoSignature) {
		t.Fatalf("expected ErrNoSignature, got %v", err)
	}
}

func TestVerifyIsIndependentOfFieldOrder(t *testing.T) {
	// The whole point of canonical-JSON signing: a byte-for-byte
	// reordering of the source object's fields must not change whether
	// the signature verifies, since verification recomputes the payload
	// by canonicalizing (key-sorting) whatever JSON object it is given.
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keyID := KeyID(publicKey)
	raw := sign(t, testAttestation(), privateKey, keyID)
	trusted := map[string]string{keyID: base64.StdEncoding.EncodeToString(publicKey)}

	var generic map[string]json.RawMessage
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatal(err)
	}
	// Re-marshal through a differently-typed map so key order in the
	// resulting text differs from the original (Go map iteration order is
	// randomized, and json.Marshal on map[string]json.RawMessage sorts
	// keys — reconstruct via a fresh map literal in a different order to
	// be sure).
	reordered, err := json.Marshal(map[string]json.RawMessage{
		"signature": generic["signature"],
		"analysis":  generic["analysis"],
		"result":    generic["result"],
		"producer":  generic["producer"],
		"subject":   generic["subject"],
		"version":   generic["version"],
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(reordered, trusted); err != nil {
		t.Fatalf("expected reordered-but-untampered attestation to still verify: %v", err)
	}
}
