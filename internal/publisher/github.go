package publisher

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type GitHubPublisher struct {
	client    *github.Client
	org       string
	repo      string
	tagFormat string // "prefixed" (<name>/v<ver>) or "plain" (v<ver>)
}

func NewGitHubPublisher(repoSlug, token, tagFormat string) (*GitHubPublisher, error) {
	parts := strings.SplitN(repoSlug, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("github publisher: repoSlug must be owner/repo, got %q", repoSlug)
	}

	var httpClient *oauth2.Transport
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = &oauth2.Transport{Source: ts}
		_ = httpClient
	}
	var client *github.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = github.NewClient(oauth2.NewClient(context.Background(), ts))
	} else {
		client = github.NewClient(nil)
	}

	if tagFormat == "" {
		tagFormat = "prefixed"
	}
	return &GitHubPublisher{
		client:    client,
		org:       parts[0],
		repo:      parts[1],
		tagFormat: tagFormat,
	}, nil
}

func (p *GitHubPublisher) Publish(ctx context.Context, name, version, zipPath, sha256 string) (*Result, error) {
	tag := p.formatTag(name, version)

	// Get or create release for this tag.
	rel, err := p.getOrCreateRelease(ctx, tag, name, version, sha256)
	if err != nil {
		return nil, err
	}

	// Upload asset.
	assetName := fmt.Sprintf("%s-%s.zip", name, version)
	downloadURL, err := p.uploadAsset(ctx, rel.GetID(), assetName, zipPath)
	if err != nil {
		return nil, err
	}

	return &Result{
		Name:        name,
		Version:     version,
		DownloadURL: downloadURL,
		SHA256:      sha256,
		Tag:         tag,
	}, nil
}

func (p *GitHubPublisher) getOrCreateRelease(ctx context.Context, tag, name, version, sha256 string) (*github.RepositoryRelease, error) {
	rel, resp, err := p.client.Repositories.GetReleaseByTag(ctx, p.org, p.repo, tag)
	if err == nil {
		return rel, nil
	}
	if resp == nil || resp.StatusCode != 404 {
		return nil, fmt.Errorf("github: get release %s: %w", tag, err)
	}

	body := fmt.Sprintf("SHA256: `%s`\n\nAdd to `agent-skills.lock`:\n```yaml\n- name: %s\n  version: %s\n  source: github\n  sha256: %s\n```",
		sha256, name, version, sha256)

	rel, _, err = p.client.Repositories.CreateRelease(ctx, p.org, p.repo, &github.RepositoryRelease{
		TagName: github.String(tag),
		Name:    github.String(fmt.Sprintf("%s %s", name, version)),
		Body:    github.String(body),
	})
	if err != nil {
		return nil, fmt.Errorf("github: create release %s: %w", tag, err)
	}
	return rel, nil
}

func (p *GitHubPublisher) uploadAsset(ctx context.Context, releaseID int64, assetName, zipPath string) (string, error) {
	f, err := os.Open(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer f.Close()

	asset, _, err := p.client.Repositories.UploadReleaseAsset(ctx, p.org, p.repo, releaseID,
		&github.UploadOptions{Name: assetName}, f)
	if err != nil {
		return "", fmt.Errorf("github: upload asset %s: %w", assetName, err)
	}
	return asset.GetBrowserDownloadURL(), nil
}

func (p *GitHubPublisher) formatTag(name, version string) string {
	if p.tagFormat == "plain" {
		return "v" + version
	}
	return name + "/v" + version
}
