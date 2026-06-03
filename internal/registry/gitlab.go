package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

type GitLabRegistry struct {
	client    *gitlab.Client
	projectID string
}

// NewGitLabRegistry creates a registry backed by GitLab Releases.
// projectID is "namespace/project" or a numeric project ID.
func NewGitLabRegistry(baseURL, projectID, token string) (*GitLabRegistry, error) {
	opts := []gitlab.ClientOptionFunc{}
	if baseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(baseURL))
	}
	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("gitlab: create client: %w", err)
	}
	return &GitLabRegistry{client: client, projectID: projectID}, nil
}

func (r *GitLabRegistry) Resolve(ctx context.Context, name, version string) (*ResolvedArtifact, error) {
	var releases []*gitlab.Release
	var err error

	if version == "" {
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
				Name:        name,
				Version:     ver,
				DownloadURL: link.URL,
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("gitlab: download %s: %w", artifact.DownloadURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gitlab: download %s: HTTP %d", artifact.DownloadURL, resp.StatusCode)
	}
	_, err = io.Copy(dest, resp.Body)
	return err
}

func extractGitLabVersion(tag, name string) string {
	tag = strings.TrimPrefix(tag, name+"/v")
	tag = strings.TrimPrefix(tag, name+"-v")
	tag = strings.TrimPrefix(tag, "v")
	return tag
}

var _ = extractGitLabVersion
