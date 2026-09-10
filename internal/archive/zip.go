// Package archive provides the single, hardened ZIP extraction
// implementation used everywhere skpm unpacks an artifact that ultimately
// comes from a registry, git host, or a local file a user pointed at — none
// of which skpm should trust to be a well-formed archive. Before this
// package existed, skpm had four independent, hand-rolled ZIP extraction
// loops (internal/installer, internal/cli/add.go, internal/cli/clone.go +
// import_cmd.go, internal/registry/refdownload.go); two of them had no path
// containment check at all, letting a malicious archive entry (e.g.
// "../../../../home/user/.ssh/authorized_keys") write outside the intended
// destination directory during `skpm add` or a GitHub/GitLab --ref
// download. This package is the fix: every extraction path in skpm now
// goes through it, and there is exactly one place a future containment
// bypass could be introduced instead of four.
package archive

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Limits guard against malicious or malformed archives (zip bombs,
// resource-exhaustion via entry count) during extraction of untrusted
// content. They are deliberately generous for legitimate skill packages
// (which are source/config bundles, not binary blobs) while still bounding
// worst-case behavior.
const (
	MaxEntries        = 100_000
	MaxEntrySizeBytes = 200 * 1024 * 1024      // 200MB per file
	MaxTotalSizeBytes = 1 * 1024 * 1024 * 1024 // 1GB total decompressed
)

// ErrZipSlip is returned when an archive entry would resolve outside the
// extraction destination directory.
var ErrZipSlip = errors.New("archive: entry escapes destination directory (zip slip)")

// SafeJoin resolves a ZIP entry's name (which always uses '/' separators,
// per the ZIP format — see archive/zip's documentation of zip.File.Name)
// against destDir and returns the local filesystem path to extract it to.
// It rejects — does not silently sanitize — absolute paths, NUL bytes, and
// any name containing a ".." component that would climb above destDir: the
// classic zip-slip primitive. Rejecting outright (rather than e.g.
// remapping "../evil.txt" to "evil.txt") is deliberate: a malformed or
// hostile entry name is a signal the whole archive shouldn't be trusted,
// and silently "fixing" it would hide that signal from the caller.
func SafeJoin(destDir, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("archive: empty entry name")
	}
	if strings.ContainsRune(name, 0) {
		return "", fmt.Errorf("%w: entry name contains a NUL byte: %q", ErrZipSlip, name)
	}
	cleaned := path.Clean(name)
	if cleaned == "." {
		return "", fmt.Errorf("%w: entry resolves to destination root: %q", ErrZipSlip, name)
	}
	if path.IsAbs(cleaned) {
		return "", fmt.Errorf("%w: absolute entry path: %q", ErrZipSlip, name)
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("%w: %q", ErrZipSlip, name)
	}
	outPath := filepath.Join(destDir, filepath.FromSlash(cleaned))

	// Belt-and-suspenders re-check with filepath.Rel: Windows path
	// semantics (drive letters, backslashes, UNC paths) differ from the
	// POSIX-style ("/"-separated) cleaning above, so this catches anything
	// that slips past it on that platform.
	if !IsWithinDir(destDir, outPath) {
		return "", fmt.Errorf("%w: %s", ErrZipSlip, name)
	}
	return outPath, nil
}

// IsWithinDir reports whether target resolves to a path at or below base.
// It's the containment primitive SafeJoin uses internally, exported
// separately for callers checking an already-resolved path rather than a
// raw zip entry name.
func IsWithinDir(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// ExtractReader extracts every entry of an already-open *zip.Reader into
// destDir, stripping prefix (if non-empty) from each entry name first —
// used when unpacking a subdirectory out of a larger archive (a GitHub/
// GitLab source zipball, or an import bundle with a single top-level
// directory). destDir is created if it does not exist. Symlink entries are
// rejected outright: skill packages have no legitimate use for them, and
// they are a well-known secondary escape vector even when the symlink's
// own target path passes the containment check (the attack extracts a
// symlink pointing outside destDir, then a later "write through the link"
// entry).
func ExtractReader(zr *zip.Reader, destDir, prefix string) error {
	if len(zr.File) > MaxEntries {
		return fmt.Errorf("archive: %d entries exceeds the %d limit", len(zr.File), MaxEntries)
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("archive: create destination: %w", err)
	}

	forbiddenModes := fs.ModeSymlink | fs.ModeDevice | fs.ModeNamedPipe | fs.ModeSocket | fs.ModeCharDevice | fs.ModeIrregular | fs.ModeSetuid | fs.ModeSetgid
	seenPaths := make(map[string]string)
	var totalSize uint64

	for _, f := range zr.File {
		name := f.Name
		if prefix != "" {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			name = strings.TrimPrefix(name, prefix)
		}
		if name == "" {
			continue
		}

		if f.Mode()&forbiddenModes != 0 {
			return fmt.Errorf("archive: refusing non-regular entry mode %v for %q", f.Mode(), f.Name)
		}

		outPath, err := SafeJoin(destDir, name)
		if err != nil {
			return err
		}

		lowerPath := strings.ToLower(filepath.Clean(outPath))
		if existing, ok := seenPaths[lowerPath]; ok && existing != outPath {
			return fmt.Errorf("archive: refusing case collision between %q and %q", existing, outPath)
		}
		seenPaths[lowerPath] = outPath

		if f.CompressedSize64 > 0 && f.UncompressedSize64 > 1024*1024 {
			ratio := f.UncompressedSize64 / f.CompressedSize64
			if ratio > 1000 {
				return fmt.Errorf("archive: refusing zip bomb (extreme compression ratio %d:1) for %q", ratio, f.Name)
			}
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			}
			continue
		}

		if f.UncompressedSize64 > MaxEntrySizeBytes {
			return fmt.Errorf("archive: entry %q (%d bytes uncompressed) exceeds the %d byte per-file limit", f.Name, f.UncompressedSize64, MaxEntrySizeBytes)
		}
		totalSize += f.UncompressedSize64
		if totalSize > MaxTotalSizeBytes {
			return fmt.Errorf("archive: total uncompressed size exceeds the %d byte limit", MaxTotalSizeBytes)
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		if err := writeEntry(f, outPath); err != nil {
			return err
		}
	}
	return nil
}

func writeEntry(f *zip.File, outPath string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Defense in depth against a mismatched/forged UncompressedSize64
	// header: cap the actual bytes copied at one byte past the declared
	// per-file limit, so a truncated read (rather than an unbounded one)
	// is the failure mode even if the header lied.
	limited := io.LimitReader(src, MaxEntrySizeBytes+1)
	n, err := io.Copy(out, limited)
	if err != nil {
		return err
	}
	if n > MaxEntrySizeBytes {
		return fmt.Errorf("archive: entry %q exceeded the %d byte per-file limit during extraction", f.Name, MaxEntrySizeBytes)
	}
	return nil
}

// Extract opens zipPath and extracts it into destDir via ExtractReader.
func Extract(zipPath, destDir, prefix string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	return ExtractReader(&r.Reader, destDir, prefix)
}

// ExtractStrippingCommonPrefix opens zipPath, detects a shared top-level
// directory across all its entries via CommonPrefix, and extracts with
// that prefix stripped — the common shape for a skill ZIP whose entries
// are all wrapped in one top-level directory (a git-forge source archive,
// or a bundle produced by re-zipping a directory).
func ExtractStrippingCommonPrefix(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	prefix := CommonPrefix(r.File)
	return ExtractReader(&r.Reader, destDir, prefix)
}

// AtomicExtract extracts zipPath into a staging directory next to destDir,
// then atomically swaps it into place: any existing destDir is renamed
// aside as a backup, the staging directory is renamed to destDir, and the
// backup is removed only after that rename succeeds. On any failure,
// destDir is left as it was (the backup is restored) rather than partially
// overwritten. This is the extraction path for anything replacing an
// existing installed skill directory — internal/installer.Install,
// `skpm add`, and `skpm import`.
func AtomicExtract(zipPath, destDir, prefix string) error {
	suffix := strconv.Itoa(rand.Int())
	staging := destDir + "~skpm-staging-" + suffix
	backup := destDir + "~skpm-backup-" + suffix

	if err := Extract(zipPath, staging, prefix); err != nil {
		os.RemoveAll(staging)
		return fmt.Errorf("extract to staging: %w", err)
	}

	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, backup); err != nil {
			os.RemoveAll(staging)
			return fmt.Errorf("back up existing directory: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, destDir)
		return fmt.Errorf("create parent directory: %w", err)
	}

	if err := os.Rename(staging, destDir); err != nil {
		os.RemoveAll(staging)
		os.Rename(backup, destDir)
		return fmt.Errorf("rename staging to destination: %w", err)
	}

	os.RemoveAll(backup)
	return nil
}

// AtomicExtractStrippingCommonPrefix combines ExtractStrippingCommonPrefix's
// prefix detection with AtomicExtract's staged, atomic swap into destDir.
func AtomicExtractStrippingCommonPrefix(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	prefix := CommonPrefix(r.File)
	r.Close()
	return AtomicExtract(zipPath, destDir, prefix)
}

// CommonPrefix returns the shared top-level directory of every entry in
// files (e.g. "owner-repo-abc123/" for a GitHub source zipball), or "" if
// the entries don't share one. Used to strip a synthetic wrapper directory
// before extraction.
func CommonPrefix(files []*zip.File) string {
	if len(files) == 0 {
		return ""
	}
	first := files[0].Name
	slash := len(first)
	for i := range first {
		if first[i] == '/' {
			slash = i + 1
			break
		}
	}
	prefix := first[:slash]
	if prefix == "" {
		return ""
	}
	for _, f := range files[1:] {
		if !strings.HasPrefix(f.Name, prefix) {
			return ""
		}
	}
	return prefix
}
