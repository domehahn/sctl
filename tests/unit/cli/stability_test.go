package cli_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// discardCmd returns a cobra.Command that silently discards all output.
func discardCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	return cmd
}

// ── skpm run / readSkillScripts ───────────────────────────────────────────────

func TestReadSkillScriptsNoFile(t *testing.T) {
	_, err := cli.ReadSkillScripts(t.TempDir())
	assert.Error(t, err)
}

func TestReadSkillScriptsNoScriptsSection(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "skill.yaml"), []byte("name: my-skill\nversion: 0.1.0\n"), 0o644))
	scripts, err := cli.ReadSkillScripts(dir)
	require.NoError(t, err)
	assert.Empty(t, scripts)
}

func TestReadSkillScriptsParsesScripts(t *testing.T) {
	dir := t.TempDir()
	content := "name: s\nversion: 0.1.0\nscripts:\n  test: skpm validate .\n  lint: skpm lint .\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "skill.yaml"), []byte(content), 0o644))
	scripts, err := cli.ReadSkillScripts(dir)
	require.NoError(t, err)
	assert.Equal(t, "skpm validate .", scripts["test"])
	assert.Equal(t, "skpm lint .", scripts["lint"])
}

func TestScriptNamesReturnsKeys(t *testing.T) {
	scripts := map[string]string{"test": "go test", "build": "go build"}
	names := cli.ScriptNames(scripts)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "test")
	assert.Contains(t, names, "build")
}

// ── skpm hooks ────────────────────────────────────────────────────────────────

func TestRunProjectHookNoManifest(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	// Should silently return nil — no manifest present.
	err := cli.RunProjectHook(discardCmd(), "pre_install")
	assert.NoError(t, err)
}

func TestRunProjectHookMissingHookKey(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	manifest := "version: 1\nskills: []\nhooks:\n  post_install: echo done\n"
	require.NoError(t, os.WriteFile("agent-skills.yaml", []byte(manifest), 0o644))

	// A hook that is not defined should be a no-op.
	err := cli.RunProjectHook(discardCmd(), "pre_install")
	assert.NoError(t, err)
}

// ── skpm template ─────────────────────────────────────────────────────────────

func TestIsTemplateFile(t *testing.T) {
	assert.True(t, cli.IsTemplateFile("SKILL.md"))
	assert.True(t, cli.IsTemplateFile("skill.yaml"))
	assert.True(t, cli.IsTemplateFile("VERSION"))
	assert.True(t, cli.IsTemplateFile("CHANGELOG.md"))
	assert.True(t, cli.IsTemplateFile("README.md"))
	assert.True(t, cli.IsTemplateFile("notes.txt"))
	assert.False(t, cli.IsTemplateFile("main.go"))
	assert.False(t, cli.IsTemplateFile("binary.zip"))
	assert.False(t, cli.IsTemplateFile("image.png"))
}

func TestLoadTemplateRegistryNotFound(t *testing.T) {
	dir := t.TempDir()
	// Point config dir to temp so templates.yaml doesn't exist.
	t.Setenv("XDG_CONFIG_HOME", dir)
	reg, err := cli.LoadTemplateRegistry()
	require.NoError(t, err)
	assert.Empty(t, reg.Templates)
}

func TestSaveAndLoadTemplateRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	reg, err := cli.LoadTemplateRegistry()
	require.NoError(t, err)

	reg.Templates["review"] = cli.SkillTemplate{
		Name:        "review",
		Description: "Code review skill",
		Files:       map[string]string{"SKILL.md": "# {{.Name}}\n"},
	}

	require.NoError(t, cli.SaveTemplateRegistry(reg))

	loaded, err := cli.LoadTemplateRegistry()
	require.NoError(t, err)
	assert.Contains(t, loaded.Templates, "review")
	assert.Equal(t, "Code review skill", loaded.Templates["review"].Description)
	assert.Contains(t, loaded.Templates["review"].Files["SKILL.md"], "{{.Name}}")
}

func TestSortedTemplateNamesAlphabetical(t *testing.T) {
	reg := &cli.TemplateRegistry{
		Templates: map[string]cli.SkillTemplate{
			"zebra": {Name: "zebra"},
			"alpha": {Name: "alpha"},
			"beta":  {Name: "beta"},
		},
	}
	names := cli.SortedTemplateNames(reg)
	assert.Equal(t, []string{"alpha", "beta", "zebra"}, names)
}

// ── skpm migrate ──────────────────────────────────────────────────────────────

func TestRenameYAMLKey(t *testing.T) {
	input := "platforms:\n  - claude-code\nversion: 0.1.0\n"
	var doc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(input), &doc))

	cli.RenameYAMLKey(&doc, "platforms", "compatible_with")

	out, err := yaml.Marshal(&doc)
	require.NoError(t, err)
	assert.Contains(t, string(out), "compatible_with")
	assert.NotContains(t, string(out), "platforms")
}

func TestRenameYAMLKeyMissingKey(t *testing.T) {
	input := "name: skill\nversion: 0.1.0\n"
	var doc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(input), &doc))
	// Should not panic or error on missing key.
	assert.NotPanics(t, func() { cli.RenameYAMLKey(&doc, "nonexistent", "new_key") })
	out, _ := yaml.Marshal(&doc)
	assert.NotContains(t, string(out), "new_key")
}

func TestMigrateSkillMDAddsMissingSpecVersion(t *testing.T) {
	dir := t.TempDir()
	content := "---\nname: my-skill\ncompatibility:\n  platforms:\n    - claude-code\n---\n\n# Body\n"
	mdPath := filepath.Join(dir, "SKILL.md")
	require.NoError(t, os.WriteFile(mdPath, []byte(content), 0o644))

	changed, err := cli.MigrateSkillMD(discardCmd(), mdPath, false)
	require.NoError(t, err)
	assert.True(t, changed)

	data, _ := os.ReadFile(mdPath)
	assert.Contains(t, string(data), "spec_version: 1")
}

func TestMigrateSkillMDIdempotent(t *testing.T) {
	dir := t.TempDir()
	content := "---\nname: my-skill\ncompatibility:\n  spec_version: 1\n---\n\n# Body\n"
	mdPath := filepath.Join(dir, "SKILL.md")
	require.NoError(t, os.WriteFile(mdPath, []byte(content), 0o644))

	changed, err := cli.MigrateSkillMD(discardCmd(), mdPath, false)
	require.NoError(t, err)
	assert.False(t, changed)
}

func TestMigrateSkillYAMLFileAddsMissingNamespace(t *testing.T) {
	dir := t.TempDir()
	content := "name: my-skill\nversion: 0.1.0\ncompatible_with:\n  - claude-code\n"
	syPath := filepath.Join(dir, "skill.yaml")
	require.NoError(t, os.WriteFile(syPath, []byte(content), 0o644))

	changed, err := cli.MigrateSkillYAMLFile(discardCmd(), syPath, false)
	require.NoError(t, err)
	assert.True(t, changed)

	data, _ := os.ReadFile(syPath)
	assert.Contains(t, string(data), "namespace")
}

func TestMigrateSkillYAMLFileRenamesPlatforms(t *testing.T) {
	dir := t.TempDir()
	content := "name: my-skill\nversion: 0.1.0\nplatforms:\n  - claude-code\n"
	syPath := filepath.Join(dir, "skill.yaml")
	require.NoError(t, os.WriteFile(syPath, []byte(content), 0o644))

	changed, err := cli.MigrateSkillYAMLFile(discardCmd(), syPath, false)
	require.NoError(t, err)
	assert.True(t, changed)

	data, _ := os.ReadFile(syPath)
	assert.Contains(t, string(data), "compatible_with")
	assert.NotContains(t, string(data), "platforms:")
}

func TestMigrateSkillYAMLFileDryRun(t *testing.T) {
	dir := t.TempDir()
	content := "name: my-skill\nversion: 0.1.0\nplatforms:\n  - all\n"
	syPath := filepath.Join(dir, "skill.yaml")
	require.NoError(t, os.WriteFile(syPath, []byte(content), 0o644))

	changed, err := cli.MigrateSkillYAMLFile(discardCmd(), syPath, true)
	require.NoError(t, err)
	assert.True(t, changed)

	// Dry run — file must not be modified.
	data, _ := os.ReadFile(syPath)
	assert.Contains(t, string(data), "platforms:")
}

// ── skpm graph / readGraphNode ────────────────────────────────────────────────

func TestReadGraphNodeNoFile(t *testing.T) {
	_, err := cli.ReadGraphNode(t.TempDir())
	assert.Error(t, err)
}

func TestReadGraphNodeParsesRequires(t *testing.T) {
	dir := t.TempDir()
	content := "name: my-skill\nversion: 1.0.0\nrequires:\n  - dep-a\n  - dep-b\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "skill.yaml"), []byte(content), 0o644))

	node, err := cli.ReadGraphNode(dir)
	require.NoError(t, err)
	assert.Equal(t, "my-skill", node.Name)
	assert.Equal(t, "1.0.0", node.Version)
	assert.Equal(t, []string{"dep-a", "dep-b"}, node.Requires)
}

func TestReadGraphNodeNoRequires(t *testing.T) {
	dir := t.TempDir()
	content := "name: leaf-skill\nversion: 0.2.0\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "skill.yaml"), []byte(content), 0o644))

	node, err := cli.ReadGraphNode(dir)
	require.NoError(t, err)
	assert.Empty(t, node.Requires)
}

// ── skpm workspace ────────────────────────────────────────────────────────────

func TestReadWorkspaceNoFile(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	_, err := cli.ReadWorkspace()
	assert.Error(t, err)
}

func TestReadWorkspaceParses(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	content := "version: 1\nskills:\n  - ./skill-a\n  - ./skill-b\n"
	require.NoError(t, os.WriteFile("skpm-workspace.yaml", []byte(content), 0o644))

	wm, err := cli.ReadWorkspace()
	require.NoError(t, err)
	assert.Equal(t, []string{"./skill-a", "./skill-b"}, wm.Skills)
}

func TestFindWorkspaceFileNotFound(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	assert.Equal(t, "", cli.FindWorkspaceFile())
}

func TestFindWorkspaceFileInCwd(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	require.NoError(t, os.WriteFile("skpm-workspace.yaml", []byte("version: 1\nskills: [./s]\n"), 0o644))
	found := cli.FindWorkspaceFile()
	assert.NotEmpty(t, found)
	assert.True(t, strings.HasSuffix(found, "skpm-workspace.yaml"))
}

func TestFindWorkspaceFileInParent(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "sub", "pkg")
	require.NoError(t, os.MkdirAll(child, 0o755))

	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(child))
	defer os.Chdir(orig) //nolint:errcheck

	require.NoError(t, os.WriteFile(filepath.Join(parent, "skpm-workspace.yaml"), []byte("version: 1\nskills: [./s]\n"), 0o644))

	found := cli.FindWorkspaceFile()
	assert.NotEmpty(t, found)
	assert.Contains(t, found, parent)
}

func TestWorkspaceDirsAll(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	require.NoError(t, os.WriteFile("skpm-workspace.yaml", []byte("version: 1\nskills:\n  - ./a\n  - ./b\n  - ./c\n"), 0o644))
	dirs, err := cli.WorkspaceDirs(nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"./a", "./b", "./c"}, dirs)
}

func TestWorkspaceDirsFilter(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	require.NoError(t, os.WriteFile("skpm-workspace.yaml", []byte("version: 1\nskills:\n  - ./a\n  - ./b\n  - ./c\n"), 0o644))
	dirs, err := cli.WorkspaceDirs([]string{"./a", "./c"})
	require.NoError(t, err)
	assert.Equal(t, []string{"./a", "./c"}, dirs)
}

func TestFilterChangedDirsNoGit(t *testing.T) {
	// In a directory that is not a git repo, filterChangedDirs should return all dirs.
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(orig) //nolint:errcheck

	input := []string{"./a", "./b"}
	result, err := cli.FilterChangedDirs(input)
	require.NoError(t, err)
	// Not a git repo → falls back to returning all dirs unchanged.
	assert.Equal(t, input, result)
}
