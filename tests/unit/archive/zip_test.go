package archive_test

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/archive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildZIP writes a well-formed archive with the given entry names/content
// and returns its path.
func buildZIP(t *testing.T, files map[string]string) string {
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

// buildMaliciousZIP writes an archive with one or more entries whose names
// are the attack payloads themselves (path traversal, absolute paths, ...).
// archive/zip itself performs no validation of entry names — this is
// exactly the property a zip-slip attack exploits, and exactly what
// tests in this file confirm this package's own extraction rejects.
func buildMaliciousZIP(t *testing.T, entryNames ...string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "evil-*.zip")
	require.NoError(t, err)
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, name := range entryNames {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
		require.NoError(t, err)
		_, err = w.Write([]byte("evil payload"))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return f.Name()
}

// ── SafeJoin ─────────────────────────────────────────────────────────────

func TestSafeJoinAllowsNormalPaths(t *testing.T) {
	dest := "/base"
	out, err := archive.SafeJoin(dest, "SKILL.md")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dest, "SKILL.md"), out)

	out, err = archive.SafeJoin(dest, "sub/dir/file.txt")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dest, "sub", "dir", "file.txt"), out)
}

func TestSafeJoinRejectsTraversal(t *testing.T) {
	dest := "/base"
	cases := []string{
		"../evil.txt",
		"../../etc/passwd",
		"../../../../home/user/.ssh/authorized_keys", // the exact payload shape from the review
		"a/../../b",
		"a/b/../../../c",
		"..",
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := archive.SafeJoin(dest, name)
			require.Error(t, err, "expected %q to be rejected as zip slip", name)
			assert.ErrorIs(t, err, archive.ErrZipSlip)
		})
	}
}

func TestSafeJoinRejectsAbsolutePath(t *testing.T) {
	_, err := archive.SafeJoin("/base", "/etc/passwd")
	require.Error(t, err)
}

func TestSafeJoinRejectsEmptyAndNulEntries(t *testing.T) {
	_, err := archive.SafeJoin("/base", "")
	assert.Error(t, err)
	_, err = archive.SafeJoin("/base", "evil\x00.txt")
	assert.Error(t, err)
}

// ── IsWithinDir ──────────────────────────────────────────────────────────

func TestIsWithinDir(t *testing.T) {
	assert.True(t, archive.IsWithinDir("/base", "/base/sub/file"))
	assert.True(t, archive.IsWithinDir("/base", "/base/file"))
	assert.False(t, archive.IsWithinDir("/base", "/base/../other"))
	assert.False(t, archive.IsWithinDir("/base", "/other/file"))
}

// ── Extract: the actual regression test for the reported vulnerability ────

func TestExtractRejectsZipSlip(t *testing.T) {
	parentDir := t.TempDir()
	destDir := filepath.Join(parentDir, "installed")
	evilZip := buildMaliciousZIP(t, "../../../../home/user/.ssh/authorized_keys")

	err := archive.Extract(evilZip, destDir, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, archive.ErrZipSlip)

	// The critical assertion: nothing was written outside destDir. Walk the
	// whole temp root and confirm no file named "authorized_keys" exists
	// anywhere outside destDir itself.
	err = filepath.WalkDir(parentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Base(path) == "authorized_keys" {
			t.Fatalf("zip-slip entry escaped extraction: found %s", path)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestExtractRejectsZipSlipMixedWithLegitimateEntries(t *testing.T) {
	// A real attack hides the traversal entry among legitimate files, hoping
	// extraction of the good entries masks the one bad one. Extract must
	// still fail the whole operation.
	destDir := filepath.Join(t.TempDir(), "installed")
	f, err := os.CreateTemp(t.TempDir(), "mixed-*.zip")
	require.NoError(t, err)
	zw := zip.NewWriter(f)
	for _, e := range []struct{ name, content string }{
		{"SKILL.md", "# hello"},
		{"skill.yaml", "name: test"},
		{"../../outside.txt", "escaped"},
	} {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: e.name, Method: zip.Deflate})
		require.NoError(t, err)
		_, err = w.Write([]byte(e.content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	f.Close()

	err = archive.Extract(f.Name(), destDir, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, archive.ErrZipSlip)
}

func TestExtractRejectsSymlinkEntries(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	f, err := os.CreateTemp(t.TempDir(), "symlink-*.zip")
	require.NoError(t, err)
	zw := zip.NewWriter(f)
	hdr := &zip.FileHeader{Name: "innocuous-looking-link", Method: zip.Deflate}
	hdr.SetMode(fs.ModeSymlink | 0o777)
	w, err := zw.CreateHeader(hdr)
	require.NoError(t, err)
	_, err = w.Write([]byte("/etc/passwd")) // symlink target
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	f.Close()

	err = archive.Extract(f.Name(), destDir, "")
	require.Error(t, err, "symlink entries must be rejected outright")
}

func TestExtractGoodArchive(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	zipPath := buildZIP(t, map[string]string{
		"SKILL.md":        "# Hello",
		"skill.yaml":      "name: test",
		"nested/dir/f.md": "nested content",
	})
	require.NoError(t, archive.Extract(zipPath, destDir, ""))
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
	assert.FileExists(t, filepath.Join(destDir, "skill.yaml"))
	assert.FileExists(t, filepath.Join(destDir, "nested", "dir", "f.md"))
}

func TestExtractWithPrefix(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	zipPath := buildZIP(t, map[string]string{
		"repo-abc123/SKILL.md":       "# Hello",
		"repo-abc123/skill.yaml":     "name: test",
		"repo-abc123/other/file.txt": "x",
	})
	require.NoError(t, archive.Extract(zipPath, destDir, "repo-abc123/"))
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
	assert.FileExists(t, filepath.Join(destDir, "other", "file.txt"))
}

func TestExtractPrefixDoesNotBypassContainment(t *testing.T) {
	// A malicious entry can't use the prefix-stripping path to smuggle a
	// traversal either — SafeJoin runs on the name *after* the prefix is
	// stripped.
	destDir := filepath.Join(t.TempDir(), "installed")
	evilZip := buildMaliciousZIP(t, "repo-abc123/../../../etc/passwd")
	err := archive.Extract(evilZip, destDir, "repo-abc123/")
	require.Error(t, err)
}

// ── AtomicExtract ────────────────────────────────────────────────────────

func TestAtomicExtractRejectsZipSlipAndLeavesDestUntouched(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "existing-skill")
	require.NoError(t, os.MkdirAll(destDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "keep-me.txt"), []byte("original"), 0o644))

	evilZip := buildMaliciousZIP(t, "../../../etc/passwd")
	err := archive.AtomicExtract(evilZip, destDir, "")
	require.Error(t, err)

	// The pre-existing directory must be untouched — this is the "no
	// partial overwrite on failure" guarantee AtomicExtract provides.
	assert.FileExists(t, filepath.Join(destDir, "keep-me.txt"))
	content, err := os.ReadFile(filepath.Join(destDir, "keep-me.txt"))
	require.NoError(t, err)
	assert.Equal(t, "original", string(content))

	// No leftover staging/backup directories.
	parent := filepath.Dir(destDir)
	entries, err := os.ReadDir(parent)
	require.NoError(t, err)
	for _, e := range entries {
		assert.NotContains(t, e.Name(), "~skpm-staging-")
		assert.NotContains(t, e.Name(), "~skpm-backup-")
	}
}

func TestAtomicExtractReplacesExistingDir(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "skill")
	require.NoError(t, os.MkdirAll(destDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "old.txt"), []byte("old"), 0o644))

	zipPath := buildZIP(t, map[string]string{"SKILL.md": "new"})
	require.NoError(t, archive.AtomicExtract(zipPath, destDir, ""))

	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
	assert.NoFileExists(t, filepath.Join(destDir, "old.txt"))
}

// ── CommonPrefix / *StrippingCommonPrefix ───────────────────────────────

func TestCommonPrefix(t *testing.T) {
	zipPath := buildZIP(t, map[string]string{
		"owner-repo-abc123/SKILL.md":   "x",
		"owner-repo-abc123/skill.yaml": "y",
	})
	r, err := zip.OpenReader(zipPath)
	require.NoError(t, err)
	defer r.Close()
	assert.Equal(t, "owner-repo-abc123/", archive.CommonPrefix(r.File))
}

func TestCommonPrefixNoneWhenEntriesDiffer(t *testing.T) {
	zipPath := buildZIP(t, map[string]string{
		"a/file.txt": "x",
		"b/file.txt": "y",
	})
	r, err := zip.OpenReader(zipPath)
	require.NoError(t, err)
	defer r.Close()
	assert.Equal(t, "", archive.CommonPrefix(r.File))
}

func TestExtractStrippingCommonPrefix(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	zipPath := buildZIP(t, map[string]string{
		"wrapper/SKILL.md":   "x",
		"wrapper/skill.yaml": "y",
	})
	require.NoError(t, archive.ExtractStrippingCommonPrefix(zipPath, destDir))
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
	assert.FileExists(t, filepath.Join(destDir, "skill.yaml"))
}

func TestAtomicExtractStrippingCommonPrefix(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	zipPath := buildZIP(t, map[string]string{
		"wrapper/SKILL.md": "x",
	})
	require.NoError(t, archive.AtomicExtractStrippingCommonPrefix(zipPath, destDir))
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
}

// ── Malformed input ──────────────────────────────────────────────────────

func TestExtractInvalidZIP(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "bad.zip")
	require.NoError(t, os.WriteFile(bad, []byte("not a zip"), 0o644))
	err := archive.Extract(bad, filepath.Join(t.TempDir(), "dest"), "")
	assert.Error(t, err)
}

func TestExtractRejectsCaseCollision(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	zipPath := buildZIP(t, map[string]string{
		"File.txt": "content 1",
		"file.txt": "content 2",
	})
	err := archive.Extract(zipPath, destDir, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "case collision")
}

func TestExtractRejectsDeviceFiles(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "installed")
	f, err := os.CreateTemp(t.TempDir(), "device-*.zip")
	require.NoError(t, err)
	zw := zip.NewWriter(f)
	hdr := &zip.FileHeader{Name: "dev-pipe", Method: zip.Deflate}
	hdr.SetMode(fs.ModeNamedPipe | 0o666)
	w, err := zw.CreateHeader(hdr)
	require.NoError(t, err)
	_, _ = w.Write([]byte("data"))
	require.NoError(t, zw.Close())
	f.Close()

	err = archive.Extract(f.Name(), destDir, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "refusing non-regular entry mode")
}
