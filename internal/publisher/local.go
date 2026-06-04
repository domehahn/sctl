package publisher

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalPublisher struct {
	baseDir string
}

// NewLocalPublisher creates a publisher that copies artifacts to a local directory.
// Layout: baseDir/<name>/<version>/<name>-<version>.zip
func NewLocalPublisher(baseDir string) *LocalPublisher {
	return &LocalPublisher{baseDir: baseDir}
}

func (p *LocalPublisher) Publish(_ context.Context, name, version, zipPath, sha256 string) (*Result, error) {
	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	destDir := filepath.Join(p.baseDir, name, version)
	destPath := filepath.Join(destDir, assetName)

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, fmt.Errorf("local: create dir %s: %w", destDir, err)
	}

	if err := copyFile(zipPath, destPath); err != nil {
		return nil, fmt.Errorf("local: copy artifact: %w", err)
	}

	absPath, _ := filepath.Abs(destPath)
	return &Result{
		Name:        name,
		Version:     version,
		DownloadURL: "file://" + absPath,
		SHA256:      sha256,
	}, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}
