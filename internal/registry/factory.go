package registry

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/internal/config"
)

// New returns a Registry for the given source identifier.
// source may be a named registry key from cfg.Registries, or one of:
// "github", "gitlab", "artifactory", "local", a file path, or a URL.
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
		return NewLocalRegistry(source), nil
	case strings.HasPrefix(source, "https://github.com") || strings.HasPrefix(source, "http://github.com"):
		return NewGitHubRegistry("", "")
	case strings.HasPrefix(source, "https://gitlab.") || strings.HasPrefix(source, "http://gitlab."):
		return NewGitLabRegistry(source, "", "")
	default:
		return nil, fmt.Errorf("factory: unknown registry source %q — use a named registry from config or one of: github, gitlab, artifactory, local", source)
	}
}

func fromRegistryConfig(name string, rc config.RegistryConfig) (Registry, error) {
	switch rc.Type {
	case "github":
		repoSlug := rc.URL
		if repoSlug == "" {
			return nil, fmt.Errorf("factory: github registry %q requires url set to owner/repo", name)
		}
		repoSlug = strings.TrimPrefix(repoSlug, "https://github.com/")
		return NewGitHubRegistry(repoSlug, rc.Token)
	case "gitlab":
		if rc.Project == "" {
			return nil, fmt.Errorf("factory: gitlab registry %q requires project set to namespace/project (e.g. \"platform/agent-skills\")", name)
		}
		return NewGitLabRegistry(rc.URL, rc.Project, rc.Token)
	case "artifactory":
		parts := strings.SplitN(rc.URL, "#", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("factory: artifactory registry %q requires url in format <base-url>#<repo>", name)
		}
		return NewArtifactoryRegistry(parts[0], parts[1], rc.Token), nil
	case "local":
		return NewLocalRegistry(rc.URL), nil
	default:
		return nil, fmt.Errorf("factory: unknown registry type %q for registry %q", rc.Type, name)
	}
}
