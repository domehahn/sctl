package cli_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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

// ── freeze ────────────────────────────────────────────────────────────────────

func TestFreezeCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "freeze" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock"))
		}
	}
	assert.True(t, found, "freeze command should be registered on root")
}

func TestFreezeUpdateManifestVersions(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	// Write lockfile with known version.
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-a",
		Version: "2.3.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Write manifest with a different version.
	mf := "version: 1\nskills:\n  - name: skill-a\n    version: 1.0.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"freeze", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "2.3.0")

	// Verify manifest was updated.
	data, err := os.ReadFile(filepath.Join(dir, "agent-skills.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "2.3.0")
}

func TestFreezeAlreadyFrozen(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-b",
		Version: "1.5.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("b", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	mf := "version: 1\nskills:\n  - name: skill-b\n    version: 1.5.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"freeze", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "already match")
}

func TestFreezeDryRun(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-c",
		Version: "3.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("c", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	mf := "version: 1\nskills:\n  - name: skill-c\n    version: 1.0.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"freeze", "--lock", lp, "--dry-run"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "Dry run")

	// Manifest should not be changed.
	data, _ := os.ReadFile(filepath.Join(dir, "agent-skills.yaml"))
	assert.Contains(t, string(data), "1.0.0")
	assert.NotContains(t, string(data), "3.0.0")
}

func TestFreezeJSONOutput(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "skill-d",
		Version: "4.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("d", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	mf := "version: 1\nskills:\n  - name: skill-d\n    version: 0.0.1\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"freeze", "--lock", lp, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var results []map[string]string
	require.NoError(t, json.Unmarshal(buf.Bytes(), &results))
	assert.Len(t, results, 1)
	assert.Equal(t, "skill-d", results[0]["skill"])
	assert.Equal(t, "4.0.0", results[0]["version"])
}

func TestFreezeSkipsSkillsNotInManifest(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "extra-skill",
		Version: "9.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("e", 64),
	})
	lp := filepath.Join(dir, "agent-skills.lock.yaml")
	require.NoError(t, lf.Write(lp))

	// Manifest has no skills.
	mf := "version: 1\nskills: []\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "agent-skills.yaml"), []byte(mf), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"freeze", "--lock", lp})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "already match")
}

// ── list-protected ────────────────────────────────────────────────────────────

func TestListProtectedCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "list-protected" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("lock-dir"))
		}
	}
	assert.True(t, found, "list-protected command should be registered on root")
}

func TestListProtectedEmpty(t *testing.T) {
	dir := t.TempDir()

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"list-protected", "--lock-dir", dir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "No protected")
}

func TestListProtectedShowsNames(t *testing.T) {
	dir := t.TempDir()

	// Add two protected skills.
	for _, name := range []string{"skill-a", "skill-b"} {
		root := cli.NewRootCmd()
		root.SetArgs([]string{"protect", name, "--lock-dir", dir})
		root.SetOut(&bytes.Buffer{})
		require.NoError(t, root.Execute())
	}

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"list-protected", "--lock-dir", dir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	out := buf.String()
	assert.Contains(t, out, "skill-a")
	assert.Contains(t, out, "skill-b")
	assert.Contains(t, out, "2 protected")
}

func TestListProtectedJSONOutput(t *testing.T) {
	dir := t.TempDir()

	root := cli.NewRootCmd()
	root.SetArgs([]string{"protect", "skill-x", "--lock-dir", dir})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	var buf bytes.Buffer
	root = cli.NewRootCmd()
	root.SetArgs([]string{"list-protected", "--lock-dir", dir, "--output", "json"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	var names []string
	require.NoError(t, json.Unmarshal(buf.Bytes(), &names))
	assert.Contains(t, names, "skill-x")
}

// ── import ────────────────────────────────────────────────────────────────────

func TestUnpackCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "unpack <tarball>" {
			found = true
			assert.NotNil(t, sub.Flags().Lookup("dir"))
			assert.NotNil(t, sub.Flags().Lookup("force"))
		}
	}
	assert.True(t, found, "unpack command should be registered on root")
}

func TestUnpackRequiresTarball(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unpack"})
	root.SetOut(&bytes.Buffer{})
	assert.Error(t, root.Execute())
}

func makeTestTarball(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "skills.tar.gz")

	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)

	content := []byte("---\nname: test-skill\nversion: 1.0.0\n---\n\n# Test Skill\n")
	hdr := &tar.Header{
		Name: "test-skill/SKILL.md",
		Mode: 0o644,
		Size: int64(len(content)),
	}
	require.NoError(t, tw.WriteHeader(hdr))
	_, err = tw.Write(content)
	require.NoError(t, err)

	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())
	return path
}

func TestUnpackExtractsFiles(t *testing.T) {
	tarball := makeTestTarball(t)
	destDir := t.TempDir()

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unpack", tarball, "--dir", destDir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())

	assert.FileExists(t, filepath.Join(destDir, "test-skill", "SKILL.md"))
	assert.Contains(t, buf.String(), "Imported")
}

func TestUnpackSkipsExistingFiles(t *testing.T) {
	tarball := makeTestTarball(t)
	destDir := t.TempDir()

	skillDir := filepath.Join(destDir, "test-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("existing"), 0o644))

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unpack", tarball, "--dir", destDir})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "skip")

	data, _ := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	assert.Equal(t, "existing", string(data))
}

func TestUnpackForceOverwrites(t *testing.T) {
	tarball := makeTestTarball(t)
	destDir := t.TempDir()

	skillDir := filepath.Join(destDir, "test-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("old"), 0o644))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"unpack", tarball, "--dir", destDir, "--force"})
	root.SetOut(&bytes.Buffer{})
	require.NoError(t, root.Execute())

	data, _ := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	assert.NotEqual(t, "old", string(data))
}

func TestUnpackDryRun(t *testing.T) {
	tarball := makeTestTarball(t)
	destDir := t.TempDir()

	var buf bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs([]string{"unpack", tarball, "--dir", destDir, "--dry-run"})
	root.SetOut(&buf)
	require.NoError(t, root.Execute())
	assert.Contains(t, buf.String(), "would create")
	assert.Contains(t, buf.String(), "Dry run")

	_, err := os.Stat(filepath.Join(destDir, "test-skill", "SKILL.md"))
	assert.True(t, os.IsNotExist(err))
}
