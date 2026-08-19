package cli_test

import (
	"bytes"
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

// ── lock-add ──────────────────────────────────────────────────────────────────

func TestLockAddCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "lock-add <name> <version> <source> <sha256>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("installed-to"))
		}
	}
	assert.True(t, found, "lock-add command should be registered")
}

func TestLockAddRequiresFourArgs(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-add", "name", "version", "source"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestLockAddInsertsEntry(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	sha := strings.Repeat("a", 64)
	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-add", "new-skill", "1.0.0", "https://r.example.com", sha, "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "new-skill")

	lf, err := lockfile.Read(lp)
	require.NoError(t, err)
	sl, found := lf.Find("new-skill")
	require.True(t, found)
	assert.Equal(t, "1.0.0", sl.Version)
	assert.Equal(t, sha, sl.SHA256)
}

func TestLockAddErrorsOnDuplicate(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "existing",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-add", "existing", "2.0.0", "https://r.example.com", strings.Repeat("b", 64), "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestLockAddRejectsShortSHA(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-add", "s", "1.0.0", "https://r.example.com", "tooshort", "--lock", lp})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func TestLockAddWithInstalledTo(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	root := cli.NewRootCmd()
	root.SetArgs([]string{
		"lock-add", "skill-x", "1.0.0", "https://r.example.com", strings.Repeat("c", 64),
		"--lock", lp,
		"--installed-to", ".agents/skills/skill-x",
	})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	lf, _ := lockfile.Read(lp)
	sl, found := lf.Find("skill-x")
	require.True(t, found)
	assert.Contains(t, sl.InstalledTo, ".agents/skills/skill-x")
}

func TestLockAddDryRun(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"lock-add", "dry-skill", "1.0.0", "https://r.example.com", strings.Repeat("d", 64), "--lock", lp, "--dry-run"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Would add")

	lf, _ := lockfile.Read(lp)
	_, found := lf.Find("dry-skill")
	assert.False(t, found)
}

// ── summarize ─────────────────────────────────────────────────────────────────

func TestSummarizeCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "summarize" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
		}
	}
	assert.True(t, found, "summarize command should be registered")
}

func TestSummarizeEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"summarize", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "empty")
}

func TestSummarizeShowsStats(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.0.0", Source: "https://r1.example.com", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "skill-b", Version: "2.0.0", Source: "https://r2.example.com", SHA256: strings.Repeat("b", 64), Signature: "sig"})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"summarize", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "2") // total skills
	assert.Contains(t, out, "1.0.0")
	assert.Contains(t, out, "2.0.0")
}

func TestSummarizeJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "s", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"summarize", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	assert.Contains(t, result, "total_skills")
	assert.Contains(t, result, "unique_sources")
}

// ── namespace ─────────────────────────────────────────────────────────────────

func TestNamespaceCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "namespace" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
			assert.NotNil(t, sub.Flags().Lookup("namespace"))
		}
	}
	assert.True(t, found, "namespace command should be registered")
}

func TestNamespaceEmptyLockfile(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"namespace", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No skills")
}

func TestNamespaceListsDefaultNamespace(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "skill-a", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lf.Upsert(lockfile.SkillLock{Name: "skill-b", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("b", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"namespace", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "default")
	assert.Contains(t, buf.String(), "1 namespace")
}

func TestNamespaceFilterByNamespace(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()

	sl := lockfile.SkillLock{Name: "ns-skill", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)}
	sl.Namespace = "my-ns"
	lf.Upsert(sl)
	lf.Upsert(lockfile.SkillLock{Name: "other", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("b", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"namespace", "--lock", lp, "--namespace", "my-ns"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "ns-skill")
	assert.NotContains(t, buf.String(), "other")
}

func TestNamespaceFilterMissing(t *testing.T) {
	dir := t.TempDir()
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lockfile.New().Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"namespace", "--lock", lp, "--namespace", "ghost-ns"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No skills")
}

func TestNamespaceJSONOutput(t *testing.T) {
	dir := t.TempDir()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{Name: "s", Version: "1.0.0", Source: "https://r.example.com", SHA256: strings.Repeat("a", 64)})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"namespace", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var result []map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	assert.NotEmpty(t, result)
	assert.Contains(t, result[0], "namespace")

	// Write the protect file needed by a separate test.
	_ = os.MkdirAll(dir, 0o755)
}
