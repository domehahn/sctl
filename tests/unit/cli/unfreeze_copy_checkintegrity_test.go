package cli_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── unfreeze ──────────────────────────────────────────────────────────────────

func TestUnfreezeCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "unfreeze" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
		}
	}
	assert.True(t, found, "unfreeze command should be registered on root")
}

func TestUnfreezeClearsFrozenVersions(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-a",
		Version: "2.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Manifest version matches locked version exactly (frozen state).
	mf := "version: 1\nskills:\n  - name: skill-a\n    version: 2.0.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unfreeze", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "skill-a")

	data, err := os.ReadFile(filepath.Join(dir, "agent-skills.yaml"))
	require.NoError(t, err)
	// Version constraint should be removed (empty) so upgrades work again.
	assert.NotContains(t, string(data), "2.0.0")
}

func TestUnfreezeNoFrozenSkills(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-b",
		Version: "3.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("b", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Manifest version differs — not frozen by 'skpm freeze'.
	mf := "version: 1\nskills:\n  - name: skill-b\n    version: ^3.0.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unfreeze", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No frozen")
}

func TestUnfreezeDryRun(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-c",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("c", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	mf := "version: 1\nskills:\n  - name: skill-c\n    version: 1.0.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unfreeze", "--lock", lp, "--dry-run"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Dry run")

	// Manifest should be unchanged.
	data, _ := os.ReadFile(filepath.Join(dir, "agent-skills.yaml"))
	assert.Contains(t, string(data), "1.0.0")
}

func TestUnfreezeJSONOutput(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-d",
		Version: "5.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("d", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	mf := "version: 1\nskills:\n  - name: skill-d\n    version: 5.0.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unfreeze", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var results []map[string]string
	require.NoError(t, json.Unmarshal(buf.Bytes(), &results))
	assert.Len(t, results, 1)
	assert.Equal(t, "skill-d", results[0]["skill"])
}

// ── copy ──────────────────────────────────────────────────────────────────────

func TestCopyCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "copy <source> <new-name>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
		}
	}
	assert.True(t, found, "copy command should be registered on root")
}

func TestCopyRequiresTwoArgs(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"copy", "only-one"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestCopyCreatesIndependentEntry(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "original",
		Version: "1.2.3",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"copy", "original", "replica", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "replica")

	result, err := lockfile.Read(lp)
	require.NoError(t, err)

	orig, ok := result.Find("original")
	require.True(t, ok)
	replica, ok := result.Find("replica")
	require.True(t, ok)

	assert.Equal(t, orig.Version, replica.Version)
	assert.Equal(t, orig.SHA256, replica.SHA256)
	assert.Equal(t, "replica", replica.Name)
}

func TestCopyErrorsOnMissingSource(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"copy", "ghost", "copy-ghost", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestCopyErrorsOnExistingDest(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "a", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "b", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("b", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"copy", "a", "b", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestCopyDryRun(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "src", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"copy", "src", "dst", "--lock", lp, "--dry-run"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Would copy")

	// dst should not be in the lockfile.
	result, _ := lockfile.Read(lp)
	_, found := result.Find("dst")
	assert.False(t, found)
}

// ── check-integrity ───────────────────────────────────────────────────────────

func TestCheckIntegrityCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "check-integrity" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("fail-on-mismatch"))
		}
	}
	assert.True(t, found, "check-integrity command should be registered on root")
}

func TestCheckIntegritySkipsNoSHA256(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "unsigned",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		// No SHA256
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"check-integrity", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "skip")
}

func TestCheckIntegrityOK(t *testing.T) {
	dir := t.TempDir()
	content := []byte("---\nname: good-skill\nversion: 1.0.0\n---\n\n# Good\n")
	h := sha256.Sum256(content)
	digest := hex.EncodeToString(h[:])

	skillDir := filepath.Join(dir, ".agents", "skills", "good-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), content, 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "good-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      digest,
		InstalledTo: []string{skillDir},
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"check-integrity", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "ok")
	assert.Contains(t, buf.String(), "All checked skills OK")
}

func TestCheckIntegrityDetectsMismatch(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, ".agents", "skills", "tampered")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("tampered content"), 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "tampered",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64), // wrong digest
		InstalledTo: []string{skillDir},
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"check-integrity", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute()) // no fail-on-mismatch, should succeed
	assert.Contains(t, buf.String(), "MISMATCH")
}

func TestCheckIntegrityFailOnMismatch(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, ".agents", "skills", "bad")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("bad content"), 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "bad",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("b", 64),
		InstalledTo: []string{skillDir},
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"check-integrity", "--lock", lp, "--fail-on-mismatch"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestCheckIntegrityJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "no-hash",
		Version: "1.0.0",
		Source:  "https://r.example.com",
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"check-integrity", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var results []map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &results))
	assert.NotEmpty(t, results)
	assert.Contains(t, results[0]["status"], "skip")
}
