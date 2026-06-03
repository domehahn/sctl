package skill

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageValidSkill(t *testing.T) {
	dir := writeSkillFixture(t, validSkillFiles)
	p := NewPackager()
	outDir := t.TempDir()

	result, err := p.Package(context.Background(), dir, outDir)
	require.NoError(t, err)
	assert.Equal(t, "my-skill", result.Name)
	assert.Equal(t, "1.2.3", result.Version)
	assert.NotEmpty(t, result.SHA256)

	expectedName := "my-skill-1.2.3.zip"
	assert.Equal(t, filepath.Join(outDir, expectedName), result.OutputPath)

	_, err = os.Stat(result.OutputPath)
	require.NoError(t, err)

	zr, err := zip.OpenReader(result.OutputPath)
	require.NoError(t, err)
	defer zr.Close()

	fileNames := make(map[string]bool)
	var manifest SkillManifest
	for _, f := range zr.File {
		fileNames[f.Name] = true
		if f.Name == "manifest.json" {
			rc, _ := f.Open()
			json.NewDecoder(rc).Decode(&manifest)
			rc.Close()
		}
	}

	assert.True(t, fileNames["SKILL.md"])
	assert.True(t, fileNames["skill.yaml"])
	assert.True(t, fileNames["manifest.json"])
	assert.Equal(t, "my-skill", manifest.Name)
	assert.Equal(t, "1.2.3", manifest.Version)
	assert.NotEmpty(t, manifest.CreatedAt)
}

func TestPackageInvalidSkill(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "SKILL.md")
	dir := writeSkillFixture(t, files)
	p := NewPackager()

	_, err := p.Package(context.Background(), dir, t.TempDir())
	require.Error(t, err)
	var vfe *ValidationFailedError
	assert.ErrorAs(t, err, &vfe)
}

func TestPackageNoTmpFileLeftOnError(t *testing.T) {
	files := copyMap(validSkillFiles)
	delete(files, "SKILL.md")
	dir := writeSkillFixture(t, files)
	outDir := t.TempDir()
	p := NewPackager()

	_, _ = p.Package(context.Background(), dir, outDir)

	entries, _ := os.ReadDir(outDir)
	assert.Empty(t, entries)
}
