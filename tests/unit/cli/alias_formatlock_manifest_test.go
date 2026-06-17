package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── alias: command registered ─────────────────────────────────────────────────

func TestAliasCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "alias <source> <alias>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			break
		}
	}
	assert.True(t, found, "alias command should be registered on root")
}

// ── alias: creates entry ──────────────────────────────────────────────────────

func TestAliasCreatesLockEntry(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "original",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"alias", "original", "my-alias", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)

	entry, ok := updated.Find("my-alias")
	require.True(t, ok, "alias entry should exist in lockfile")
	assert.Equal(t, "my-alias", entry.Name)
	assert.Equal(t, "1.0.0", entry.Version)
	assert.Equal(t, strings.Repeat("a", 64), entry.SHA256)
}

// ── alias: error on missing source ────────────────────────────────────────────

func TestAliasMissingSourceErrors(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"alias", "nonexistent", "my-alias", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

// ── alias: error on duplicate alias ──────────────────────────────────────────

func TestAliasDuplicateErrors(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	sha := strings.Repeat("b", 64)
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lf.Upsert(lockfile.SkillLock{Name: "skill-b", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"alias", "skill-a", "skill-b", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

// ── alias: dry-run makes no changes ──────────────────────────────────────────

func TestAliasDryRun(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "src", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("c", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))
	before, _ := os.ReadFile(lp)

	root := cli.NewRootCmd()
	root.SetArgs([]string{"alias", "src", "new-alias", "--lock", lp, "--dry-run"})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	after, _ := os.ReadFile(lp)
	assert.Equal(t, string(before), string(after))
}

// ── format-lock: command registered ──────────────────────────────────────────

func TestFormatLockCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "format-lock" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("check"))
			break
		}
	}
	assert.True(t, found, "format-lock command should be registered on root")
}

// ── format-lock: sorts entries ────────────────────────────────────────────────

func TestFormatLockSortsEntries(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "zebra", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("z", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "alpha", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"format-lock", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	require.Equal(t, 2, len(updated.Skills))
	assert.Equal(t, "alpha", updated.Skills[0].Name)
	assert.Equal(t, "zebra", updated.Skills[1].Name)
}

// ── format-lock: already canonical ───────────────────────────────────────────

func TestFormatLockAlreadyCanonical(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "alpha", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lf.Sort()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"format-lock", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "canonical")
}

// ── format-lock: --check exits non-zero when not canonical ───────────────────

func TestFormatLockCheckFailsWhenNotCanonical(t *testing.T) {
	dir := t.TempDir()
	// Write entries out of order by manipulating the file directly.
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "zebra", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("z", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "alpha", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Force the file to be non-canonical by writing zebra first.
	// The lockfile.Write already sorts, so we need to write raw YAML.
	raw := "version: 1\nskills:\n  - name: zebra\n    version: 1.0.0\n    source: https://r.example.com\n    sha256: " + strings.Repeat("z", 64) + "\n  - name: alpha\n    version: 1.0.0\n    source: https://r.example.com\n    sha256: " + strings.Repeat("a", 64) + "\n"
	require.NoError(t, os.WriteFile(lp, []byte(raw), 0o644))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"format-lock", "--lock", lp, "--check"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

// ── manifest: command registered ──────────────────────────────────────────────

func TestManifestCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "manifest" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("format"))
			assert.NotNil(t, sub.Flags().Lookup("out"))
			break
		}
	}
	assert.True(t, found, "manifest command should be registered on root")
}

// ── manifest: JSON output ──────────────────────────────────────────────────────

func TestManifestJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-one",
		Version: "1.2.3",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"manifest", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, `"skill-one"`)
	assert.Contains(t, out, `"1.2.3"`)
	assert.Contains(t, out, `"generated_at"`)
}

// ── manifest: YAML output ──────────────────────────────────────────────────────

func TestManifestYAMLOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-yaml",
		Version: "2.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("b", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"manifest", "--lock", lp, "--format", "yaml"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "skill-yaml")
	assert.Contains(t, out, "2.0.0")
	assert.Contains(t, out, "generated_at")
}

// ── manifest: writes to file with --out ──────────────────────────────────────

func TestManifestWritesToFile(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "file-skill",
		Version: "3.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("c", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	outFile := filepath.Join(dir, "manifest.json")
	root := cli.NewRootCmd()
	root.SetArgs([]string{"manifest", "--lock", lp, "--out", outFile})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "file-skill")
}

// ── manifest: empty lockfile ──────────────────────────────────────────────────

func TestManifestEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"manifest", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "generated_at")
}
