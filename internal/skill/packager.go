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
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	return &Packager{validator: NewValidator()}
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

	if err := addDirToZIP(zw, srcDir, srcDir); err != nil {
		f.Close()
		return "", err
	}

	manifest := SkillManifest{
		Name:           sy.Name,
		Version:        sy.Version,
		SHA256:         "",
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		CompatibleWith: sy.CompatibleWith,
	}
	if commit := os.Getenv("GIT_COMMIT"); commit != "" {
		manifest.SourceCommit = commit
	}

	mData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		f.Close()
		return "", fmt.Errorf("marshal manifest: %w", err)
	}
	mw2, err := zw.Create("manifest.json")
	if err != nil {
		f.Close()
		return "", fmt.Errorf("create manifest entry: %w", err)
	}
	if _, err := mw2.Write(mData); err != nil {
		f.Close()
		return "", fmt.Errorf("write manifest: %w", err)
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

func addDirToZIP(zw *zip.Writer, baseDir, currentDir string) error {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", currentDir, err)
	}
	for _, entry := range entries {
		fullPath := filepath.Join(currentDir, entry.Name())
		relPath, _ := filepath.Rel(baseDir, fullPath)
		relPath = filepath.ToSlash(relPath)

		if shouldSkip(entry.Name(), relPath) {
			continue
		}

		if entry.IsDir() {
			if err := addDirToZIP(zw, baseDir, fullPath); err != nil {
				return err
			}
			continue
		}

		if err := addFileToZIP(zw, fullPath, relPath); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZIP(zw *zip.Writer, fullPath, relPath string) error {
	src, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", fullPath, err)
	}
	defer src.Close()

	dst, err := zw.Create(relPath)
	if err != nil {
		return fmt.Errorf("create zip entry %s: %w", relPath, err)
	}
	_, err = io.Copy(dst, src)
	return err
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
	var sy SkillYAML
	if err := yaml.Unmarshal(data, &sy); err != nil {
		return nil, fmt.Errorf("parse skill.yaml: %w", err)
	}
	return &sy, nil
}

type ValidationFailedError struct {
	Errors []string
}

func (e *ValidationFailedError) Error() string {
	return fmt.Sprintf("skill validation failed: %s", strings.Join(e.Errors, "; "))
}
