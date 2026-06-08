package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
)

type ValidationError struct {
	Field    string
	Message  string
	Severity ValidationSeverity
}

type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
}

type Validator interface {
	Validate(ctx context.Context, dir string) (*ValidationResult, error)
}

type ValidationOptions struct {
	Strict   bool
	Publish  bool
	Platform Platform
}

type StructuredValidator struct {
	options ValidationOptions
}

func NewValidator() Validator {
	return &StructuredValidator{}
}

func NewValidatorWithOptions(options ValidationOptions) Validator {
	return &StructuredValidator{options: options}
}

func (v *StructuredValidator) Validate(_ context.Context, dir string) (*ValidationResult, error) {
	res := &ValidationResult{Valid: true}

	skillMD := filepath.Join(dir, "SKILL.md")
	if info, err := os.Stat(skillMD); err != nil || info.Size() == 0 {
		res.addError("SKILL.md", "file is missing or empty")
	}

	version, versionOK := v.validateVersion(dir, res)

	skillYAML, yamlOK := v.validateSkillYAML(dir, res)
	if yamlOK && versionOK {
		if skillYAML.Version != version {
			res.addError("skill.yaml", fmt.Sprintf("version %q does not match VERSION file %q", skillYAML.Version, version))
		}
		for _, p := range skillYAML.CompatibleWith {
			if !KnownPlatforms[p] {
				res.addError("skill.yaml", fmt.Sprintf("unknown platform %q in compatible_with", p))
			}
		}
		if v.options.Platform != "" && !skillYAML.supportsPlatform(v.options.Platform) {
			res.addError("skill.yaml", fmt.Sprintf("skill is not compatible with platform %q", v.options.Platform))
		}
	}

	if versionOK {
		v.validateChangelog(dir, version, res)
	}
	v.validateHygiene(dir, res)
	if v.options.Strict || v.options.Publish {
		v.promoteWarnings(res)
	}

	res.Valid = len(res.Errors) == 0
	return res, nil
}

var stableVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

func (v *StructuredValidator) validateVersion(dir string, res *ValidationResult) (string, bool) {
	path := filepath.Join(dir, "VERSION")
	data, err := os.ReadFile(path)
	if err != nil {
		res.addError("VERSION", "file is missing")
		return "", false
	}
	version := strings.TrimSpace(string(data))
	if version == "" {
		res.addError("VERSION", "file is empty")
		return "", false
	}
	if !stableVersionPattern.MatchString(version) {
		res.addError("VERSION", fmt.Sprintf("%q is not a stable SemVer string (expected 1.2.3)", version))
		return version, false
	}
	return version, true
}

func (v *StructuredValidator) validateSkillYAML(dir string, res *ValidationResult) (*SkillYAML, bool) {
	path := filepath.Join(dir, "skill.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		res.addError("skill.yaml", "file is missing")
		return nil, false
	}
	var sy SkillYAML
	if err := yaml.Unmarshal(data, &sy); err != nil {
		res.addError("skill.yaml", fmt.Sprintf("parse error: %v", err))
		return nil, false
	}
	if strings.TrimSpace(sy.Name) == "" {
		res.addError("skill.yaml", "name is required")
	}
	if strings.TrimSpace(sy.Version) == "" {
		res.addError("skill.yaml", "version is required")
	}
	if strings.TrimSpace(sy.Description) == "" {
		res.addError("skill.yaml", "description is required")
	}
	if len(sy.CompatibleWith) == 0 {
		res.addError("skill.yaml", "compatible_with must include at least one platform")
	}
	return &sy, true
}

var changelogVersionPattern = regexp.MustCompile(`(?m)^##\s+v?(\S+)`)

func (v *StructuredValidator) validateChangelog(dir, version string, res *ValidationResult) {
	path := filepath.Join(dir, "CHANGELOG.md")
	data, err := os.ReadFile(path)
	if err != nil {
		res.addWarning("CHANGELOG.md", "file is missing — recommended for production skills")
		return
	}
	matches := changelogVersionPattern.FindAllStringSubmatch(string(data), -1)
	for _, m := range matches {
		if m[1] == version || m[1] == "v"+version {
			return
		}
	}
	res.addWarning("CHANGELOG.md", fmt.Sprintf("no entry found for version %q", version))
}

func (v *StructuredValidator) validateHygiene(dir string, res *ValidationResult) {
	forbidden := map[string]bool{".env": true, "id_rsa": true, "id_ed25519": true}
	secretPattern := regexp.MustCompile(`(?i)(api[_-]?key|secret|token|password)\s*[:=]\s*['"]?[A-Za-z0-9_\-]{16,}`)
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			if name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if forbidden[name] {
			res.addError(rel, "forbidden file is not allowed in a skill package")
		}
		info, err := d.Info()
		if err == nil && info.Size() > 5*1024*1024 {
			res.addWarning(rel, "large file exceeds 5 MiB")
		}
		data, err := os.ReadFile(path)
		if err == nil {
			if strings.Contains(string(data), "/Users/") || strings.Contains(string(data), "C:\\") {
				res.addWarning(rel, "contains an absolute local path")
			}
			if secretPattern.Match(data) {
				res.addError(rel, "possible secret detected")
			}
		}
		return nil
	})
}

func (v *StructuredValidator) promoteWarnings(res *ValidationResult) {
	for _, warning := range res.Warnings {
		res.Errors = append(res.Errors, ValidationError{
			Field:    warning.Field,
			Message:  warning.Message,
			Severity: SeverityError,
		})
	}
	res.Warnings = nil
}

func (res *ValidationResult) addError(field, msg string) {
	res.Errors = append(res.Errors, ValidationError{Field: field, Message: msg, Severity: SeverityError})
}

func (res *ValidationResult) addWarning(field, msg string) {
	res.Warnings = append(res.Warnings, ValidationError{Field: field, Message: msg, Severity: SeverityWarning})
}

func (sy *SkillYAML) supportsPlatform(platform Platform) bool {
	for _, p := range sy.CompatibleWith {
		if p == platform || p == PlatformAll {
			return true
		}
	}
	return false
}
