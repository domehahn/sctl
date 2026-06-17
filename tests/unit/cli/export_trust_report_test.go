package cli_test

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── export: command registered ────────────────────────────────────────────────

func TestExportCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "export" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("out"))
			assert.NotNil(t, sub.Flags().Lookup("include-snapshots"))
			break
		}
	}
	assert.True(t, found, "export command should be registered on root")
}

func TestExportCreatesArchive(t *testing.T) {
	dir := t.TempDir()

	// Create a fake installed skill directory.
	skillDir := filepath.Join(dir, ".claude", "skills", "my-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# My Skill"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "skill.yaml"), []byte("name: my-skill\n"), 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	outPath := filepath.Join(dir, "bundle.tar.gz")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"export", "--lock", lockPath, "--out", outPath})
	require.NoError(t, root.Execute())

	// Verify the archive exists and contains expected files.
	f, err := os.Open(outPath)
	require.NoError(t, err)
	defer f.Close()

	gr, err := gzip.NewReader(f)
	require.NoError(t, err)
	tr := tar.NewReader(gr)

	var names []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		names = append(names, hdr.Name)
	}

	assert.Contains(t, names, "agent-skills.lock")
	hasSkillMD := false
	for _, n := range names {
		if strings.Contains(n, "SKILL.md") {
			hasSkillMD = true
			break
		}
	}
	assert.True(t, hasSkillMD, "archive should contain SKILL.md from installed skill")
}

func TestExportDryRunDoesNotCreateFile(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	outPath := filepath.Join(dir, "bundle.tar.gz")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"export", "--dry-run", "--lock", lockPath, "--out", outPath})
	require.NoError(t, root.Execute())

	_, err := os.Stat(outPath)
	assert.True(t, os.IsNotExist(err), "dry-run should not create archive file")
}

// ── trust: LoadTrustStore ─────────────────────────────────────────────────────

func TestLoadTrustStoreNonExistentReturnsEmpty(t *testing.T) {
	ts, err := cli.LoadTrustStore("/nonexistent/.skpm-trust.yaml")
	require.NoError(t, err)
	assert.Empty(t, ts.Keys)
	assert.Equal(t, 1, ts.Version)
}

func TestLoadTrustStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	trustPath := filepath.Join(dir, ".skpm-trust.yaml")

	// Use trust add command to write an entry.
	t.Setenv("TEST_TRUST_KEY", testSigningKeyB64())
	root := cli.NewRootCmd()
	root.SetArgs([]string{"trust", "add", "--trust-file", trustPath, "--name", "ci-signer", "--key-env", "TEST_TRUST_KEY"})
	require.NoError(t, root.Execute())

	ts, err := cli.LoadTrustStore(trustPath)
	require.NoError(t, err)
	require.Len(t, ts.Keys, 1)
	assert.Equal(t, "ci-signer", ts.Keys[0].Name)
	assert.Equal(t, testSigningKeyB64(), ts.Keys[0].Key)
	assert.NotEmpty(t, ts.Keys[0].Added)
}

func TestTrustAddRejectsDuplicateName(t *testing.T) {
	dir := t.TempDir()
	trustPath := filepath.Join(dir, ".skpm-trust.yaml")
	t.Setenv("TEST_TRUST_KEY", testSigningKeyB64())

	root := cli.NewRootCmd()
	root.SetArgs([]string{"trust", "add", "--trust-file", trustPath, "--name", "dup", "--key-env", "TEST_TRUST_KEY"})
	require.NoError(t, root.Execute())

	root2 := cli.NewRootCmd()
	root2.SetArgs([]string{"trust", "add", "--trust-file", trustPath, "--name", "dup", "--key-env", "TEST_TRUST_KEY"})
	assert.Error(t, root2.Execute())
}

func TestTrustAddRejectsShortKey(t *testing.T) {
	dir := t.TempDir()
	trustPath := filepath.Join(dir, ".skpm-trust.yaml")
	import64 := "dG9vc2hvcnQ=" // "tooshort" in base64 — only 8 bytes

	root := cli.NewRootCmd()
	root.SetArgs([]string{"trust", "add", "--trust-file", trustPath, "--name", "short", "--key-value", import64})
	assert.Error(t, root.Execute())
}

func TestTrustRevokeRemovesEntry(t *testing.T) {
	dir := t.TempDir()
	trustPath := filepath.Join(dir, ".skpm-trust.yaml")
	t.Setenv("TEST_TRUST_KEY", testSigningKeyB64())

	root := cli.NewRootCmd()
	root.SetArgs([]string{"trust", "add", "--trust-file", trustPath, "--name", "to-remove", "--key-env", "TEST_TRUST_KEY"})
	require.NoError(t, root.Execute())

	root2 := cli.NewRootCmd()
	root2.SetArgs([]string{"trust", "revoke", "--trust-file", trustPath, "--name", "to-remove"})
	require.NoError(t, root2.Execute())

	ts, err := cli.LoadTrustStore(trustPath)
	require.NoError(t, err)
	assert.Empty(t, ts.Keys)
}

func TestTrustRevokeFailsForUnknownName(t *testing.T) {
	dir := t.TempDir()
	trustPath := filepath.Join(dir, ".skpm-trust.yaml")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"trust", "revoke", "--trust-file", trustPath, "--name", "ghost"})
	assert.Error(t, root.Execute())
}

// ── trust: command registered ─────────────────────────────────────────────────

func TestTrustCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "trust" {
			found = true
			subNames := make(map[string]bool)
			for _, s := range sub.Commands() {
				subNames[s.Use] = true
			}
			assert.True(t, subNames["list"])
			assert.True(t, subNames["add"])
			assert.True(t, subNames["revoke"])
			break
		}
	}
	assert.True(t, found, "trust command should be registered on root")
}

// ── report: command registered and runs ───────────────────────────────────────

func TestReportCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "report" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("policy"))
			assert.NotNil(t, sub.Flags().Lookup("key-file"))
			assert.NotNil(t, sub.Flags().Lookup("fail"))
			break
		}
	}
	assert.True(t, found, "report command should be registered on root")
}

func TestReportRunsWithEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"report", "--lock", lockPath})
	// Should not error — just report skipped/ok sections.
	assert.NoError(t, root.Execute())
}

func TestReportSkipsSignaturesWithoutKey(t *testing.T) {
	dir := t.TempDir()
	lf := buildTestLockfile(t)
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetArgs([]string{"report", "--lock", lockPath})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "SIGNATURES")
	assert.Contains(t, buf.String(), "--key-file or --key-env")
}

func TestReportSkipsPolicyWhenNoFile(t *testing.T) {
	dir := t.TempDir()
	lf := buildTestLockfile(t)
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetArgs([]string{"report", "--lock", lockPath, "--policy", filepath.Join(dir, "no-such-policy.yaml")})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "POLICY")
	assert.Contains(t, buf.String(), "no policy file found")
}

func TestReportDetectsPolicyViolations(t *testing.T) {
	dir := t.TempDir()
	lf := buildTestLockfile(t)
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	// Policy that bans the first skill.
	policyPath := filepath.Join(dir, "skpm-policy.yaml")
	require.NoError(t, os.WriteFile(policyPath, []byte(
		"version: 1\nrules:\n  banned_skills:\n    - security-reviewer\n"), 0o644))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetArgs([]string{"report", "--lock", lockPath, "--policy", policyPath})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "POLICY")
	assert.Contains(t, buf.String(), "banned")
}
