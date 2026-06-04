package publisher

import "context"

// Result holds the outcome of a successful publish.
type Result struct {
	Name        string
	Version     string
	DownloadURL string
	SHA256      string
	TagCreated  bool
	Tag         string
}

// Publisher uploads a packaged skill artifact to a registry.
type Publisher interface {
	// Publish uploads zipPath and returns the public download URL.
	Publish(ctx context.Context, name, version, zipPath, sha256 string) (*Result, error)
}
