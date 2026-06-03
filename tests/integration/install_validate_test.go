//go:build integration

package integration

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/domehahn/sctl/internal/cache"
	"github.com/domehahn/sctl/internal/installer"
	"github.com/domehahn/sctl/internal/lockfile"
	"github.com/domehahn/sctl/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var exampleSkillFiles = map[string]string{
	"SKILL.md": `---
name: example-skill
description: Integration test skill.
---
# Example Skill
`,
	"VERSION": "1.0.0",
	"skill.yaml": `name: example-skill
version: "1.0.0"
description: Integration test skill.
owners:
  - test-team
compatible_with:
  - claude-code
  - gitlab-duo
`,
	"CHANGELOG.md": "# Changelog\n\n## 1.0.0\n\n- Initial release\n",
}

func TestInstallAndValidate(t *testing.T) {
	srv := newTestServer(t)
	urlPath, sha256sum := srv.addSkill(t, "example-skill", "1.0.0", exampleSkillFiles)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:      "example-skill",
		Version:   "1.0.0",
		Source:    "http",
		SourceURL: srv.skillURL(urlPath),
		SHA256:    sha256sum,
		InstalledTo: []string{
			".claude/skills/example-skill",
			"skills/example-skill",
		},
	})

	workDir := t.TempDir()
	cacheDir := t.TempDir()
	c := cache.New(cacheDir)
	ins := installer.New(c)

	result, err := ins.Install(context.Background(), lf, installer.Options{
		WorkDir:     workDir,
		Concurrency: 2,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"example-skill"}, result.Installed)
	assert.Empty(t, result.FromCache)

	claudePath := filepath.Join(workDir, ".claude", "skills", "example-skill")
	skillsPath := filepath.Join(workDir, "skills", "example-skill")
	assert.DirExists(t, claudePath)
	assert.DirExists(t, skillsPath)
	assert.FileExists(t, filepath.Join(claudePath, "SKILL.md"))
	assert.FileExists(t, filepath.Join(skillsPath, "skill.yaml"))

	v := skill.NewValidator()
	vResult, err := v.Validate(context.Background(), claudePath)
	require.NoError(t, err)
	assert.True(t, vResult.Valid, "installed skill should be valid: %+v", vResult.Errors)
}

func TestInstallFromCache(t *testing.T) {
	srv := newTestServer(t)
	urlPath, sha256sum := srv.addSkill(t, "example-skill", "1.0.0", exampleSkillFiles)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "example-skill",
		Version:     "1.0.0",
		Source:      "http",
		SourceURL:   srv.skillURL(urlPath),
		SHA256:      sha256sum,
		InstalledTo: []string{".claude/skills/example-skill"},
	})

	cacheDir := t.TempDir()
	c := cache.New(cacheDir)
	ins := installer.New(c)
	ctx := context.Background()
	opts := installer.Options{WorkDir: t.TempDir(), Concurrency: 1}

	// First install — downloads from server.
	_, err := ins.Install(ctx, lf, opts)
	require.NoError(t, err)
	assert.True(t, c.Has(sha256sum))

	// Second install — must use cache, not re-download.
	srv.Server.Close()
	result, err := ins.Install(ctx, lf, installer.Options{WorkDir: t.TempDir(), Concurrency: 1})
	require.NoError(t, err)
	assert.Equal(t, []string{"example-skill"}, result.FromCache)
}

func TestInstallDryRun(t *testing.T) {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "example-skill",
		Version:     "1.0.0",
		Source:      "http",
		SourceURL:   "http://localhost:9999/nonexistent.zip",
		SHA256:      "deadbeef",
		InstalledTo: []string{".claude/skills/example-skill"},
	})

	workDir := t.TempDir()
	ins := installer.New(cache.New(t.TempDir()))
	result, err := ins.Install(context.Background(), lf, installer.Options{
		DryRun:  true,
		WorkDir: workDir,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"example-skill"}, result.Skipped)
	assert.NoDirExists(t, filepath.Join(workDir, ".claude"))
}

func TestPackageAndValidate(t *testing.T) {
	skillDir := writeFixtureSkill(t)

	p := skill.NewPackager()
	outDir := t.TempDir()
	result, err := p.Package(context.Background(), skillDir, outDir)
	require.NoError(t, err)
	assert.Equal(t, "example-skill", result.Name)
	assert.Equal(t, "1.0.0", result.Version)
	assert.NotEmpty(t, result.SHA256)

	v := skill.NewValidator()
	vResult, err := v.Validate(context.Background(), skillDir)
	require.NoError(t, err)
	assert.True(t, vResult.Valid)
}
