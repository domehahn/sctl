package registry

import (
	"context"
	"io"
)

type ResolvedArtifact struct {
	Name           string
	Version        string
	DownloadURL    string
	SHA256         string
	CompatibleWith []string
}

type VersionInfo struct {
	Version        string   `json:"version"`
	Deprecated     bool     `json:"deprecated,omitempty"`
	Yanked         bool     `json:"yanked,omitempty"`
	CompatibleWith []string `json:"compatible_with,omitempty"`
}

type SkillSearchResult struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"`
}

type SkillInfo struct {
	Name           string        `json:"name"`
	Description    string        `json:"description,omitempty"`
	LatestVersion  string        `json:"latest_version,omitempty"`
	Versions       []VersionInfo `json:"versions,omitempty"`
	Source         string        `json:"source,omitempty"`
	CompatibleWith []string      `json:"compatible_with,omitempty"`
}

type PublishMetadata struct {
	Name        string
	Version     string
	Description string
}

type RegistryCapabilities struct {
	Resolve   bool `json:"resolve"`
	Download  bool `json:"download"`
	Search    bool `json:"search"`
	Info      bool `json:"info"`
	Publish   bool `json:"publish"`
	Deprecate bool `json:"deprecate"`
	Yank      bool `json:"yank"`
}

// Registry abstracts any source that can serve skill artifacts.
type Registry interface {
	// Resolve finds a matching version. version may be "" (latest) or an exact semver like "1.5.0".
	Resolve(ctx context.Context, name, version string) (*ResolvedArtifact, error)
	// Download streams the artifact bytes to dest.
	Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error
}

type DiscoveryRegistry interface {
	ListVersions(ctx context.Context, name string) ([]VersionInfo, error)
	Search(ctx context.Context, query string) ([]SkillSearchResult, error)
	Info(ctx context.Context, name string) (*SkillInfo, error)
	Capabilities(ctx context.Context) RegistryCapabilities
}

type PublishingRegistry interface {
	Publish(ctx context.Context, artifactPath string, metadata PublishMetadata) error
	Deprecate(ctx context.Context, name string, version string, reason string) error
	Yank(ctx context.Context, name string, version string, reason string) error
}
