package cli_test

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func buildTestLockfile(t *testing.T) *lockfile.LockFile {
	t.Helper()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "security-reviewer",
		Version:     "1.2.0",
		Source:      "https://registry.example.com",
		DownloadURL: "https://registry.example.com/security-reviewer-1.2.0.zip",
		SHA256:      "aabbccdd" + strings.Repeat("00", 28),
		InstalledTo: []string{".claude/skills/security-reviewer"},
	})
	lf.Upsert(lockfile.SkillLock{
		Name:        "code-formatter",
		Namespace:   "myorg",
		Version:     "0.5.1",
		Source:      "https://registry.example.com",
		SHA256:      "11223344" + strings.Repeat("ff", 28),
		InstalledTo: []string{".github/skills/code-formatter"},
	})
	return lf
}

func testSigningKey() []byte {
	// 32 bytes → 256-bit key.
	return []byte("test-signing-key-32-bytes-pad000")
}

func testSigningKeyB64() string {
	return base64.StdEncoding.EncodeToString(testSigningKey())
}

// ── sbom: SPDX ────────────────────────────────────────────────────────────────

func TestSbomSPDXContainsRequiredHeaders(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomSPDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, "SPDXVersion: SPDX-2.3")
	assert.Contains(t, out, "DataLicense: CC0-1.0")
	assert.Contains(t, out, "Creator: Tool: skpm")
}

func TestSbomSPDXContainsSkillEntries(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomSPDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, "PackageName: security-reviewer")
	assert.Contains(t, out, "PackageVersion: 1.2.0")
	assert.Contains(t, out, "PackageName: code-formatter")
	assert.Contains(t, out, "PackageVersion: 0.5.1")
}

func TestSbomSPDXContainsSHA256Checksums(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomSPDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, "PackageChecksum: SHA256:")
}

func TestSbomSPDXContainsPURL(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomSPDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, "pkg:skpm/security-reviewer@1.2.0")
	// Namespaced skill gets namespace in purl.
	assert.Contains(t, out, "pkg:skpm/myorg/code-formatter@0.5.1")
}

func TestSbomSPDXEmptyLockfile(t *testing.T) {
	lf := lockfile.New()
	out, err := cli.SbomSPDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, "SPDXVersion:")
	// No package sections.
	assert.NotContains(t, out, "PackageName:")
}

// ── sbom: CycloneDX ───────────────────────────────────────────────────────────

func TestSbomCycloneDXValidJSON(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomCycloneDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, `"bomFormat": "CycloneDX"`)
	assert.Contains(t, out, `"specVersion": "1.5"`)
}

func TestSbomCycloneDXContainsComponents(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomCycloneDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, `"name": "security-reviewer"`)
	assert.Contains(t, out, `"version": "1.2.0"`)
	assert.Contains(t, out, `"purl": "pkg:skpm/security-reviewer@1.2.0"`)
}

func TestSbomCycloneDXContainsHashes(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomCycloneDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, `"alg": "SHA-256"`)
}

func TestSbomCycloneDXContainsTool(t *testing.T) {
	lf := buildTestLockfile(t)
	out, err := cli.SbomCycloneDX(lf)
	require.NoError(t, err)
	assert.Contains(t, out, `"name": "skpm"`)
}

// ── sign: LoadSigningKey ──────────────────────────────────────────────────────

func TestLoadSigningKeyFromEnv(t *testing.T) {
	t.Setenv("TEST_SKPM_KEY", testSigningKeyB64())
	key, err := cli.LoadSigningKey("", "TEST_SKPM_KEY")
	require.NoError(t, err)
	assert.Equal(t, testSigningKey(), key)
}

func TestLoadSigningKeyFromFile(t *testing.T) {
	dir := t.TempDir()
	keyFile := filepath.Join(dir, "signing.key")
	require.NoError(t, os.WriteFile(keyFile, []byte(testSigningKeyB64()), 0o600))

	key, err := cli.LoadSigningKey(keyFile, "")
	require.NoError(t, err)
	assert.Equal(t, testSigningKey(), key)
}

func TestLoadSigningKeyFileHasPriority(t *testing.T) {
	t.Setenv("TEST_SKPM_KEY", base64.StdEncoding.EncodeToString([]byte("env-key-that-should-be-ignored-!!")))
	dir := t.TempDir()
	keyFile := filepath.Join(dir, "signing.key")
	require.NoError(t, os.WriteFile(keyFile, []byte(testSigningKeyB64()), 0o600))

	key, err := cli.LoadSigningKey(keyFile, "TEST_SKPM_KEY")
	require.NoError(t, err)
	assert.Equal(t, testSigningKey(), key)
}

func TestLoadSigningKeyMissingEnvAndFile(t *testing.T) {
	t.Setenv("SKPM_SIGNING_KEY", "")
	_, err := cli.LoadSigningKey("", "SKPM_SIGNING_KEY")
	assert.Error(t, err)
}

func TestLoadSigningKeyInvalidBase64(t *testing.T) {
	t.Setenv("TEST_SKPM_KEY", "not-valid-base64!!!")
	_, err := cli.LoadSigningKey("", "TEST_SKPM_KEY")
	assert.Error(t, err)
}

func TestLoadSigningKeyTooShort(t *testing.T) {
	// 8 bytes — below the 16-byte minimum.
	short := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	t.Setenv("TEST_SKPM_KEY", short)
	_, err := cli.LoadSigningKey("", "TEST_SKPM_KEY")
	assert.Error(t, err)
}

// ── sign: SkillSignature ──────────────────────────────────────────────────────

func TestSkillSignatureIsDeterministic(t *testing.T) {
	key := testSigningKey()
	sig1 := cli.SkillSignature(key, "my-skill", "1.0.0", "abc123")
	sig2 := cli.SkillSignature(key, "my-skill", "1.0.0", "abc123")
	assert.Equal(t, sig1, sig2)
}

func TestSkillSignatureChangesWithDifferentInputs(t *testing.T) {
	key := testSigningKey()
	base := cli.SkillSignature(key, "my-skill", "1.0.0", "abc123")
	assert.NotEqual(t, base, cli.SkillSignature(key, "other-skill", "1.0.0", "abc123"))
	assert.NotEqual(t, base, cli.SkillSignature(key, "my-skill", "2.0.0", "abc123"))
	assert.NotEqual(t, base, cli.SkillSignature(key, "my-skill", "1.0.0", "different"))
}

func TestSkillSignatureChangesWithDifferentKey(t *testing.T) {
	key1 := testSigningKey()
	key2 := []byte("different-key-32-bytes-padded000")
	sig1 := cli.SkillSignature(key1, "my-skill", "1.0.0", "abc123")
	sig2 := cli.SkillSignature(key2, "my-skill", "1.0.0", "abc123")
	assert.NotEqual(t, sig1, sig2)
}

func TestSkillSignatureIsValidBase64(t *testing.T) {
	sig := cli.SkillSignature(testSigningKey(), "my-skill", "1.0.0", "digest")
	_, err := base64.StdEncoding.DecodeString(sig)
	assert.NoError(t, err, "signature should be valid base64")
}

// ── sign: apply + verify round-trip ──────────────────────────────────────────

func TestSignApplyAndVerifyRoundTrip(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, "agent-skills.lock")
	lf := buildTestLockfile(t)
	require.NoError(t, lf.Write(lockPath))

	key := testSigningKey()

	// Apply signatures.
	for i, sl := range lf.Skills {
		lf.Skills[i].Signature = cli.SkillSignature(key, sl.Name, sl.Version, sl.SHA256)
	}
	require.NoError(t, lf.Write(lockPath))

	// Reload and verify each signature.
	loaded, err := lockfile.Read(lockPath)
	require.NoError(t, err)
	for _, sl := range loaded.Skills {
		require.NotEmpty(t, sl.Signature, "skill %s should have signature", sl.Name)
		expected := cli.SkillSignature(key, sl.Name, sl.Version, sl.SHA256)
		assert.Equal(t, expected, sl.Signature, "signature mismatch for %s", sl.Name)
	}
}

func TestSignatureDetectsTampering(t *testing.T) {
	key := testSigningKey()
	original := cli.SkillSignature(key, "my-skill", "1.0.0", "abc123")
	tampered := cli.SkillSignature(key, "my-skill", "1.0.0", "TAMPERED")
	assert.NotEqual(t, original, tampered)
}

// ── rollback: SnapshotPath ────────────────────────────────────────────────────

func TestSnapshotPathFormat(t *testing.T) {
	p := cli.SnapshotPath("my-snap")
	assert.Equal(t, filepath.Join(".skpm-snapshots", "my-snap.lock"), p)
}

// ── ci: command registered and flags parse ────────────────────────────────────

func TestCICommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "ci" {
			found = true
			break
		}
	}
	assert.True(t, found, "ci command should be registered on root")
}

func TestCICommandHasExpectedFlags(t *testing.T) {
	root := cli.NewRootCmd()
	var ciCmd interface {
		Flags() interface{ Lookup(string) interface{} }
	}
	_ = ciCmd
	for _, sub := range root.Commands() {
		if sub.Use == "ci" {
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("platform"))
			assert.NotNil(t, sub.Flags().Lookup("no-verify"))
			break
		}
	}
}

// ── sbom: command registered ──────────────────────────────────────────────────

func TestSBOMCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "sbom" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("format"))
			assert.NotNil(t, sub.Flags().Lookup("out"))
			break
		}
	}
	assert.True(t, found, "sbom command should be registered on root")
}

// ── sign: command registered ──────────────────────────────────────────────────

func TestSignCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "sign" {
			found = true
			subNames := []string{}
			for _, s := range sub.Commands() {
				subNames = append(subNames, s.Use)
			}
			assert.Contains(t, subNames, "apply")
			assert.Contains(t, subNames, "verify")
			break
		}
	}
	assert.True(t, found, "sign command should be registered on root")
}

// ── rollback: command registered ─────────────────────────────────────────────

func TestRollbackCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "rollback <skill-name>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("to"))
			assert.NotNil(t, sub.Flags().Lookup("list"))
			break
		}
	}
	assert.True(t, found, "rollback command should be registered on root")
}
