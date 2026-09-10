// Package attestation independently verifies the Ed25519 signature on a
// skil-produced attestation (https://github.com/domehahn/skil, `skil attest
// --sign`) directly from its JSON wire form — without importing skil as a
// Go library dependency. skpm already stores an attestation's predicate as
// opaque JSON (see internal/registry.AttestationRequest/Record); this
// package adds the ability to check that predicate's signature against a
// configured set of trusted signing keys.
//
// This mirrors skil's own internal/signing package byte-for-byte:
//   - KeyID(pubkey) = "sha256:" + hex(sha256(pubkey))
//   - the signed payload is a canonical JSON form of the attestation with
//     its "signature" field removed: marshal, then round-trip through a
//     generic decode (json.Number keeps exact numeric literals) and
//     re-marshal, which encoding/json emits with object keys sorted
//     lexicographically. This canonical form is what makes independent
//     verification possible at all — Go's normal struct-order marshal
//     would require mirroring skil's exact internal struct shape
//     field-for-field, which is exactly the hard library coupling this
//     package exists to avoid.
//   - a signature is only accepted if Provider == "builtin.ed25519" and
//     Algorithm equals "Ed25519" (case-insensitively), and its KeyID is
//     present in the caller-supplied trustedKeys map.
package attestation

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	// Provider and Algorithm are the only signature scheme this package
	// understands, matching skil's internal/signing package. A signature
	// naming any other provider or algorithm is rejected rather than
	// silently ignored, since skpm has no way to verify it.
	Provider  = "builtin.ed25519"
	Algorithm = "Ed25519"
)

// Signature mirrors skil's pkg/skil.Signature. It is the only part of a
// skil attestation's shape this package needs to know explicitly — every
// other field is handled generically.
type Signature struct {
	Provider  string `json:"provider"`
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"key_id"`
	Value     string `json:"value"`
}

// KeyID reproduces skil's internal/signing.KeyID: the trust-anchor
// identifier for an Ed25519 public key, used to look it up in a
// trusted-keys map and to match against a Signature's KeyID field.
func KeyID(publicKey ed25519.PublicKey) string {
	sum := sha256.Sum256(publicKey)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ErrNoSignature is returned when the predicate JSON has no top-level
// "signature" field to verify.
var ErrNoSignature = errors.New("attestation has no signature field")

// ErrUntrustedKey is returned when the signature names a KeyID that is not
// present in the supplied trustedKeys map.
var ErrUntrustedKey = errors.New("attestation signed by an untrusted key")

// Verify checks predicate's Ed25519 signature against trustedKeys (a map
// from KeyID to base64-encoded Ed25519 public key — the same shape as
// skil's policy.TrustedSigners config field, and skpm's own
// config.Config.TrustedSigners). It returns the verified Signature on
// success, or a descriptive error otherwise.
//
// predicate is treated as a generic JSON object: this function does not
// require (or benefit from) knowing skil's Attestation/EvidenceBundle/
// PackageStatement struct shapes — any JSON object with a "signature" key
// shaped like Signature works, which is what lets skpm verify skil's
// evidence without depending on skil's Go types.
func Verify(predicate json.RawMessage, trustedKeys map[string]string) (*Signature, error) {
	var envelope map[string]json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(predicate))
	dec.UseNumber()
	if err := dec.Decode(&envelope); err != nil {
		return nil, fmt.Errorf("attestation is not a JSON object: %w", err)
	}

	if envelope == nil {
		return nil, errors.New("attestation must be a JSON object")
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return nil, errors.New("attestation contains trailing JSON data")
	}
	sigRaw, ok := envelope["signature"]
	if !ok {
		return nil, ErrNoSignature
	}
	var sig Signature
	if err := json.Unmarshal(sigRaw, &sig); err != nil {
		return nil, fmt.Errorf("parse signature field: %w", err)
	}
	if sig.Provider != Provider || !strings.EqualFold(sig.Algorithm, Algorithm) {
		return nil, fmt.Errorf("unsupported signature provider or algorithm %q/%q", sig.Provider, sig.Algorithm)
	}

	encodedKey, ok := trustedKeys[sig.KeyID]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUntrustedKey, sig.KeyID)
	}
	publicKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("trusted key %q is not a valid base64 Ed25519 public key", sig.KeyID)
	}
	if KeyID(publicKey) != sig.KeyID {
		return nil, errors.New("signature key ID does not match trusted public key")
	}
	value, err := base64.StdEncoding.DecodeString(sig.Value)
	if err != nil || len(value) != ed25519.SignatureSize {
		return nil, errors.New("signature value is not valid base64 Ed25519 data")
	}

	delete(envelope, "signature")
	payload, err := canonicalJSON(envelope)
	if err != nil {
		return nil, fmt.Errorf("canonicalize attestation payload: %w", err)
	}
	if !ed25519.Verify(publicKey, payload, value) {
		return nil, errors.New("cryptographic signature verification failed")
	}
	return &sig, nil
}

// canonicalJSON reproduces skil's internal/signing.CanonicalJSON: marshal
// v, then decode generically with UseNumber (so numeric literals keep
// their exact original digit sequence) and re-marshal. encoding/json
// marshals map[string]any with keys sorted lexicographically, which is
// what makes this independent of the original field order. Both sides
// (skil signing, skpm verifying) must apply exactly this algorithm to the
// same underlying JSON object for a signature to check out.
func canonicalJSON(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var generic any
	if err := dec.Decode(&generic); err != nil {
		return nil, err
	}
	return json.Marshal(generic)
}
