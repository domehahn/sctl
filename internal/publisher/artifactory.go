package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ArtifactoryPublisher struct {
	baseURL string
	repo    string
	token   string
}

func NewArtifactoryPublisher(baseURL, repo, token string) *ArtifactoryPublisher {
	return &ArtifactoryPublisher{
		baseURL: strings.TrimRight(baseURL, "/"),
		repo:    repo,
		token:   token,
	}
}

func (p *ArtifactoryPublisher) Publish(ctx context.Context, name, version, zipPath, sha256 string) (*Result, error) {
	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	uploadURL := fmt.Sprintf("%s/%s/%s/%s/%s", p.baseURL, p.repo, name, version, assetName)

	f, err := os.Open(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/zip")
	if sha256 != "" {
		req.Header.Set("X-Checksum-Sha256", sha256)
	}
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("artifactory: upload: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("artifactory: upload HTTP %d: %s", resp.StatusCode, string(body))
	}

	return &Result{
		Name:        name,
		Version:     version,
		DownloadURL: uploadURL,
		SHA256:      sha256,
	}, nil
}
