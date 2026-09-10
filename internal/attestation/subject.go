package attestation

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// BindSubject checks the identity inside the signed predicate, without interpreting
// security findings. Attestations must bind to exact artifact/package digests.
func BindSubject(predicate json.RawMessage, digest string) error {
	cleanDigest := strings.TrimPrefix(digest, "sha256:")
	cleanDigest = strings.ToLower(cleanDigest)
	raw, err := hex.DecodeString(cleanDigest)
	if err != nil || len(raw) != 32 || cleanDigest != digest {
		return fmt.Errorf("ATTESTATION_DIGEST_INVALID: expected lowercase SHA256 hex")
	}

	var doc struct {
		Version int `json:"version"`
		Subject struct {
			SHA256 string `json:"sha256"`
		} `json:"subject"`
	}
	if err := json.Unmarshal(predicate, &doc); err != nil || doc.Version == 0 {
		return fmt.Errorf("ATTESTATION_INVALID: %w", err)
	}
	if doc.Version != 1 {
		return fmt.Errorf("ATTESTATION_SCHEMA_UNSUPPORTED: version %d", doc.Version)
	}
	if strings.TrimPrefix(strings.ToLower(doc.Subject.SHA256), "sha256:") != cleanDigest {
		return fmt.Errorf("ATTESTATION_SUBJECT_MISMATCH: signed subject must identify the exact package")
	}
	return nil
}
