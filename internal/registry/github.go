package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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
				Artifact: asset.GetName(),
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

func (r *GitHubRegistry) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	name := req.Manifest.Name
	version := req.Manifest.Version
	tag := formatReleaseTag(name, version, req.TagFormat)

	rel, err := r.getOrCreateRelease(ctx, tag, name, version, req.SHA256)
	if err != nil {
		return nil, err
	}

	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	downloadURL, err := r.uploadAsset(ctx, rel.GetID(), assetName, req.ArtifactPath)
	if err != nil {
		return nil, err
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

func (r *GitHubRegistry) getOrCreateRelease(ctx context.Context, tag, name, version, sha256 string) (*github.RepositoryRelease, error) {
	rel, resp, err := r.client.Repositories.GetReleaseByTag(ctx, r.org, r.repo, tag)
	if err == nil {
		return rel, nil
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("github: get release %s: %w", tag, err)
	}

	body := fmt.Sprintf("SHA256: `%s`\n\nAdd to `agent-skills.lock`:\n```yaml\n- name: %s\n  version: %s\n  source: github\n  sha256: %s\n```",
		sha256, name, version, sha256)

	rel, _, err = r.client.Repositories.CreateRelease(ctx, r.org, r.repo, &github.RepositoryRelease{
		TagName: github.String(tag),
		Name:    github.String(fmt.Sprintf("%s %s", name, version)),
		Body:    github.String(body),
	})
	if err != nil {
		return nil, fmt.Errorf("github: create release %s: %w", tag, err)
	}
	return rel, nil
}

func (r *GitHubRegistry) uploadAsset(ctx context.Context, releaseID int64, assetName, zipPath string) (string, error) {
	f, err := os.Open(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer f.Close()

	asset, _, err := r.client.Repositories.UploadReleaseAsset(ctx, r.org, r.repo, releaseID,
		&github.UploadOptions{Name: assetName}, f)
	if err != nil {
		return "", fmt.Errorf("github: upload asset %s: %w", assetName, err)
	}
	return asset.GetBrowserDownloadURL(), nil
}

// formatReleaseTag is shared by the github and gitlab backends: "prefixed"
// (the default) yields <name>/v<ver>, "plain" yields v<ver>.
func formatReleaseTag(name, version, tagFormat string) string {
	if tagFormat == "plain" {
		return "v" + version
	}
	return name + "/v" + version
}
