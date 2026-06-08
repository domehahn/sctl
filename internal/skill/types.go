package skill

import (
	"github.com/domehahn/sklib/spec"
	"github.com/domehahn/sklib/validate"
)

// Platform is aliased from sklib/spec to ensure canonical platform identifiers
// are shared across all tools in the ecosystem.
type Platform = spec.Platform

// Platform constants — direct aliases from sklib/spec.
const (
	PlatformClaudeCode    Platform = spec.PlatformClaudeCode
	PlatformGitLabDuo     Platform = spec.PlatformGitLabDuo
	PlatformGitHubCopilot Platform = spec.PlatformGitHubCopilot
	PlatformCodex         Platform = spec.PlatformCodex
	PlatformCursor        Platform = spec.PlatformCursor
	PlatformWindsurf      Platform = spec.PlatformWindsurf
	PlatformOpenHands     Platform = spec.PlatformOpenHands
	PlatformOpenCode      Platform = spec.PlatformOpenCode
	PlatformOllama        Platform = spec.PlatformOllama
	PlatformGeneric       Platform = spec.PlatformGeneric
	PlatformAll           Platform = spec.PlatformAll
)

// KnownPlatforms is the canonical set of valid platform identifiers.
// Kept for backward compatibility; use spec.IsKnownPlatform for new code.
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

// NormalizePlatform normalizes a platform identifier using sklib/spec rules.
// Returns p unchanged if it cannot be resolved (logs-safe; callers validate separately).
func NormalizePlatform(p Platform) Platform {
	canonical, err := spec.NormalizePlatform(string(p))
	if err != nil {
		return p
	}
	return canonical
}

// ValidationSeverity is aliased from sklib/validate.
type ValidationSeverity = validate.Severity

// Severity constants — aliases from sklib/validate.
const (
	SeverityError   ValidationSeverity = validate.SeverityError
	SeverityWarning ValidationSeverity = validate.SeverityWarning
	SeverityInfo    ValidationSeverity = validate.SeverityInfo
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
// It retains skpm-specific fields (Profile, Path, and separate slices per severity)
// beyond what sklib/validate.Result provides.
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

// supportsPlatform reports whether the skill declares support for the given platform.
func (s *SkillYAML) supportsPlatform(p Platform) bool {
	return spec.SupportsPlatform(s.CompatibleWith, p)
}

// SkillManifest is aliased from sklib/spec so the canonical manifest.json schema
// is defined in one place and shared across skpm, skcr, and SkillForge.
type SkillManifest = spec.PackageManifest
