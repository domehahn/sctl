package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

type GitLabPublisher struct {
	client    *gitlab.Client
	projectID string
	token     string
	baseURL   string
	tagFormat string
}

func NewGitLabPublisher(baseURL, projectID, token, tagFormat string) (*GitLabPublisher, error) {
	opts := []gitlab.ClientOptionFunc{}
	if baseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(baseURL))
	}
	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("gitlab publisher: create client: %w", err)
	}
	if tagFormat == "" {
		tagFormat = "prefixed"
	}
	return &GitLabPublisher{
		client:    client,
		projectID: projectID,
		token:     token,
		baseURL:   strings.TrimRight(baseURL, "/"),
		tagFormat: tagFormat,
	}, nil
}

func (p *GitLabPublisher) Publish(ctx context.Context, name, version, zipPath, sha256 string) (*Result, error) {
	tag := p.formatTag(name, version)
	assetName := fmt.Sprintf("%s-%s.zip", name, version)

	// Upload to GitLab Package Registry (generic packages).
	downloadURL, err := p.uploadGenericPackage(ctx, name, version, assetName, zipPath)
	if err != nil {
		return nil, fmt.Errorf("gitlab: upload package: %w", err)
	}

	// Create or update release pointing to the package.
	if err := p.createOrUpdateRelease(ctx, tag, name, version, sha256, assetName, downloadURL); err != nil {
		return nil, fmt.Errorf("gitlab: create release: %w", err)
	}

	return &Result{
		Name:        name,
		Version:     version,
		DownloadURL: downloadURL,
		SHA256:      sha256,
		Tag:         tag,
	}, nil
}

func (p *GitLabPublisher) uploadGenericPackage(ctx context.Context, pkgName, version, fileName, zipPath string) (string, error) {
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
	encodedProject := strings.ReplaceAll(p.projectID, "/", "%2F")
	url := fmt.Sprintf("%s/api/v4/projects/%s/packages/generic/%s/%s/%s",
		p.baseURL, encodedProject, pkgName, version, fileName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("PRIVATE-TOKEN", p.token)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("upload HTTP %d", resp.StatusCode)
	}

	downloadURL := fmt.Sprintf("%s/api/v4/projects/%s/packages/generic/%s/%s/%s",
		p.baseURL, encodedProject, pkgName, version, fileName)
	return downloadURL, nil
}

func (p *GitLabPublisher) createOrUpdateRelease(ctx context.Context, tag, name, version, sha256, assetName, downloadURL string) error {
	description := fmt.Sprintf("SHA256: `%s`\n\nAdd to `agent-skills.lock`:\n```yaml\n- name: %s\n  version: %s\n  source: gitlab\n  sha256: %s\n```",
		sha256, name, version, sha256)

	links := []*gitlab.ReleaseAssetLinkOptions{{
		Name:     gitlab.String(assetName),
		URL:      gitlab.String(downloadURL),
		LinkType: gitlab.LinkType(gitlab.PackageLinkType),
	}}

	_, resp, err := p.client.Releases.GetRelease(p.projectID, tag, gitlab.WithContext(ctx))
	if err == nil {
		_, _, err = p.client.Releases.UpdateRelease(p.projectID, tag, &gitlab.UpdateReleaseOptions{
			Description: gitlab.String(description),
		}, gitlab.WithContext(ctx))
		return err
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		return err
	}

	_, _, err = p.client.Releases.CreateRelease(p.projectID, &gitlab.CreateReleaseOptions{
		Name:        gitlab.String(fmt.Sprintf("%s %s", name, version)),
		TagName:     gitlab.String(tag),
		Description: gitlab.String(description),
		Assets:      &gitlab.ReleaseAssetsOptions{Links: links},
	}, gitlab.WithContext(ctx))
	return err
}

func (p *GitLabPublisher) formatTag(name, version string) string {
	if p.tagFormat == "plain" {
		return "v" + version
	}
	return name + "/v" + version
}
