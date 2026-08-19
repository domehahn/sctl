package registry

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/httpclient"
	gitlab "github.com/xanzy/go-gitlab"
)

func init() {
	Register("gitlab", func(name string, rc config.RegistryConfig) (Registry, error) {
		if rc.Project == "" {
			return nil, fmt.Errorf("factory: gitlab registry %q requires project set to namespace/project (e.g. \"platform/agent-skills\")", name)
		}
		reg, err := NewGitLabRegistry(rc.URL, rc.Project, authToken(rc))
		if err != nil {
			return nil, err
		}
		return reg.WithName(name), nil
	})
}

type GitLabRegistry struct {
	client    *gitlab.Client
	name      string
	projectID string
	baseURL   string
	token     string
}

// NewGitLabRegistry creates a registry backed by GitLab Releases.
// projectID is "namespace/project" or a numeric project ID.
func NewGitLabRegistry(baseURL, projectID, token string) (*GitLabRegistry, error) {
	opts := []gitlab.ClientOptionFunc{gitlab.WithHTTPClient(httpclient.New())}
	if baseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(baseURL))
	}
	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("gitlab: create client: %w", err)
	}
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	return &GitLabRegistry{client: client, name: "gitlab", projectID: projectID, baseURL: strings.TrimRight(baseURL, "/"), token: token}, nil
}

func (r *GitLabRegistry) Type() string { return "gitlab" }

func (r *GitLabRegistry) Name() string { return r.name }

func (r *GitLabRegistry) WithName(name string) *GitLabRegistry {
	r.name = name
	return r
}

func (r *GitLabRegistry) Capabilities(context.Context) (*RegistryCapabilities, error) {
	return &RegistryCapabilities{Resolve: true, Download: true, Checksums: false}, nil
}

func (r *GitLabRegistry) Resolve(ctx context.Context, req ResolveRequest) (*ResolvedArtifact, error) {
	name := req.Ref.Name
	version := req.Constraint
	var releases []*gitlab.Release
	var err error

	if version == "" || version == "latest" {
		releases, _, err = r.client.Releases.ListReleases(r.projectID, &gitlab.ListReleasesOptions{
			ListOptions: gitlab.ListOptions{PerPage: 1},
		}, gitlab.WithContext(ctx))
	} else {
		tag := name + "/v" + version
		var rel *gitlab.Release
		rel, _, err = r.client.Releases.GetRelease(r.projectID, tag, gitlab.WithContext(ctx))
		if err != nil {
			tag = name + "-v" + version
			rel, _, err = r.client.Releases.GetRelease(r.projectID, tag, gitlab.WithContext(ctx))
		}
		if err == nil {
			releases = []*gitlab.Release{rel}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("gitlab: get release for %s@%s: %w", name, version, err)
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("gitlab: no release found for %s@%s", name, version)
	}

	rel := releases[0]
	ver := extractVersion(rel.TagName, name)
	assetName := fmt.Sprintf("%s-%s.zip", name, ver)

	for _, link := range rel.Assets.Links {
		if link.Name == assetName {
			return &ResolvedArtifact{
				Namespace:    req.Ref.Namespace,
				Name:         name,
				Version:      ver,
				Registry:     r.name,
				RegistryType: r.Type(),
				DownloadURL:  link.URL,
				Artifact:     link.Name,
				PackageType:  "zip",
			}, nil
		}
	}
	return nil, fmt.Errorf("gitlab: no asset %s found in release %s", assetName, rel.TagName)
}

func (r *GitLabRegistry) Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artifact.DownloadURL, nil)
	if err != nil {
		return err
	}
	resp, err := sharedHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("gitlab: download %s: %w", artifact.DownloadURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gitlab: download %s: HTTP %d", artifact.DownloadURL, resp.StatusCode)
	}
	_, err = httpclient.CopyLimited(dest, resp.Body)
	return err
}

func extractGitLabVersion(tag, name string) string {
	tag = strings.TrimPrefix(tag, name+"/v")
	tag = strings.TrimPrefix(tag, name+"-v")
	tag = strings.TrimPrefix(tag, "v")
	return tag
}

var _ = extractGitLabVersion

func (r *GitLabRegistry) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	name := req.Manifest.Name
	version := req.Manifest.Version
	tag := formatReleaseTag(name, version, req.TagFormat)
	assetName := fmt.Sprintf("%s-%s.zip", name, version)

	downloadURL, err := r.uploadGenericPackage(ctx, name, version, assetName, req.ArtifactPath)
	if err != nil {
		return nil, fmt.Errorf("gitlab: upload package: %w", err)
	}

	if err := r.createOrUpdateRelease(ctx, tag, name, version, req.SHA256, assetName, downloadURL); err != nil {
		return nil, fmt.Errorf("gitlab: create release: %w", err)
	}

	return &PublishResult{
		Name:        name,
		Version:     version,
		DownloadURL: downloadURL,
		SHA256:      req.SHA256,
		Registry:    r.name,
		Created:     true,
	}, nil
}

func (r *GitLabRegistry) uploadGenericPackage(ctx context.Context, pkgName, version, fileName, zipPath string) (string, error) {
	f, err := os.Open(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	// GitLab Generic Package API: PUT /projects/:id/packages/generic/:name/:version/:file
	encodedProject := strings.ReplaceAll(r.projectID, "/", "%2F")
	url := fmt.Sprintf("%s/api/v4/projects/%s/packages/generic/%s/%s/%s",
		r.baseURL, encodedProject, pkgName, version, fileName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("PRIVATE-TOKEN", r.token)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := sharedHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("upload HTTP %d", resp.StatusCode)
	}

	downloadURL := fmt.Sprintf("%s/api/v4/projects/%s/packages/generic/%s/%s/%s",
		r.baseURL, encodedProject, pkgName, version, fileName)
	return downloadURL, nil
}

func (r *GitLabRegistry) createOrUpdateRelease(ctx context.Context, tag, name, version, sha256, assetName, downloadURL string) error {
	description := fmt.Sprintf("SHA256: `%s`\n\nAdd to `agent-skills.lock`:\n```yaml\n- name: %s\n  version: %s\n  source: gitlab\n  sha256: %s\n```",
		sha256, name, version, sha256)

	links := []*gitlab.ReleaseAssetLinkOptions{{
		Name:     gitlab.String(assetName),
		URL:      gitlab.String(downloadURL),
		LinkType: gitlab.LinkType(gitlab.PackageLinkType),
	}}

	_, resp, err := r.client.Releases.GetRelease(r.projectID, tag, gitlab.WithContext(ctx))
	if err == nil {
		_, _, err = r.client.Releases.UpdateRelease(r.projectID, tag, &gitlab.UpdateReleaseOptions{
			Description: gitlab.String(description),
		}, gitlab.WithContext(ctx))
		return err
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		return err
	}

	_, _, err = r.client.Releases.CreateRelease(r.projectID, &gitlab.CreateReleaseOptions{
		Name:        gitlab.String(fmt.Sprintf("%s %s", name, version)),
		TagName:     gitlab.String(tag),
		Description: gitlab.String(description),
		Assets:      &gitlab.ReleaseAssetsOptions{Links: links},
	}, gitlab.WithContext(ctx))
	return err
}
