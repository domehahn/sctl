package cli_test

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── sources command ───────────────────────────────────────────────────────────

func TestSourcesGroupsByURL(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "https://registry.example.com", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "lint-skill", Version: "1.0.0", Source: "https://registry.example.com", SHA256: strings.Repeat("b", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "scan-skill", Version: "1.0.0", Source: "https://other.example.com", SHA256: strings.Repeat("c", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"sources", "--lock", lp})
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "registry.example.com")
	assert.Contains(t, out, "other.example.com")
	// registry.example.com has 2 skills — should appear first (sorted by count desc).
	idxRegistry := strings.Index(out, "registry.example.com")
	idxOther := strings.Index(out, "other.example.com")
	assert.Less(t, idxRegistry, idxOther, "higher-count source should appear first")
}

func TestSourcesVerboseLists(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "https://registry.example.com", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"sources", "--lock", lp, "--list-skills"})
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "auth-skill")
}

func TestSourcesJSON(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", Source: "https://registry.example.com", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"sources", "--lock", lp, "--output", "json"})
	require.NoError(t, root.Execute())

	var result []map[string]any
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &result))
	require.Len(t, result, 1)
	assert.Equal(t, "https://registry.example.com", result[0]["source"])
	assert.Equal(t, float64(1), result[0]["count"])
}

func TestSourcesUnknownSource(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "auth-skill", Version: "1.0.0", SHA256: strings.Repeat("a", 64)})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"sources", "--lock", lp})
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "(unknown)")
}

func TestSourcesEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"sources", "--lock", lp})
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "No skills")
}

// ── installed-to command ──────────────────────────────────────────────────────

func TestInstalledToFindsSkill(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64),
		InstalledTo: []string{".agents/skills", ".claude/skills"},
	})
	lf.Upsert(lockfile.SkillLock{
		Name: "lint-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("b", 64),
		InstalledTo: []string{".claude/skills"},
	})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"installed-to", ".agents/skills", "--lock", lp})
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "auth-skill")
	assert.NotContains(t, out, "lint-skill")
}

func TestInstalledToNoMatch(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64),
		InstalledTo: []string{".agents/skills"},
	})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"installed-to", ".other/path", "--lock", lp})
	require.NoError(t, root.Execute())

	assert.Contains(t, buf.String(), "No skills installed to")
}

func TestInstalledToJSON(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name: "auth-skill", Version: "1.0.0", Source: "local", SHA256: strings.Repeat("a", 64),
		InstalledTo: []string{".agents/skills"},
	})
	require.NoError(t, lf.Write(lp))

	var buf strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&buf)
	root.SetArgs([]string{"installed-to", ".agents/skills", "--lock", lp, "--output", "json"})
	require.NoError(t, root.Execute())

	var result []map[string]any
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &result))
	require.Len(t, result, 1)
	assert.Equal(t, "auth-skill", result[0]["name"])
}

func TestInstalledToRequiresArg(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"installed-to", "--lock", lp})
	assert.Error(t, root.Execute())
}
