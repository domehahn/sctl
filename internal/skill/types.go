package skill

// Platform identifies an AI coding platform.
type Platform string

const (
	PlatformClaudeCode    Platform = "claude-code"
	PlatformGitLabDuo     Platform = "gitlab-duo"
	PlatformGitHubCopilot Platform = "github-copilot"
	PlatformCodex         Platform = "codex"
	PlatformCursor        Platform = "cursor"
	PlatformWindsurf      Platform = "windsurf"
	PlatformOpenHands     Platform = "openhands"
	PlatformOpenCode      Platform = "opencode"
	PlatformOllama        Platform = "ollama"
	PlatformGeneric       Platform = "generic"
	PlatformAll           Platform = "all"
)

// KnownPlatforms is the canonical set of valid platform identifiers.
var KnownPlatforms = map[Platform]bool{
	PlatformClaudeCode:    true,
	PlatformGitLabDuo:     true,
	PlatformGitHubCopilot: true,
	PlatformCodex:         true,
	PlatformCursor:        true,
	PlatformWindsurf:      true,
	PlatformOpenHands:     true,
	PlatformOpenCode:      true,
	PlatformOllama:        true,
	PlatformGeneric:       true,
	PlatformAll:           true,
}

// PlatformAliases maps non-canonical platform names to their canonical form.
var PlatformAliases = map[Platform]Platform{
	"gitlab":  PlatformGitLabDuo,
	"duo":     PlatformGitLabDuo,
	"github":  PlatformGitHubCopilot,
	"copilot": PlatformGitHubCopilot,
	"claude":  PlatformClaudeCode,
}

// NormalizePlatform returns the canonical name for p.
// Returns p unchanged if it is already canonical or unknown.
func NormalizePlatform(p Platform) Platform {
	if KnownPlatforms[p] {
		return p
	}
	if canonical, ok := PlatformAliases[p]; ok {
		return canonical
	}
	return p
}

// ValidationSeverity classifies a validation finding.
type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
	SeverityInfo    ValidationSeverity = "info"
)

// ValidationFinding is a single validation issue with its context.
type ValidationFinding struct {
	Field    string             `json:"field,omitempty"`
	Path     string             `json:"path,omitempty"`
	Message  string             `json:"message"`
	Severity ValidationSeverity `json:"severity"`
	Code     string             `json:"code,omitempty"`
}

// ValidationError is an alias for ValidationFinding kept for compatibility.
type ValidationError = ValidationFinding

// ValidationResult holds the outcome of a validation run.
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Profile  string              `json:"profile"`
	Path     string              `json:"path"`
	Errors   []ValidationFinding `json:"errors,omitempty"`
	Warnings []ValidationFinding `json:"warnings,omitempty"`
	Infos    []ValidationFinding `json:"infos,omitempty"`
}

// ValidationOptions configures a validation run.
type ValidationOptions struct {
	Strict              bool
	Publish             bool
	Platform            Platform
	MaxFileSizeBytes    int64
	MaxFileCount        int
	AllowPrerelease     bool
	AllowMissingTests   bool
	AllowMissingLicense bool
}

// SkillYAML is the parsed representation of a skill.yaml file.
// Field order defines canonical YAML output order for skpm format.
type SkillYAML struct {
	Name           string            `yaml:"name"`
	Version        string            `yaml:"version"`
	Description    string            `yaml:"description"`
	Namespace      string            `yaml:"namespace,omitempty"`
	Owners         []string          `yaml:"owners,omitempty"`
	License        string            `yaml:"license,omitempty"`
	Entrypoint     string            `yaml:"entrypoint,omitempty"`
	Tags           []string          `yaml:"tags,omitempty"`
	CompatibleWith []Platform        `yaml:"compatible_with"`
	Security       *SkillSecurity    `yaml:"security,omitempty"`
	Metadata       map[string]string `yaml:"metadata,omitempty"`
}

// SkillSecurity declares the runtime requirements of a skill.
type SkillSecurity struct {
	RequiresNetwork bool `yaml:"requires_network"`
	RequiresSecrets bool `yaml:"requires_secrets"`
	WritesFiles     bool `yaml:"writes_files"`
	RunsCommands    bool `yaml:"runs_commands"`
}

// SkillManifest is embedded in the packaged ZIP as manifest.json.
type SkillManifest struct {
	Name           string     `json:"name"`
	Version        string     `json:"version"`
	SHA256         string     `json:"sha256"`
	CreatedAt      string     `json:"created_at"`
	SourceCommit   string     `json:"source_commit,omitempty"`
	CompatibleWith []Platform `json:"compatible_with"`
}
