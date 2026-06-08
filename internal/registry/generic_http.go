package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
)

type GenericHTTPRegistry struct {
	name       string
	regType    string
	baseURL    string
	namespace  string
	auth       config.AuthConfig
	headers    map[string]string
	endpoints  map[string]string
	staticCaps map[string]bool
	client     *http.Client
}

func NewGenericHTTPRegistry(name string, rc config.RegistryConfig) *GenericHTTPRegistry {
	return &GenericHTTPRegistry{
		name:       name,
		regType:    rc.Type,
		baseURL:    strings.TrimRight(rc.URL, "/"),
		namespace:  rc.Namespace,
		auth:       rc.Auth,
		headers:    rc.Headers,
		endpoints:  rc.Endpoints,
		staticCaps: rc.Capabilities,
		client:     &http.Client{},
	}
}

func NewSkillForgeRegistry(name string, rc config.RegistryConfig) *GenericHTTPRegistry {
	if rc.Endpoints == nil {
		rc.Endpoints = map[string]string{}
	}
	defaults := map[string]string{
		"capabilities": "/api/v1/capabilities",
		"search":       "/api/v1/skills?q={query}",
		"info":         "/api/v1/skills/{namespace}/{name}",
		"versions":     "/api/v1/skills/{namespace}/{name}/versions",
		"resolve":      "/api/v1/skills/{namespace}/{name}/resolve?constraint={constraint}",
		"download":     "/api/v1/skills/{namespace}/{name}/versions/{version}/download",
		"publish":      "/api/v1/skills/{namespace}/{name}/versions/{version}",
		"deprecate":    "/api/v1/skills/{namespace}/{name}/versions/{version}/deprecate",
		"yank":         "/api/v1/skills/{namespace}/{name}/versions/{version}/yank",
		"unyank":       "/api/v1/skills/{namespace}/{name}/versions/{version}/unyank",
	}
	for k, v := range defaults {
		if rc.Endpoints[k] == "" {
			rc.Endpoints[k] = v
		}
	}
	rc.Type = "skillforge"
	return NewGenericHTTPRegistry(name, rc)
}

func (r *GenericHTTPRegistry) Type() string { return r.regType }

func (r *GenericHTTPRegistry) Name() string { return r.name }

func (r *GenericHTTPRegistry) Capabilities(ctx context.Context) (*RegistryCapabilities, error) {
	if endpoint := r.endpoints["capabilities"]; endpoint != "" {
		req, err := r.newRequest(ctx, http.MethodGet, r.endpointURL(endpoint, nil), nil)
		if err != nil {
			return nil, err
		}
		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%s capabilities: %w", r.name, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var caps RegistryCapabilities
			if err := json.NewDecoder(resp.Body).Decode(&caps); err != nil {
				return nil, fmt.Errorf("%s capabilities decode: %w", r.name, err)
			}
			return &caps, nil
		}
	}
	return &RegistryCapabilities{
		Resolve:           r.supports("resolve"),
		Download:          r.supports("download"),
		Search:            r.supports("search"),
		Info:              r.supports("info"),
		ListVersions:      r.supports("versions"),
		Publish:           r.supports("publish"),
		Deprecate:         r.supports("deprecate"),
		Yank:              r.supports("yank"),
		Unyank:            r.supports("unyank"),
		SemVerConstraints: r.supports("semver_constraints"),
		Checksums:         r.supports("checksums"),
	}, nil
}

func (r *GenericHTTPRegistry) Resolve(ctx context.Context, req ResolveRequest) (*ResolvedArtifact, error) {
	endpoint := r.endpoints["resolve"]
	if endpoint == "" {
		return nil, Unsupported(r.name, "resolve", "configure endpoints.resolve or use another registry")
	}
	values := r.values(req.Ref, req.Constraint, "")
	httpReq, err := r.newRequest(ctx, http.MethodGet, r.endpointURL(endpoint, values), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s resolve: %w", r.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s resolve: HTTP %d", r.name, resp.StatusCode)
	}
	var artifact ResolvedArtifact
	if err := json.NewDecoder(resp.Body).Decode(&artifact); err != nil {
		return nil, fmt.Errorf("%s resolve decode: %w", r.name, err)
	}
	r.fillArtifact(&artifact, req.Ref)
	return &artifact, nil
}

func (r *GenericHTTPRegistry) Download(ctx context.Context, artifact *ResolvedArtifact, dest io.Writer) error {
	downloadURL := artifact.DownloadURL
	if downloadURL == "" {
		endpoint := r.endpoints["download"]
		if endpoint == "" {
			return Unsupported(r.name, "download", "configure endpoints.download or use another registry")
		}
		downloadURL = r.endpointURL(endpoint, r.values(SkillRef{Namespace: artifact.Namespace, Name: artifact.Name}, artifact.Version, artifact.Artifact))
	}
	req, err := r.newRequest(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s download: %w", r.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s download: HTTP %d", r.name, resp.StatusCode)
	}
	_, err = io.Copy(dest, resp.Body)
	return err
}

func (r *GenericHTTPRegistry) Search(ctx context.Context, req SearchRequest) ([]SkillSearchResult, error) {
	endpoint := r.endpoints["search"]
	if endpoint == "" {
		return nil, Unsupported(r.name, "search", "configure endpoints.search or use another registry")
	}
	httpReq, err := r.newRequest(ctx, http.MethodGet, r.endpointURL(endpoint, map[string]string{"query": req.Query, "namespace": req.Namespace}), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s search: %w", r.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s search: HTTP %d", r.name, resp.StatusCode)
	}
	var payload struct {
		Skills []SkillSearchResult `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%s search decode: %w", r.name, err)
	}
	return payload.Skills, nil
}

func (r *GenericHTTPRegistry) Info(ctx context.Context, ref SkillRef) (*SkillInfo, error) {
	endpoint := r.endpoints["info"]
	if endpoint == "" {
		return nil, Unsupported(r.name, "info", "configure endpoints.info or use another registry")
	}
	httpReq, err := r.newRequest(ctx, http.MethodGet, r.endpointURL(endpoint, r.values(ref, "", "")), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s info: %w", r.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s info: HTTP %d", r.name, resp.StatusCode)
	}
	var info SkillInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("%s info decode: %w", r.name, err)
	}
	return &info, nil
}

func (r *GenericHTTPRegistry) ListVersions(ctx context.Context, ref SkillRef) ([]VersionInfo, error) {
	endpoint := r.endpoints["versions"]
	if endpoint == "" {
		return nil, Unsupported(r.name, "versions", "configure endpoints.versions or use another registry")
	}
	httpReq, err := r.newRequest(ctx, http.MethodGet, r.endpointURL(endpoint, r.values(ref, "", "")), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s versions: %w", r.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s versions: HTTP %d", r.name, resp.StatusCode)
	}
	var payload struct {
		Versions []VersionInfo `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%s versions decode: %w", r.name, err)
	}
	return payload.Versions, nil
}

func (r *GenericHTTPRegistry) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	endpoint := r.endpoints["publish"]
	if endpoint == "" {
		return nil, Unsupported(r.name, "publish", "configure endpoints.publish or use another registry")
	}
	ref := SkillRef{Namespace: r.namespace, Name: req.Manifest.Name}
	values := r.values(ref, req.Manifest.Version, "")
	body, err := os.ReadFile(req.ArtifactPath)
	if err != nil {
		return nil, fmt.Errorf("read artifact: %w", err)
	}
	httpReq, err := r.newRequest(ctx, http.MethodPut, r.endpointURL(endpoint, values), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/zip")
	httpReq.Header.Set("X-SKPM-SHA256", req.SHA256)
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s publish: %w", r.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s publish: HTTP %d", r.name, resp.StatusCode)
	}
	var result PublishResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result = PublishResult{Name: req.Manifest.Name, Version: req.Manifest.Version, SHA256: req.SHA256, Registry: r.name, Created: resp.StatusCode == http.StatusCreated}
	}
	return &result, nil
}

func (r *GenericHTTPRegistry) Deprecate(ctx context.Context, ref SkillVersionRef, reason string) error {
	return r.postGovernance(ctx, "deprecate", ref, reason)
}

func (r *GenericHTTPRegistry) Yank(ctx context.Context, ref SkillVersionRef, reason string) error {
	return r.postGovernance(ctx, "yank", ref, reason)
}

func (r *GenericHTTPRegistry) Unyank(ctx context.Context, ref SkillVersionRef) error {
	return r.postGovernance(ctx, "unyank", ref, "")
}

func (r *GenericHTTPRegistry) postGovernance(ctx context.Context, op string, ref SkillVersionRef, reason string) error {
	endpoint := r.endpoints[op]
	if endpoint == "" {
		return Unsupported(r.name, op, "configure endpoints."+op+" or use another registry")
	}
	payload, _ := json.Marshal(map[string]string{"reason": reason})
	httpReq, err := r.newRequest(ctx, http.MethodPost, r.endpointURL(endpoint, r.values(SkillRef{Namespace: ref.Namespace, Name: ref.Name}, ref.Version, "")), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("%s %s: %w", r.name, op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: HTTP %d", r.name, op, resp.StatusCode)
	}
	return nil
}

func (r *GenericHTTPRegistry) newRequest(ctx context.Context, method, rawURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	switch r.auth.Type {
	case "bearer":
		if r.auth.Token != "" {
			req.Header.Set("Authorization", "Bearer "+r.auth.Token)
		}
	case "basic":
		req.SetBasicAuth(r.auth.Username, r.auth.Password)
	case "header":
		if r.auth.HeaderName != "" && r.auth.Token != "" {
			req.Header.Set(r.auth.HeaderName, r.auth.Token)
		}
	}
	return req, nil
}

func (r *GenericHTTPRegistry) endpointURL(template string, values map[string]string) string {
	out := template
	for k, v := range values {
		out = strings.ReplaceAll(out, "{"+k+"}", url.QueryEscape(v))
	}
	if strings.HasPrefix(out, "http://") || strings.HasPrefix(out, "https://") {
		return out
	}
	return r.baseURL + "/" + strings.TrimLeft(out, "/")
}

func (r *GenericHTTPRegistry) values(ref SkillRef, version, artifact string) map[string]string {
	namespace := ref.Namespace
	if namespace == "" {
		namespace = r.namespace
	}
	if namespace == "" {
		namespace = "default"
	}
	return map[string]string{
		"namespace":  namespace,
		"name":       ref.Name,
		"version":    version,
		"constraint": version,
		"artifact":   artifact,
	}
}

func (r *GenericHTTPRegistry) supports(op string) bool {
	if r.staticCaps != nil {
		if supported, ok := r.staticCaps[op]; ok {
			return supported
		}
	}
	return r.endpoints[op] != ""
}

func (r *GenericHTTPRegistry) fillArtifact(artifact *ResolvedArtifact, ref SkillRef) {
	if artifact.Namespace == "" {
		artifact.Namespace = ref.Namespace
	}
	if artifact.Namespace == "" {
		artifact.Namespace = r.namespace
	}
	if artifact.Name == "" {
		artifact.Name = ref.Name
	}
	if artifact.Registry == "" {
		artifact.Registry = r.name
	}
	if artifact.RegistryType == "" {
		artifact.RegistryType = r.regType
	}
	if artifact.PackageType == "" {
		artifact.PackageType = "zip"
	}
}
