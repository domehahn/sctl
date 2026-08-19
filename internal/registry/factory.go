package registry

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
)

// BackendFactory constructs a named Registry adapter from its type-specific
// config.RegistryConfig. Each built-in backend registers its own factory via
// Register() in an init() in its own file (see github.go, gitlab.go,
// artifactory.go, local.go, generic_http.go) — adding a new registry type
// means implementing Registry and calling Register(), not editing a
// dispatch switch here. This is the pattern that let skillforge/generic-http
// support pull/resolve without also being wired into publish by hand; new
// backends now get both for free.
type BackendFactory func(name string, rc config.RegistryConfig) (Registry, error)

var factories = map[string]BackendFactory{}

// Register makes a registry type constructible via New() / a named
// cfg.Registries[...].type entry. Intended to be called from an init() in
// the backend's own file. Panics on a duplicate type name — that's a
// programming error (two backends claiming the same config type), not a
// runtime condition to handle gracefully.
func Register(typeName string, factory BackendFactory) {
	if _, exists := factories[typeName]; exists {
		panic(fmt.Sprintf("registry: factory for type %q already registered", typeName))
	}
	factories[typeName] = factory
}

// New returns a Registry for the given source identifier.
// source may be a named registry key from cfg.Registries, or one of:
// "github", "gitlab", "artifactory", "local", "skillforge",
// "generic-http", a file path, or a URL.
func New(source string, cfg *config.Config) (Registry, error) {
	if rc, ok := cfg.Registries[source]; ok {
		return fromRegistryConfig(source, rc)
	}

	switch {
	case source == "github":
		return NewGitHubRegistry("", "")
	case source == "gitlab":
		return NewGitLabRegistry("", "", "")
	case source == "artifactory":
		return nil, fmt.Errorf("factory: artifactory requires a named registry config with url and repo")
	case source == "local" || strings.HasPrefix(source, "/") || strings.HasPrefix(source, "./") || strings.HasPrefix(source, "../"):
		return NewLocalRegistry(source).WithName(source), nil
	case strings.HasPrefix(source, "https://github.com") || strings.HasPrefix(source, "http://github.com"):
		return NewGitHubRegistry("", "")
	case strings.HasPrefix(source, "https://gitlab.") || strings.HasPrefix(source, "http://gitlab."):
		return NewGitLabRegistry(source, "", "")
	default:
		return nil, fmt.Errorf("factory: unknown registry source %q — use a named registry from config or one of: github, gitlab, artifactory, local, skillforge, generic-http", source)
	}
}

func fromRegistryConfig(name string, rc config.RegistryConfig) (Registry, error) {
	factory, ok := factories[rc.Type]
	if !ok {
		return nil, fmt.Errorf("factory: unknown registry type %q for registry %q", rc.Type, name)
	}
	return factory(name, rc)
}

func authToken(rc config.RegistryConfig) string {
	if rc.Auth.Token != "" {
		return rc.Auth.Token
	}
	return rc.Token
}
