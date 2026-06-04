package publisher

import "context"

var (
	CopyFile = copyFile
)

func (p *GitLabPublisher) UploadGenericPackage(ctx context.Context, pkgName, version, fileName, zipPath string) (string, error) {
	return p.uploadGenericPackage(ctx, pkgName, version, fileName, zipPath)
}

func (p *GitLabPublisher) CreateOrUpdateRelease(ctx context.Context, tag, name, version, sha256, assetName, downloadURL string) error {
	return p.createOrUpdateRelease(ctx, tag, name, version, sha256, assetName, downloadURL)
}

func (p *GitLabPublisher) FormatTag(name, version string) string {
	return p.formatTag(name, version)
}

func (p *GitHubPublisher) FormatTag(name, version string) string {
	return p.formatTag(name, version)
}
