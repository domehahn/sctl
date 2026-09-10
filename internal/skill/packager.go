package skill

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/domehahn/sklib/packageio"
	"github.com/domehahn/sklib/spec"
)

type PackageResult struct {
	OutputPath string
	SHA256     string
	Name       string
	Version    string
}

type Packager struct {
	validator Validator
}

func NewPackager() *Packager {
	return &Packager{validator: NewValidatorWithOptions(ValidationOptions{Strict: true})}
}

func (p *Packager) Package(ctx context.Context, dir, outputDir string) (*PackageResult, error) {
	res, err := p.validator.Validate(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if !res.Valid {
		msgs := make([]string, len(res.Errors))
		for i, e := range res.Errors {
			msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
		}
		return nil, &ValidationFailedError{Errors: msgs}
	}

	sy, err := readSkillYAML(dir)
	if err != nil {
		return nil, err
	}

	if outputDir == "" {
		outputDir = dir
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}
	outName := fmt.Sprintf("%s-%s.zip", sy.Name, sy.Version)
	outPath := filepath.Join(outputDir, outName)

	sha, err := buildZIP(ctx, dir, outPath, sy)
	if err != nil {
		return nil, err
	}

	return &PackageResult{
		OutputPath: outPath,
		SHA256:     sha,
		Name:       sy.Name,
		Version:    sy.Version,
	}, nil
}

func buildZIP(_ context.Context, srcDir, outPath string, sy *SkillYAML) (string, error) {
	tmp := outPath + ".tmp"
	defer os.Remove(tmp)

	f, err := os.Create(tmp)
	if err != nil {
		return "", fmt.Errorf("create zip: %w", err)
	}

	h := sha256.New()
	mw := io.MultiWriter(f, h)
	zw := zip.NewWriter(mw)

	checksums, err := addDirToZIP(zw, srcDir, srcDir)
	if err != nil {
		f.Close()
		return "", err
	}

	// Collect sorted file list for manifest.json.
	sortedFiles := make([]string, 0, len(checksums))
	for p := range checksums {
		sortedFiles = append(sortedFiles, p)
	}
	sort.Strings(sortedFiles)
	sortedFiles = append(sortedFiles, "manifest.json", "checksums.txt")

	manifest := SkillManifest{
		SpecVersion:    1,
		Name:           sy.Name,
		Namespace:      spec.DefaultNamespace(sy.Namespace),
		Version:        sy.Version,
		Description:    sy.Description,
		Entrypoint:     spec.DefaultEntrypoint(sy.Entrypoint),
		CompatibleWith: sy.CompatibleWith,
		PackageType:    "zip",
		License:        sy.License,
		PackagedBy:     "skpm",
		PackagedAt:     reproducibleTime().Format(time.RFC3339),
		Files:          sortedFiles,
	}
	if commit := os.Getenv("GIT_COMMIT"); commit != "" {
		manifest.SourceCommit = commit
	}

	mData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		f.Close()
		return "", fmt.Errorf("marshal manifest: %w", err)
	}
	if err := addBytesToZIP(zw, "manifest.json", mData); err != nil {
		f.Close()
		return "", err
	}
	checksums["manifest.json"] = sha256Hex(mData)

	// Build checksums.txt using sklib/packageio canonical format: "<sha256>  <path>".
	var csEntries []spec.ChecksumEntry
	for _, p := range append(sortedFiles[:len(sortedFiles)-1], "manifest.json") {
		if sum, ok := checksums[p]; ok {
			csEntries = append(csEntries, spec.ChecksumEntry{Path: p, SHA256: sum})
		}
	}
	sort.Slice(csEntries, func(i, j int) bool { return csEntries[i].Path < csEntries[j].Path })
	checksumData := packageio.FormatChecksumsText(csEntries)
	if err := addBytesToZIP(zw, "checksums.txt", checksumData); err != nil {
		f.Close()
		return "", err
	}

	if err := zw.Close(); err != nil {
		f.Close()
		return "", fmt.Errorf("close zip: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close zip file: %w", err)
	}

	sha := hex.EncodeToString(h.Sum(nil))

	if err := os.Rename(tmp, outPath); err != nil {
		return "", fmt.Errorf("atomic rename zip: %w", err)
	}
	return sha, nil
}

func addDirToZIP(zw *zip.Writer, baseDir, currentDir string) (map[string]string, error) {
	var files []string
	if err := filepath.WalkDir(currentDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(baseDir, path)
		relPath = filepath.ToSlash(relPath)
		if relPath == "." {
			return nil
		}
		if relPath == "manifest.json" || relPath == "checksums.txt" {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink is not allowed in a package: %s", relPath)
		}
		if shouldSkip(entry.Name(), relPath) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk %s: %w", currentDir, err)
	}
	sort.Strings(files)
	checksums := map[string]string{}
	for _, fullPath := range files {
		relPath, _ := filepath.Rel(baseDir, fullPath)
		relPath = filepath.ToSlash(relPath)
		sum, err := addFileToZIP(zw, fullPath, relPath)
		if err != nil {
			return nil, err
		}
		checksums[relPath] = sum
	}
	return checksums, nil
}

func addFileToZIP(zw *zip.Writer, fullPath, relPath string) (string, error) {
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", fullPath, err)
	}
	if err := addBytesToZIP(zw, relPath, data); err != nil {
		return "", err
	}
	return sha256Hex(data), nil
}

func addBytesToZIP(zw *zip.Writer, relPath string, data []byte) error {
	header := &zip.FileHeader{
		Name:     relPath,
		Method:   zip.Deflate,
		Modified: reproducibleTime(),
	}
	header.SetMode(0o644)
	dst, err := zw.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create zip entry %s: %w", relPath, err)
	}
	if _, err := dst.Write(data); err != nil {
		return fmt.Errorf("write zip entry %s: %w", relPath, err)
	}
	return nil
}

func shouldSkip(name, relPath string) bool {
	skip := []string{".git", ".DS_Store", "*.tmp", "*.part"}
	for _, pattern := range skip {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return strings.HasPrefix(relPath, ".git/")
}

func readSkillYAML(dir string) (*SkillYAML, error) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read skill.yaml: %w", err)
	}
	sy, _, err := decodeSkillMetadata(data)
	if err != nil {
		return nil, fmt.Errorf("parse skill.yaml: %w", err)
	}
	return sy, nil
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func reproducibleTime() time.Time {
	return time.Unix(0, 0).UTC()
}

type ValidationFailedError struct {
	Errors []string
}

func (e *ValidationFailedError) Error() string {
	return fmt.Sprintf("skill validation failed: %s", strings.Join(e.Errors, "; "))
}
