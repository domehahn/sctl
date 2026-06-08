package registry

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/domehahn/sklib/registryapi"
	"github.com/domehahn/skpm/v2/internal/skill"
)

type SkillRef struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Name      string `json:"name" yaml:"name"`
}

type SkillVersionRef struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Name      string `json:"name" yaml:"name"`
	Version   string `json:"version" yaml:"version"`
}

type ResolveRequest struct {
	Ref               SkillRef `json:"ref"`
	Constraint        string   `json:"constraint,omitempty"`
	IncludePrerelease bool     `json:"include_prerelease,omitempty"`
	AllowDeprecated   bool     `json:"allow_deprecated,omitempty"`
	AllowYanked       bool     `json:"allow_yanked,omitempty"`
}

type SearchRequest struct {
	Query     string `json:"query,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// ResolvedArtifact is the result of resolving a skill reference via a registry adapter.
// It carries skpm-internal fields (Registry, RegistryType, Metadata) beyond the wire format.
type ResolvedArtifact struct {
	Namespace      string             `json:"namespace,omitempty"`
	Name           string             `json:"name"`
	Version        string             `json:"version"`
	Registry       string             `json:"registry,omitempty"`
	RegistryType   string             `json:"registry_type,omitempty"`
	DownloadURL    string             `json:"download_url"`
	Artifact       string             `json:"artifact,omitempty"`
	SHA256         string             `json:"sha256,omitempty"`
	PackageType    string             `json:"package_type,omitempty"`
	CompatibleWith []skill.Platform   `json:"compatible_with,omitempty"`
	Metadata       map[string]string  `json:"metadata,omitempty"`
}

// VersionInfo is aliased from sklib/registryapi so version list responses share
// the canonical schema with the Skill Registry OpenAPI contract.
type VersionInfo = registryapi.SkillVersion

type SkillSearchResult struct {
	Namespace     string `json:"namespace,omitempty"`
	Name          string `json:"name"`
	Version       string `json:"version,omitempty"`
	LatestVersion string `json:"latest_version,omitempty"`
	Description   string `json:"description,omitempty"`
	Source        string `json:"source,omitempty"`
}

type SkillInfo struct {
	Namespace     string        `json:"namespace,omitempty"`
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	LatestVersion string        `json:"latest_version,omitempty"`
	Versions      []VersionInfo `json:"versions,omitempty"`
	Source        string        `json:"source,omitempty"`
}

type PublishRequest struct {
	ArtifactPath string              `json:"artifact_path"`
	Manifest     skill.SkillManifest `json:"manifest"`
	SHA256       string              `json:"sha256"`
	PackageType  string              `json:"package_type,omitempty"`
	Force        bool                `json:"force,omitempty"`
	DryRun       bool                `json:"dry_run,omitempty"`
}

type PublishResult struct {
	Namespace   string `json:"namespace,omitempty"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
	SHA256      string `json:"sha256"`
	Registry    string `json:"registry,omitempty"`
	Created     bool   `json:"created"`
}

type PublishMetadata struct {
	Name        string
	Version     string
	Description string
}

// RegistryCapabilities is aliased from sklib/registryapi so the canonical
// capabilities schema is shared across skpm, skcr, and SkillForge.
type RegistryCapabilities = registryapi.RegistryCapabilities

type Registry interface {
	Type() string
	Name() string
	Capabilities(ctx context.Context) (*RegistryCapabilities, error)
	Resolve(ctx context.Context, req ResolveRequest) (*ResolvedArtifact, error)
	Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error
}

type DiscoveryRegistry interface {
	Search(ctx context.Context, req SearchRequest) ([]SkillSearchResult, error)
	Info(ctx context.Context, ref SkillRef) (*SkillInfo, error)
	ListVersions(ctx context.Context, ref SkillRef) ([]VersionInfo, error)
}

type PublishingRegistry interface {
	Publish(ctx context.Context, req PublishRequest) (*PublishResult, error)
}

type GovernanceRegistry interface {
	Deprecate(ctx context.Context, ref SkillVersionRef, reason string) error
	Yank(ctx context.Context, ref SkillVersionRef, reason string) error
	Unyank(ctx context.Context, ref SkillVersionRef) error
}

func ParseSkillRef(raw, defaultNamespace string) SkillRef {
	raw = strings.TrimSpace(raw)
	if defaultNamespace == "" {
		defaultNamespace = "default"
	}
	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		if parts[0] != "" && parts[1] != "" {
			return SkillRef{Namespace: parts[0], Name: parts[1]}
		}
	}
	return SkillRef{Namespace: defaultNamespace, Name: raw}
}

func Unsupported(registryName, operation, hint string) error {
	if hint != "" {
		return fmt.Errorf("registry %q does not support %s; %s", registryName, operation, hint)
	}
	return fmt.Errorf("registry %q does not support %s", registryName, operation)
}
