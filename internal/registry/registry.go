package registry

import (
	"context"
	"io"
)

type ResolvedArtifact struct {
	Name        string
	Version     string
	DownloadURL string
	SHA256      string
}

// Registry abstracts any source that can serve skill artifacts.
type Registry interface {
	// Resolve finds a matching version. version may be "" (latest) or an exact semver like "1.5.0".
	Resolve(ctx context.Context, name, version string) (*ResolvedArtifact, error)
	// Download streams the artifact bytes to dest.
	Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error
}
