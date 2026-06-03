package installer

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/sctl/internal/lockfile"
	"github.com/domehahn/sctl/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestZIP(t *testing.T, files map[string]string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.zip")
	require.NoError(t, err)
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, content := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return f.Name()
}

func TestAtomicUnzip(t *testing.T) {
	zipPath := buildTestZIP(t, map[string]string{
		"SKILL.md":   "# Hello",
		"skill.yaml": "name: test",
	})

	dest := filepath.Join(t.TempDir(), "installed")
	require.NoError(t, atomicUnzip(zipPath, dest))

	assert.FileExists(t, filepath.Join(dest, "SKILL.md"))
	assert.FileExists(t, filepath.Join(dest, "skill.yaml"))
}

func TestAtomicUnzipReplaces(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "skill")
	require.NoError(t, os.MkdirAll(dest, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dest, "old.txt"), []byte("old"), 0o644))

	zipPath := buildTestZIP(t, map[string]string{"SKILL.md": "new"})
	require.NoError(t, atomicUnzip(zipPath, dest))

	assert.FileExists(t, filepath.Join(dest, "SKILL.md"))
	assert.NoFileExists(t, filepath.Join(dest, "old.txt"))
}

func TestIsWithinDir(t *testing.T) {
	assert.True(t, isWithinDir("/base", "/base/sub/file"))
	assert.False(t, isWithinDir("/base", "/other"))
}

func TestResolvePaths(t *testing.T) {
	tests := []struct {
		name      string
		platforms []skill.Platform
		wantPaths []string
	}{
		{
			name:      "claude-code",
			platforms: []skill.Platform{skill.PlatformClaudeCode},
			wantPaths: []string{".claude/skills/my-skill"},
		},
		{
			name:      "gitlab-duo",
			platforms: []skill.Platform{skill.PlatformGitLabDuo},
			wantPaths: []string{"skills/my-skill", ".agents/skills/my-skill"},
		},
		{
			name:      "dedup: gitlab-duo + codex share .agents/skills",
			platforms: []skill.Platform{skill.PlatformGitLabDuo, skill.PlatformCodex},
			wantPaths: []string{"skills/my-skill", ".agents/skills/my-skill"},
		},
		{
			name:      "github-copilot",
			platforms: []skill.Platform{skill.PlatformGitHubCopilot},
			wantPaths: []string{".github/skills/my-skill"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paths, err := ResolvePaths("my-skill", tc.platforms)
			require.NoError(t, err)
			assert.Equal(t, tc.wantPaths, paths)
		})
	}
}

func buildTestLockFile(name, version, sourceURL, sha256 string) *lockfile.LockFile {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        name,
		Version:     version,
		Source:      "github",
		SourceURL:   sourceURL,
		SHA256:      sha256,
		InstalledTo: []string{"skills/" + name},
	})
	return lf
}

func TestInstallDryRun(t *testing.T) {
	dir := t.TempDir()
	ins := New(nil)

	lf := buildTestLockFile("my-skill", "1.0.0", "http://example.com/my-skill.zip", "abc123")
	result, err := ins.Install(context.Background(), lf, Options{
		DryRun:  true,
		WorkDir: dir,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"my-skill"}, result.Skipped)
	assert.Empty(t, result.Installed)

	entries, _ := os.ReadDir(dir)
	assert.Empty(t, entries, "dry run must not write any files")
}
