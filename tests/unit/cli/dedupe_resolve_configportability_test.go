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

// ── dedupe: command registered ────────────────────────────────────────────────

func TestDedupeCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "dedupe" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("fix"))
			break
		}
	}
	assert.True(t, found, "dedupe command should be registered on root")
}

// ── dedupe: no duplicates ─────────────────────────────────────────────────────

func TestDedupeNoDuplicates(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "skill-b", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("b", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"dedupe", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No duplicate")
}

// ── dedupe: detects duplicates ────────────────────────────────────────────────

func TestDedupeDetectsDuplicates(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("c", 64)
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lf.Upsert(lockfile.SkillLock{Name: "skill-alias", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"dedupe", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "duplicate")
	assert.Contains(t, out, "--fix")
}

// ── dedupe: fix removes duplicate ────────────────────────────────────────────

func TestDedupeFixRemovesDuplicate(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("d", 64)
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "alpha", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lf.Upsert(lockfile.SkillLock{Name: "beta", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"dedupe", "--lock", lp, "--fix"})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	updated, err := lockfile.Read(lp)
	require.NoError(t, err)
	assert.Equal(t, 1, len(updated.Skills), "one duplicate should be removed")
}

// ── dedupe: JSON output ───────────────────────────────────────────────────────

func TestDedupeJSONOutput(t *testing.T) {
	dir := t.TempDir()
	sha := strings.Repeat("e", 64)
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "x", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lf.Upsert(lockfile.SkillLock{Name: "y", Version: "1.0.0", Source: "https://r.example.com", SHA256: sha})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"dedupe", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), `"SHA256"`)
}

// ── resolve: command registered ──────────────────────────────────────────────

func TestResolveCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "resolve <name>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			break
		}
	}
	assert.True(t, found, "resolve command should be registered on root")
}

// ── resolve: missing skill ───────────────────────────────────────────────────

func TestResolveMissingSkillErrors(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"resolve", "nonexistent", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

// ── resolve: shows skill details ─────────────────────────────────────────────

func TestResolveShowsSkillDetails(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "2.1.0",
		Source:      "https://registry.example.com/my-skill",
		SHA256:      strings.Repeat("f", 64),
		InstalledTo: []string{".agents/skills/my-skill"},
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"resolve", "my-skill", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "my-skill")
	assert.Contains(t, out, "2.1.0")
	assert.Contains(t, out, "registry.example.com")
	assert.Contains(t, out, ".agents/skills/my-skill")
	assert.Contains(t, out, "(none)") // signature and provenance
}

// ── resolve: JSON output ──────────────────────────────────────────────────────

func TestResolveJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "json-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"resolve", "json-skill", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), `"json-skill"`)
	assert.Contains(t, buf.String(), `"version"`)
}

// ── config export/import: commands registered ────────────────────────────────

func TestConfigExportImportRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	names := map[string]bool{}
	for _, sub := range root.Commands() {
		if sub.Use == "config" {
			for _, s := range sub.Commands() {
				names[s.Use] = true
			}
			break
		}
	}
	assert.True(t, names["export"], "config export should be registered")
	assert.True(t, names["import"], "config import should be registered")
}

// ── config export: stdout when no --out ──────────────────────────────────────

func TestConfigExportPrintsToStdout(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)
	// Write a minimal config.
	skpmDir := filepath.Join(cfgDir, "skpm")
	require.NoError(t, os.MkdirAll(skpmDir, 0o755))
	cfg := "default_registry: main\nregistries:\n  main:\n    type: http\n    url: https://r.example.com\n"
	require.NoError(t, os.WriteFile(filepath.Join(skpmDir, "config.yaml"), []byte(cfg), 0o600))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"config", "export"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "main")
	assert.Contains(t, out, "r.example.com")
	assert.NotContains(t, out, "token") // no token set
}

// ── config export: writes file when --out given ──────────────────────────────

func TestConfigExportWritesFile(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)
	skpmDir := filepath.Join(cfgDir, "skpm")
	require.NoError(t, os.MkdirAll(skpmDir, 0o755))
	cfg := "registries:\n  r1:\n    type: http\n    url: https://r1.example.com\n"
	require.NoError(t, os.WriteFile(filepath.Join(skpmDir, "config.yaml"), []byte(cfg), 0o600))

	outFile := filepath.Join(t.TempDir(), "exported.yaml")
	root := cli.NewRootCmd()
	root.SetArgs([]string{"config", "export", "--out", outFile})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "r1.example.com")
}

// ── config import: round-trip ────────────────────────────────────────────────

func TestConfigImportRoundTrip(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	// Create export file
	exportFile := filepath.Join(t.TempDir(), "export.yaml")
	exportContent := "default_registry: imported\nregistries:\n  imported:\n    url: https://imported.example.com\n    insecure_skip_verify: false\n"
	require.NoError(t, os.WriteFile(exportFile, []byte(exportContent), 0o600))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"config", "import", "--in", exportFile})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	cfgPath := filepath.Join(cfgDir, "skpm", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "imported")
	assert.Contains(t, string(data), "imported.example.com")
}

// ── config import: errors without --in ───────────────────────────────────────

func TestConfigImportRequiresInFlag(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)

	root := cli.NewRootCmd()
	root.SetArgs([]string{"config", "import"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

// ── config import: --merge preserves existing entries ────────────────────────

func TestConfigImportMergePreservesExisting(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfgDir)
	skpmDir := filepath.Join(cfgDir, "skpm")
	require.NoError(t, os.MkdirAll(skpmDir, 0o755))
	existing := "registries:\n  existing:\n    type: http\n    url: https://existing.example.com\n"
	require.NoError(t, os.WriteFile(filepath.Join(skpmDir, "config.yaml"), []byte(existing), 0o600))

	importFile := filepath.Join(t.TempDir(), "import.yaml")
	require.NoError(t, os.WriteFile(importFile, []byte("registries:\n  new:\n    url: https://new.example.com\n"), 0o600))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"config", "import", "--in", importFile, "--merge"})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	cfgPath := filepath.Join(cfgDir, "skpm", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "existing.example.com")
	assert.Contains(t, content, "new.example.com")
}
