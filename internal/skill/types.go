package skill

type Platform string

const (
	PlatformClaudeCode    Platform = "claude-code"
	PlatformGitLabDuo     Platform = "gitlab-duo"
	PlatformGitHubCopilot Platform = "github-copilot"
	PlatformCodex         Platform = "codex"
	PlatformAll           Platform = "all"
)

var KnownPlatforms = map[Platform]bool{
	PlatformClaudeCode:    true,
	PlatformGitLabDuo:     true,
	PlatformGitHubCopilot: true,
	PlatformCodex:         true,
	PlatformAll:           true,
}

type SkillYAML struct {
	Name           string     `yaml:"name"`
	Version        string     `yaml:"version"`
	Description    string     `yaml:"description"`
	Owners         []string   `yaml:"owners"`
	CompatibleWith []Platform `yaml:"compatible_with"`
}

type SkillManifest struct {
	Name           string     `json:"name"`
	Version        string     `json:"version"`
	SHA256         string     `json:"sha256"`
	CreatedAt      string     `json:"created_at"`
	SourceCommit   string     `json:"source_commit,omitempty"`
	CompatibleWith []Platform `json:"compatible_with"`
}
