package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

type ArtifactoryRegistry struct {
	name       string
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
		name:       "artifactory",
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

func (r *ArtifactoryRegistry) Type() string { return "artifactory" }

func (r *ArtifactoryRegistry) Name() string { return r.name }

func (r *ArtifactoryRegistry) WithName(name string) *ArtifactoryRegistry {
	r.name = name
	return r
}

func (r *ArtifactoryRegistry) Capabilities(context.Context) (*RegistryCapabilities, error) {
	return &RegistryCapabilities{Resolve: true, Download: true, Info: false, SemVerConstraints: true, Checksums: false}, nil
}

func (r *ArtifactoryRegistry) Resolve(ctx context.Context, req ResolveRequest) (*ResolvedArtifact, error) {
	name := req.Ref.Name
	version := req.Constraint
	if version == "" || version == "latest" || strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") || strings.Contains(version, " ") {
		versions, err := r.ListVersions(ctx, req.Ref)
		if err != nil {
			return nil, err
		}
		version, err = SelectVersion(versions, req.Constraint, ConstraintOptions{
			IncludePrerelease: req.IncludePrerelease,
			AllowDeprecated:   req.AllowDeprecated,
			AllowYanked:       req.AllowYanked,
		})
		if err != nil {
			return nil, err
		}
	}

	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	downloadURL := fmt.Sprintf("%s/%s/%s/%s/%s", r.baseURL, r.repo, name, version, assetName)

	return &ResolvedArtifact{
		Namespace:    req.Ref.Namespace,
		Name:         name,
		Version:      version,
		Registry:     r.name,
		RegistryType: r.Type(),
		DownloadURL:  downloadURL,
		Artifact: assetName,
		PackageType:  "zip",
	}, nil
}

func (r *ArtifactoryRegistry) ListVersions(ctx context.Context, ref SkillRef) ([]VersionInfo, error) {
	url := fmt.Sprintf("%s/api/storage/%s/%s", r.baseURL, r.repo, ref.Name)
	req, err := r.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("artifactory: list versions for %s: %w", ref.Name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artifactory: list versions HTTP %d", resp.StatusCode)
	}

	var info artifactoryFolderInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("artifactory: parse folder info: %w", err)
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
		return nil, fmt.Errorf("artifactory: no versions found for %s", ref.Name)
	}

	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare("v"+versions[i], "v"+versions[j]) > 0
	})
	out := make([]VersionInfo, 0, len(versions))
	for _, v := range versions {
		out = append(out, VersionInfo{Version: v})
	}
	return out, nil
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

func (r *ArtifactoryRegistry) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	name := req.Manifest.Name
	version := req.Manifest.Version
	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	uploadURL := fmt.Sprintf("%s/%s/%s/%s/%s", r.baseURL, r.repo, name, version, assetName)

	f, err := os.Open(req.ArtifactPath)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/zip")
	if req.SHA256 != "" {
		httpReq.Header.Set("X-Checksum-Sha256", req.SHA256)
	}
	if r.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+r.token)
	}

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("artifactory: upload: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("artifactory: upload HTTP %d: %s", resp.StatusCode, string(body))
	}

	return &PublishResult{
		Name:        name,
		Version:     version,
		DownloadURL: uploadURL,
		SHA256:      req.SHA256,
		Registry:    r.name,
		Created:     true,
	}, nil
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
