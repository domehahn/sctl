package installer

import (
	"fmt"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/skill"
)

// Platform maps a skill name to filesystem install paths.
type Platform interface {
	Name() skill.Platform
	InstallPaths(skillName string) []string
}

type claudeCodePlatform struct{}
type gitlabDuoPlatform struct{}
type githubCopilotPlatform struct{}
type codexPlatform struct{}
type piPlatform struct{}

func (claudeCodePlatform) Name() skill.Platform    { return skill.PlatformClaudeCode }
func (gitlabDuoPlatform) Name() skill.Platform     { return skill.PlatformGitLabDuo }
func (githubCopilotPlatform) Name() skill.Platform { return skill.PlatformGitHubCopilot }
func (codexPlatform) Name() skill.Platform         { return skill.PlatformCodex }
func (piPlatform) Name() skill.Platform              { return skill.PlatformPi }

func (claudeCodePlatform) InstallPaths(skillName string) []string {
	return []string{filepath.Join(".claude", "skills", skillName)}
}

func (gitlabDuoPlatform) InstallPaths(skillName string) []string {
	return []string{
		filepath.Join("skills", skillName),
		filepath.Join(".agents", "skills", skillName),
	}
}

func (githubCopilotPlatform) InstallPaths(skillName string) []string {
	return []string{filepath.Join(".github", "skills", skillName)}
}

func (codexPlatform) InstallPaths(skillName string) []string {
	return []string{filepath.Join(".agents", "skills", skillName)}
}

func (piPlatform) InstallPaths(skillName string) []string {
	return []string{
		filepath.Join(".pi", "skills", skillName),
		filepath.Join(".pi", "agent", "skills", skillName),
	}
}

var knownPlatforms = map[skill.Platform]Platform{
	skill.PlatformClaudeCode:    claudeCodePlatform{},
	skill.PlatformGitLabDuo:     gitlabDuoPlatform{},
	skill.PlatformGitHubCopilot: githubCopilotPlatform{},
	skill.PlatformCodex:         codexPlatform{},
	skill.PlatformPi:            piPlatform{},
}

// ResolvePaths returns the deduplicated set of filesystem paths for the given compatible_with list.
// If platforms contains "all", paths for every known platform are returned.
func ResolvePaths(skillName string, platforms []skill.Platform) ([]string, error) {
	effective := platforms
	for _, p := range platforms {
		if p == skill.PlatformAll {
			effective = []skill.Platform{
				skill.PlatformClaudeCode,
				skill.PlatformGitLabDuo,
				skill.PlatformGitHubCopilot,
				skill.PlatformCodex,
				skill.PlatformPi,
			}
			break
		}
	}

	seen := make(map[string]bool)
	var paths []string
	for _, p := range effective {
		impl, ok := knownPlatforms[p]
		if !ok {
			return nil, fmt.Errorf("unknown platform %q", p)
		}
		for _, path := range impl.InstallPaths(skillName) {
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	return paths, nil
}
