package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type GitHubRegistry struct {
	client *github.Client
	name   string
	org    string
	repo   string
	token  string
}

// NewGitHubRegistry creates a registry backed by GitHub Releases.
// repoSlug must be in "owner/repo" format.
func NewGitHubRegistry(repoSlug, token string) (*GitHubRegistry, error) {
	parts := strings.SplitN(repoSlug, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("github registry: repoSlug must be owner/repo, got %q", repoSlug)
	}

	var httpClient *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = oauth2.NewClient(context.Background(), ts)
	}
	return &GitHubRegistry{
		client: github.NewClient(httpClient),
		name:   "github",
		org:    parts[0],
		repo:   parts[1],
		token:  token,
	}, nil
}

func (r *GitHubRegistry) Type() string { return "github" }

func (r *GitHubRegistry) Name() string { return r.name }

func (r *GitHubRegistry) WithName(name string) *GitHubRegistry {
	r.name = name
	return r
}

func (r *GitHubRegistry) Capabilities(context.Context) (*RegistryCapabilities, error) {
	return &RegistryCapabilities{Resolve: true, Download: true, Checksums: false}, nil
}

func (r *GitHubRegistry) Resolve(ctx context.Context, req ResolveRequest) (*ResolvedArtifact, error) {
	name := req.Ref.Name
	version := req.Constraint
	var release *github.RepositoryRelease
	var err error

	if version == "" || version == "latest" {
		release, _, err = r.client.Repositories.GetLatestRelease(ctx, r.org, r.repo)
	} else {
		release, _, err = r.client.Repositories.GetReleaseByTag(ctx, r.org, r.repo, name+"/v"+version)
		if err != nil {
			release, _, err = r.client.Repositories.GetReleaseByTag(ctx, r.org, r.repo, name+"-v"+version)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("github: get release for %s@%s: %w", name, version, err)
	}

	assetName := fmt.Sprintf("%s-%s.zip", name, release.GetTagName())
	assetName = strings.TrimPrefix(assetName, name+"-"+name+"/v")

	for _, asset := range release.Assets {
		if asset.GetName() == fmt.Sprintf("%s-%s.zip", name, extractVersion(release.GetTagName(), name)) {
			return &ResolvedArtifact{
				Namespace:    req.Ref.Namespace,
				Name:         name,
				Version:      extractVersion(release.GetTagName(), name),
				Registry:     r.name,
				RegistryType: r.Type(),
				DownloadURL:  asset.GetBrowserDownloadURL(),
				ArtifactName: asset.GetName(),
				PackageType:  "zip",
			}, nil
		}
	}
	return nil, fmt.Errorf("github: no asset matching %s-*.zip found in release %s", name, release.GetTagName())
}

func (r *GitHubRegistry) Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artifact.DownloadURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("github: download %s: %w", artifact.DownloadURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github: download %s: HTTP %d", artifact.DownloadURL, resp.StatusCode)
	}
	_, err = io.Copy(dest, resp.Body)
	return err
}

func extractVersion(tag, name string) string {
	tag = strings.TrimPrefix(tag, name+"/v")
	tag = strings.TrimPrefix(tag, name+"-v")
	tag = strings.TrimPrefix(tag, "v")
	return tag
}
