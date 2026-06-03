package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"
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

type StructuredValidator struct{}

func NewValidator() Validator {
	return &StructuredValidator{}
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
		if len(skillYAML.CompatibleWith) == 0 {
			res.addWarning("skill.yaml", "compatible_with is empty — skill will not be installed to any platform path")
		}
	}

	if versionOK {
		v.validateChangelog(dir, version, res)
	}

	res.Valid = len(res.Errors) == 0
	return res, nil
}

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
	canonical := "v" + version
	if !semver.IsValid(canonical) {
		res.addError("VERSION", fmt.Sprintf("%q is not a valid semver string", version))
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
	if sy.Name == "" {
		res.addError("skill.yaml", "name is required")
	}
	if sy.Version == "" {
		res.addError("skill.yaml", "version is required")
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

func (res *ValidationResult) addError(field, msg string) {
	res.Errors = append(res.Errors, ValidationError{Field: field, Message: msg, Severity: SeverityError})
}

func (res *ValidationResult) addWarning(field, msg string) {
	res.Warnings = append(res.Warnings, ValidationError{Field: field, Message: msg, Severity: SeverityWarning})
}
