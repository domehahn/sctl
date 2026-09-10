//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNegativeE2E_TamperedArchiveRejected(t *testing.T) {
	srv := newTestServer(t)
	urlPath, _ := srv.addSkill(t, "example-skill", "1.0.0", exampleSkillFiles)

	// Lockfile specifies wrong SHA256 digest
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "example-skill",
		Version:     "1.0.0",
		Source:      "http",
		DownloadURL: srv.skillURL(urlPath),
		SHA256:      "0000000000000000000000000000000000000000000000000000000000000000",
		InstalledTo: []string{".claude/skills/example-skill"},
	})

	workDir := t.TempDir()
	cacheDir := t.TempDir()
	c := cache.New(cacheDir)
	ins := installer.New(c)

	_, err := ins.Install(context.Background(), lf, installer.Options{
		WorkDir: workDir,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SHA256 mismatch")
}

func TestInstallHandoffMetadataSchema(t *testing.T) {
	srv := newTestServer(t)
	urlPath, sha256sum := srv.addSkill(t, "example-skill", "1.0.0", exampleSkillFiles)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "example-skill",
		Namespace:   "default",
		Version:     "1.0.0",
		Source:      "http",
		DownloadURL: srv.skillURL(urlPath),
		SHA256:      sha256sum,
		InstalledTo: []string{"skills/example-skill"},
	})

	workDir := t.TempDir()
	cacheDir := t.TempDir()
	c := cache.New(cacheDir)
	ins := installer.New(c)

	_, err := ins.Install(context.Background(), lf, installer.Options{
		WorkDir:             workDir,
		LockfileDigest:      "sha256:lock123",
		AdmissionDecisionID: "adm-dec-789",
	})
	require.NoError(t, err)

	manifestPath := filepath.Join(workDir, "skills", "example-skill", ".skpm-installed.json")
	assert.FileExists(t, manifestPath)

	data, err := os.ReadFile(manifestPath)
	require.NoError(t, err)

	var id installer.MaterializedIdentity
	require.NoError(t, json.Unmarshal(data, &id))

	assert.Equal(t, "1.0", id.SchemaVersion)
	assert.Equal(t, "example-skill@1.0.0", id.Package)
	assert.Equal(t, "sha256:"+sha256sum, id.PackageDigest)
	assert.NotEmpty(t, id.CompiledDigest)
	assert.Equal(t, "http", id.Registry)
	assert.Equal(t, "sha256:lock123", id.LockfileDigest)
	assert.Equal(t, "adm-dec-789", id.AdmissionDecisionID)
	assert.True(t, id.Verified)
}
