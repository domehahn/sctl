package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── humanBytes ───────────────────────────────────────────────────────────

func TestHumanBytes(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1024 * 1024, "1.0 MiB"},
		{1024 * 1024 * 1024, "1.0 GiB"},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.want, cli.HumanBytes(tc.in), "input=%d", tc.in)
	}
}

// ── parseSkillAtVersion ──────────────────────────────────────────────────

func TestParseSkillAtVersion(t *testing.T) {
	name, ver, err := cli.ParseSkillAtVersion("my-skill@1.2.3")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", name)
	assert.Equal(t, "1.2.3", ver)
}

func TestParseSkillAtVersionMissingVersion(t *testing.T) {
	_, _, err := cli.ParseSkillAtVersion("my-skill")
	assert.Error(t, err)
}

func TestParseSkillAtVersionEmptyName(t *testing.T) {
	_, _, err := cli.ParseSkillAtVersion("@1.0.0")
	assert.Error(t, err)
}

func TestParseSkillAtVersionNamespaced(t *testing.T) {
	name, ver, err := cli.ParseSkillAtVersion("org/skill@2.0.0")
	require.NoError(t, err)
	assert.Equal(t, "org/skill", name)
	assert.Equal(t, "2.0.0", ver)
}

// ── parseSkillRef (clone) ─────────────────────────────────────────────────

func TestParseSkillRef(t *testing.T) {
	name, ver, err := cli.ParseSkillRef("my-skill@1.0.0")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", name)
	assert.Equal(t, "1.0.0", ver)
}

func TestParseSkillRefNoVersion(t *testing.T) {
	name, ver, err := cli.ParseSkillRef("my-skill")
	require.NoError(t, err)
	assert.Equal(t, "my-skill", name)
	assert.Equal(t, "", ver)
}

// ── scaffoldFiles (create) ────────────────────────────────────────────────

func TestScaffoldFilesContainsRequiredFiles(t *testing.T) {
	files := cli.ScaffoldFiles("my-skill", "does something", "MIT", "default", []string{"claude-code"})
	assert.Contains(t, files, "SKILL.md")
	assert.Contains(t, files, "skill.yaml")
	assert.Contains(t, files, "VERSION")
	assert.Contains(t, files, "CHANGELOG.md")
}

func TestScaffoldFilesSkillYAMLContent(t *testing.T) {
	files := cli.ScaffoldFiles("my-skill", "test desc", "Apache-2.0", "default", []string{"claude-code"})
	yaml := files["skill.yaml"]
	assert.Contains(t, yaml, "name: my-skill")
	assert.Contains(t, yaml, "version: 0.1.0")
	assert.Contains(t, yaml, "description: test desc")
	assert.Contains(t, yaml, "license: Apache-2.0")
	assert.Contains(t, yaml, "claude-code")
}

func TestScaffoldFilesVERSION(t *testing.T) {
	files := cli.ScaffoldFiles("x", "", "MIT", "default", nil)
	assert.Equal(t, "0.1.0\n", files["VERSION"])
}

func TestScaffoldFilesSKILLMDHasFrontmatter(t *testing.T) {
	files := cli.ScaffoldFiles("my-skill", "desc", "MIT", "default", []string{"all"})
	md := files["SKILL.md"]
	assert.True(t, strings.HasPrefix(md, "---\n"), "SKILL.md should start with YAML frontmatter")
	assert.Contains(t, md, "name: my-skill")
}

// ── parseChangelogEntries ─────────────────────────────────────────────────

func TestParseChangelogEntries(t *testing.T) {
	content := `# Changelog

## 1.2.0

- added thing

## 1.1.0

- fixed bug
`
	entries := cli.ParseChangelogEntries(content)
	require.Len(t, entries, 2)
	assert.Contains(t, entries[0], "1.2.0")
	assert.Contains(t, entries[1], "1.1.0")
}

func TestParseChangelogEntriesEmpty(t *testing.T) {
	entries := cli.ParseChangelogEntries("# Changelog\n\nNo entries yet.\n")
	assert.Empty(t, entries)
}

// ── addChangelogMessage ───────────────────────────────────────────────────

func TestAddChangelogMessageCreatesSection(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte("# Changelog\n"), 0o644))

	require.NoError(t, cli.AddChangelogMessage(dir, "1.0.0", "initial release"))

	data, _ := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	content := string(data)
	assert.Contains(t, content, "## 1.0.0")
	assert.Contains(t, content, "- initial release")
}

func TestAddChangelogMessageAppendsToExistingSection(t *testing.T) {
	dir := t.TempDir()
	initial := "# Changelog\n\n## 1.0.0\n\n- first entry\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte(initial), 0o644))

	require.NoError(t, cli.AddChangelogMessage(dir, "1.0.0", "second entry"))

	data, _ := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	content := string(data)
	assert.Contains(t, content, "- first entry")
	assert.Contains(t, content, "- second entry")
}

func TestAddChangelogMessageDoesNotDuplicateSection(t *testing.T) {
	dir := t.TempDir()
	initial := "# Changelog\n\n## 1.0.0\n\n- existing\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte(initial), 0o644))

	require.NoError(t, cli.AddChangelogMessage(dir, "1.0.0", "new entry"))

	data, _ := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	// Only one ## 1.0.0 section.
	assert.Equal(t, 1, strings.Count(string(data), "## 1.0.0"))
}

// ── snapshotPath ─────────────────────────────────────────────────────────

func TestSnapshotPath(t *testing.T) {
	p := cli.SnapshotPath("my-snap")
	assert.Equal(t, filepath.Join(".skpm-snapshots", "my-snap.lock"), p)
}

// ── links manifest round-trip ─────────────────────────────────────────────

func TestLinksManifestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	require.NoError(t, cli.UpsertLink("my-skill", "/abs/path", []string{".claude/skills/my-skill"}))

	linked, path := cli.IsLinkedSkill("my-skill")
	assert.True(t, linked)
	assert.Equal(t, "/abs/path", path)

	notLinked, _ := cli.IsLinkedSkill("other-skill")
	assert.False(t, notLinked)
}

// ── defaultSkillRoots ─────────────────────────────────────────────────────

func TestDefaultSkillRootsContainsExpected(t *testing.T) {
	roots := cli.DefaultSkillRoots()
	assert.Contains(t, roots, ".agents/skills")
	assert.Contains(t, roots, ".claude/skills")
	assert.Contains(t, roots, ".github/skills")
}

// ── configGetScalar / configSetScalar ────────────────────────────────────

func TestConfigGetScalar(t *testing.T) {
	cfg := cli.TestConfig()
	cfg.DefaultRegistry = "my-reg"
	cfg.CacheDir = "/tmp/cache"
	cfg.LogLevel = "debug"
	cfg.Concurrency = 8

	assert.Equal(t, "my-reg", cli.ConfigGetScalar(cfg, "default_registry"))
	assert.Equal(t, "/tmp/cache", cli.ConfigGetScalar(cfg, "cache_dir"))
	assert.Equal(t, "debug", cli.ConfigGetScalar(cfg, "log_level"))
	assert.Equal(t, "8", cli.ConfigGetScalar(cfg, "concurrency"))
	assert.Equal(t, "", cli.ConfigGetScalar(cfg, "unknown_key"))
}

func TestConfigSetScalar(t *testing.T) {
	cfg := cli.TestConfig()
	require.NoError(t, cli.ConfigSetScalar(cfg, "default_registry", "new-reg"))
	assert.Equal(t, "new-reg", cfg.DefaultRegistry)

	require.NoError(t, cli.ConfigSetScalar(cfg, "concurrency", "12"))
	assert.Equal(t, 12, cfg.Concurrency)

	err := cli.ConfigSetScalar(cfg, "concurrency", "not-a-number")
	assert.Error(t, err)

	err = cli.ConfigSetScalar(cfg, "unknown", "value")
	assert.Error(t, err)
}

// ── semverDelta ───────────────────────────────────────────────────────────

func TestSemverDelta(t *testing.T) {
	assert.Equal(t, "major", cli.SemverDelta("1.0.0", "2.0.0"))
	assert.Equal(t, "minor", cli.SemverDelta("1.0.0", "1.1.0"))
	assert.Equal(t, "patch", cli.SemverDelta("1.0.0", "1.0.1"))
	assert.Equal(t, "patch", cli.SemverDelta("1.2.3", "1.2.4"))
}

// ── isWithinBase ─────────────────────────────────────────────────────────

func TestIsWithinBase(t *testing.T) {
	assert.True(t, cli.IsWithinBase("/base", "/base/sub/file"))
	assert.True(t, cli.IsWithinBase("/base", "/base/file"))
	assert.False(t, cli.IsWithinBase("/base", "/base/../other"))
	assert.False(t, cli.IsWithinBase("/base", "/other/file"))
}
