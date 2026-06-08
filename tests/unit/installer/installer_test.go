package installer_test

import (
	"archive/zip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/internal/cache"
	"github.com/domehahn/skpm/internal/installer"
	"github.com/domehahn/skpm/internal/lockfile"
	"github.com/domehahn/skpm/internal/skill"
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

func TestAtomicUnzip(t *testing.T) {
	zipPath := buildTestZIP(t, map[string]string{
		"SKILL.md":   "# Hello",
		"skill.yaml": "name: test",
	})
	dest := filepath.Join(t.TempDir(), "installed")
	require.NoError(t, installer.AtomicUnzip(zipPath, dest))
	assert.FileExists(t, filepath.Join(dest, "SKILL.md"))
	assert.FileExists(t, filepath.Join(dest, "skill.yaml"))
}

func TestAtomicUnzipReplaces(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "skill")
	require.NoError(t, os.MkdirAll(dest, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dest, "old.txt"), []byte("old"), 0o644))

	zipPath := buildTestZIP(t, map[string]string{"SKILL.md": "new"})
	require.NoError(t, installer.AtomicUnzip(zipPath, dest))

	assert.FileExists(t, filepath.Join(dest, "SKILL.md"))
	assert.NoFileExists(t, filepath.Join(dest, "old.txt"))
}

func TestAtomicUnzipInvalidZIP(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "bad.zip")
	require.NoError(t, os.WriteFile(bad, []byte("not a zip"), 0o644))
	err := installer.AtomicUnzip(bad, filepath.Join(t.TempDir(), "dest"))
	assert.Error(t, err)
}

func TestIsWithinDir(t *testing.T) {
	assert.True(t, installer.IsWithinDir("/base", "/base/sub/file"))
	assert.False(t, installer.IsWithinDir("/base", "/other"))
	assert.False(t, installer.IsWithinDir("/base", "/base/../escape"))
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
		{
			name:      "all expands to every platform",
			platforms: []skill.Platform{skill.PlatformAll},
			wantPaths: []string{".claude/skills/my-skill", "skills/my-skill", ".agents/skills/my-skill", ".github/skills/my-skill"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paths, err := installer.ResolvePaths("my-skill", tc.platforms)
			require.NoError(t, err)
			assert.Equal(t, tc.wantPaths, paths)
		})
	}
}

func TestResolvePathsUnknownPlatform(t *testing.T) {
	_, err := installer.ResolvePaths("skill", []skill.Platform{"unknown"})
	require.Error(t, err)
}

func TestInstallDryRun(t *testing.T) {
	dir := t.TempDir()
	ins := installer.New(nil)
	lf := buildTestLockFile("my-skill", "1.0.0", "http://example.com/my-skill.zip", "abc123")

	result, err := ins.Install(context.Background(), lf, installer.Options{DryRun: true, WorkDir: dir})
	require.NoError(t, err)
	assert.Equal(t, []string{"my-skill"}, result.Skipped)
	assert.Empty(t, result.Installed)
	entries, _ := os.ReadDir(dir)
	assert.Empty(t, entries)
}

func TestInstallFromHTTP(t *testing.T) {
	zipPath := buildTestZIP(t, map[string]string{
		"SKILL.md":   "# Test Skill",
		"skill.yaml": "name: test-skill\nversion: \"1.0.0\"\n",
	})
	zipData, err := os.ReadFile(zipPath)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(zipData)
	}))
	defer srv.Close()

	workDir := t.TempDir()
	ins := installer.New(cache.New(t.TempDir()))
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "test-skill",
		Version:     "1.0.0",
		SourceURL:   srv.URL + "/test-skill-1.0.0.zip",
		InstalledTo: []string{"skills/test-skill"},
	})

	result, err := ins.Install(context.Background(), lf, installer.Options{WorkDir: workDir, Concurrency: 1})
	require.NoError(t, err)
	assert.Equal(t, []string{"test-skill"}, result.Installed)
	assert.FileExists(t, filepath.Join(workDir, "skills", "test-skill", "SKILL.md"))
}

func TestInstallSHA256Mismatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fake content"))
	}))
	defer srv.Close()

	ins := installer.New(cache.New(t.TempDir()))
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "skill",
		SourceURL:   srv.URL + "/skill.zip",
		SHA256:      "wrong-sha",
		InstalledTo: []string{"skills/skill"},
	})

	_, err := ins.Install(context.Background(), lf, installer.Options{WorkDir: t.TempDir(), Concurrency: 1})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SHA256 mismatch")
}
