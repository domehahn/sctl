package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/domehahn/sklib/registryapi"
	"github.com/domehahn/skpm/v2/internal/httpclient"
	"github.com/domehahn/skpm/v2/internal/skill"
)

// sharedHTTPClient is the package-wide default HTTP client for registry
// backends that make ad-hoc requests (asset downloads, generic package
// uploads) outside of a purpose-built API client like go-github/go-gitlab.
// See internal/httpclient's package doc for why every such call goes
// through a bounded client instead of http.DefaultClient.
var sharedHTTPClient = httpclient.New()

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
	Namespace      string            `json:"namespace,omitempty"`
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Registry       string            `json:"registry,omitempty"`
	RegistryType   string            `json:"registry_type,omitempty"`
	DownloadURL    string            `json:"download_url"`
	Artifact       string            `json:"artifact,omitempty"`
	SHA256         string            `json:"sha256,omitempty"`
	PackageType    string            `json:"package_type,omitempty"`
	CompatibleWith []skill.Platform  `json:"compatible_with,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
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
	// TagFormat controls the git-release tag shape backends derive their
	// release/tag from: "prefixed" (<name>/v<ver>, the default) or "plain"
	// (v<ver>). Only consulted by tag-based backends (github, gitlab).
	TagFormat string `json:"tag_format,omitempty"`
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

// AttestationRequest carries a third-party attestation (e.g. a skil
// Attestation produced by `skil attest --output attestation.json`) to attach
// to a published skill version. Predicate is opaque to skpm — it's stored
// and returned as-is by the registry, letting skil's evidence schema evolve
// independently of skpm's registry client.
type AttestationRequest struct {
	// Type identifies the predicate's kind. What values are accepted is
	// registry-specific — there is no universal predicate-type URI
	// convention across backends here. SkillForge, for example, accepts a
	// fixed set: "signature", "scan", "provenance", "sbom" (skpm's CLI
	// defaults to "scan" for a skil evidence file — see
	// internal/cli/attestation.go).
	Type string `json:"type"`
	// Digest is the sha256 of the artifact the attestation is about
	// (normally the same digest skpm already computed when packaging).
	Digest    string          `json:"digest"`
	Predicate json.RawMessage `json:"predicate"`
}

// AttestationRecord is a stored attestation as returned by the registry.
type AttestationRecord struct {
	ID        int64           `json:"id,omitempty"`
	Type      string          `json:"type"`
	Digest    string          `json:"digest"`
	Predicate json.RawMessage `json:"predicate,omitempty"`
	CreatedBy string          `json:"created_by,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
}

// AttestationRegistry is an optional capability: a registry that can store
// and return third-party attestations (skil scan/eval evidence, provenance,
// signing records, ...) as first-class metadata attached to a published
// skill version, addressable independently of the artifact bytes.
type AttestationRegistry interface {
	Attest(ctx context.Context, ref SkillVersionRef, req AttestationRequest) (*AttestationRecord, error)
	ListAttestations(ctx context.Context, ref SkillVersionRef) ([]AttestationRecord, error)
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
