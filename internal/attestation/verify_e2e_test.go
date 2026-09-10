//go:build e2e

package attestation

// Proves interoperability with a real skil binary, not just this package's
// own idea of skil's signing scheme: it shells out to `skil key generate`
// and `skil attest --signing-key ... --sign` (built from a pinned skil
// commit — see .github/workflows/ci.yml's skil-interop job, mirroring the
// same pattern skcr uses for its skil-interop CI job) and verifies the
// resulting attestation with this package's independent Verify, with zero
// dependency on skil's Go types.
//
// Requires a `skil` binary on PATH (or SKIL_BINARY set to its path);
// skipped otherwise, since this is an environment-dependent interop test,
// not something every local `go test ./...` run needs to gate on.

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func skilBinary(t *testing.T) string {
	t.Helper()
	if bin := os.Getenv("SKIL_BINARY"); bin != "" {
		return bin
	}
	bin, err := exec.LookPath("skil")
	if err != nil {
		t.Fatal("skil binary not found on PATH (set SKIL_BINARY) — skipping skil interop test")
	}
	return bin
}

func TestVerifyInteropsWithRealSkilBinary(t *testing.T) {
	skil := skilBinary(t)
	dir := t.TempDir()

	keyPath := filepath.Join(dir, "key.pem")
	keyOut, err := exec.Command(skil, "key", "generate", "--output", keyPath).Output()
	if err != nil {
		t.Fatalf("skil key generate: %v", exitErr(err))
	}
	var keyInfo struct {
		KeyID     string `json:"key_id"`
		PublicKey string `json:"public_key"`
	}
	if err := json.Unmarshal(keyOut, &keyInfo); err != nil {
		t.Fatalf("parse skil key generate output %q: %v", keyOut, err)
	}

	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "name: demo-skill\nversion: 1.0.0\ndescription: interop test fixture\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\n"+manifest+"---\n\n# demo-skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	attestationPath := filepath.Join(dir, "attestation.json")
	cmd := exec.Command(skil, "attest", skillDir,
		"--signing-key", keyPath,
		"--output", attestationPath,
		"--static-only")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("skil attest: %v\n%s", err, out)
	}

	predicate, err := os.ReadFile(attestationPath)
	if err != nil {
		t.Fatalf("read attestation output: %v", err)
	}

	trusted := map[string]string{keyInfo.KeyID: keyInfo.PublicKey}
	sig, err := Verify(json.RawMessage(predicate), trusted)
	if err != nil {
		t.Fatalf("Verify a real skil-signed attestation: %v", err)
	}
	if sig.KeyID != keyInfo.KeyID {
		t.Fatalf("verified signature KeyID = %q, want %q", sig.KeyID, keyInfo.KeyID)
	}

	if _, err := Verify(json.RawMessage(predicate), map[string]string{}); err == nil {
		t.Fatal("expected verification against an empty trust store to fail")
	}

	tampered := append([]byte(nil), predicate...)
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(tampered, &generic); err != nil {
		t.Fatal(err)
	}
	generic["analysis"] = json.RawMessage(`["tampered"]`)
	tamperedRaw, err := json.Marshal(generic)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(json.RawMessage(tamperedRaw), trusted); err == nil {
		t.Fatal("expected verification of a tampered attestation to fail")
	}
}

func exitErr(err error) error {
	if ee, ok := err.(*exec.ExitError); ok {
		return &wrappedExitErr{err: ee}
	}
	return err
}

type wrappedExitErr struct{ err *exec.ExitError }

func (w *wrappedExitErr) Error() string {
	return w.err.Error() + ": " + string(w.err.Stderr)
}
