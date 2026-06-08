package registry

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
)

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
	switch rc.Type {
	case "github":
		repoSlug := rc.Repo
		if repoSlug == "" {
			repoSlug = rc.URL
		}
		if repoSlug == "" {
			return nil, fmt.Errorf("factory: github registry %q requires repo set to owner/repo", name)
		}
		repoSlug = strings.TrimPrefix(repoSlug, "https://github.com/")
		reg, err := NewGitHubRegistry(repoSlug, authToken(rc))
		if err != nil {
			return nil, err
		}
		return reg.WithName(name), nil
	case "gitlab":
		if rc.Project == "" {
			return nil, fmt.Errorf("factory: gitlab registry %q requires project set to namespace/project (e.g. \"platform/agent-skills\")", name)
		}
		reg, err := NewGitLabRegistry(rc.URL, rc.Project, authToken(rc))
		if err != nil {
			return nil, err
		}
		return reg.WithName(name), nil
	case "artifactory":
		if rc.URL == "" || rc.Repo == "" {
			return nil, fmt.Errorf("factory: artifactory registry %q requires url and repo", name)
		}
		return NewArtifactoryRegistry(rc.URL, rc.Repo, authToken(rc)).WithName(name), nil
	case "local":
		path := rc.Path
		if path == "" {
			path = rc.URL
		}
		return NewLocalRegistry(path).WithName(name), nil
	case "generic-http":
		return NewGenericHTTPRegistry(name, rc), nil
	case "skillforge":
		return NewSkillForgeRegistry(name, rc), nil
	default:
		return nil, fmt.Errorf("factory: unknown registry type %q for registry %q", rc.Type, name)
	}
}

func authToken(rc config.RegistryConfig) string {
	if rc.Auth.Token != "" {
		return rc.Auth.Token
	}
	return rc.Token
}
