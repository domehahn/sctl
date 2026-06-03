package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

type ArtifactoryRegistry struct {
	baseURL    string
	repo       string
	token      string
	httpClient *http.Client
}

// NewArtifactoryRegistry creates a registry backed by JFrog Artifactory Generic Repository.
// baseURL is like "https://artifactory.company.com/artifactory".
// repo is the repository name like "agent-skills".
func NewArtifactoryRegistry(baseURL, repo, token string) *ArtifactoryRegistry {
	return &ArtifactoryRegistry{
		baseURL:    strings.TrimRight(baseURL, "/"),
		repo:       repo,
		token:      token,
		httpClient: &http.Client{},
	}
}

type artifactoryFolderInfo struct {
	Children []struct {
		URI    string `json:"uri"`
		Folder bool   `json:"folder"`
	} `json:"children"`
}

func (r *ArtifactoryRegistry) Resolve(ctx context.Context, name, version string) (*ResolvedArtifact, error) {
	if version == "" {
		var err error
		version, err = r.resolveLatestVersion(ctx, name)
		if err != nil {
			return nil, err
		}
	}

	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	downloadURL := fmt.Sprintf("%s/%s/%s/%s/%s", r.baseURL, r.repo, name, version, assetName)

	return &ResolvedArtifact{
		Name:        name,
		Version:     version,
		DownloadURL: downloadURL,
	}, nil
}

func (r *ArtifactoryRegistry) resolveLatestVersion(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s/api/storage/%s/%s", r.baseURL, r.repo, name)
	req, err := r.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return "", err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("artifactory: list versions for %s: %w", name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("artifactory: list versions HTTP %d", resp.StatusCode)
	}

	var info artifactoryFolderInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("artifactory: parse folder info: %w", err)
	}

	versions := make([]string, 0, len(info.Children))
	for _, child := range info.Children {
		if child.Folder {
			v := strings.Trim(child.URI, "/")
			if semver.IsValid("v" + v) {
				versions = append(versions, v)
			}
		}
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("artifactory: no versions found for %s", name)
	}

	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare("v"+versions[i], "v"+versions[j]) > 0
	})
	return versions[0], nil
}

func (r *ArtifactoryRegistry) Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error {
	req, err := r.newRequest(ctx, http.MethodGet, artifact.DownloadURL)
	if err != nil {
		return err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("artifactory: download %s: %w", artifact.DownloadURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("artifactory: download HTTP %d", resp.StatusCode)
	}
	_, err = io.Copy(dest, resp.Body)
	return err
}

func (r *ArtifactoryRegistry) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	if r.token != "" {
		req.Header.Set("Authorization", "Bearer "+r.token)
	}
	return req, nil
}
