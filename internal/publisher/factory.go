package publisher

import (
	"fmt"
	"strings"

	"github.com/domehahn/sctl/internal/config"
)

// New returns a Publisher for the given registry source and tag format.
// tagFormat is "prefixed" (<name>/v<ver>) or "plain" (v<ver>).
func New(source, tagFormat string, cfg *config.Config) (Publisher, error) {
	rc, ok := cfg.Registries[source]
	if !ok {
		return nil, fmt.Errorf("publisher: registry %q not found in config", source)
	}

	token := rc.Token
	if token == "" {
		if env := registryTokenFromEnv(source, cfg); env != "" {
			token = env
		}
	}

	switch rc.Type {
	case "github":
		repoSlug := strings.TrimPrefix(rc.URL, "https://github.com/")
		repoSlug = strings.TrimPrefix(repoSlug, "http://github.com/")
		return NewGitHubPublisher(repoSlug, token, tagFormat)

	case "gitlab":
		if rc.Project == "" {
			return nil, fmt.Errorf("publisher: gitlab registry %q requires project field", source)
		}
		return NewGitLabPublisher(rc.URL, rc.Project, token, tagFormat)

	case "artifactory":
		parts := strings.SplitN(rc.URL, "#", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("publisher: artifactory url must be <base-url>#<repo>")
		}
		return NewArtifactoryPublisher(parts[0], parts[1], token), nil

	case "local":
		return NewLocalPublisher(rc.URL), nil

	default:
		return nil, fmt.Errorf("publisher: unsupported registry type %q", rc.Type)
	}
}

func registryTokenFromEnv(source string, cfg *config.Config) string {
	if source == cfg.DefaultRegistry {
		// SKPM_REGISTRY_TOKEN is applied by config.Load() already, but in case
		// of direct factory calls check applyEnv wasn't called yet.
		return ""
	}
	return ""
}
