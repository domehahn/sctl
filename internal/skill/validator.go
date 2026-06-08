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

type Validator interface {
	Validate(ctx context.Context, dir string) (*ValidationResult, error)
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

func (v *StructuredValidator) profile() string {
	if v.options.Publish {
		return "publish"
	}
	if v.options.Strict {
		return "strict"
	}
	return "default"
}

func (v *StructuredValidator) strictOrPublish() bool {
	return v.options.Strict || v.options.Publish
}

// strictWarn adds an error in strict/publish mode, a warning in default mode.
func (v *StructuredValidator) strictWarn(res *ValidationResult, field, msg, code string) {
	if v.strictOrPublish() {
		res.addError(field, msg, code)
	} else {
		res.addWarning(field, msg, code)
	}
}

func (v *StructuredValidator) Validate(_ context.Context, dir string) (*ValidationResult, error) {
	res := &ValidationResult{
		Valid:   true,
		Profile: v.profile(),
		Path:    dir,
	}

	info, err := os.Stat(filepath.Join(dir, "SKILL.md"))
	if err != nil {
		res.addError("SKILL.md", "file is missing", "missing_skill_md")
	} else if info.Size() == 0 {
		res.addError("SKILL.md", "file is empty", "empty_skill_md")
	}

	version, versionOK := v.validateVersion(dir, res)
	skillYAML, yamlOK := v.validateSkillYAML(dir, res)

	if yamlOK && versionOK {
		if skillYAML.Version != version {
			res.addError("skill.yaml", fmt.Sprintf("version %q does not match VERSION file %q", skillYAML.Version, version), "version_mismatch")
		}
		v.validatePlatforms(skillYAML, res)
		if v.options.Platform != "" && !skillYAML.supportsPlatform(v.options.Platform) {
			res.addError("skill.yaml", fmt.Sprintf("skill is not compatible with platform %q", v.options.Platform), "platform_not_compatible")
		}
		if v.strictOrPublish() {
			v.validateDuplicates(skillYAML, res)
			v.validateUnknownFields(dir, res)
		}
		if v.options.Publish {
			v.validatePublishRequirements(dir, skillYAML, res)
		}
	}

	if versionOK {
		v.validateChangelog(dir, version, res)
	}
	v.validateOptionalFiles(dir, res)
	v.validateHygiene(dir, res)

	if v.options.Publish {
		v.promoteWarnings(res)
	}

	res.Valid = len(res.Errors) == 0
	return res, nil
}

var stableVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
var prereleaseVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-[0-9A-Za-z.-]+(\+[0-9A-Za-z.-]+)?$`)

func (v *StructuredValidator) validateVersion(dir string, res *ValidationResult) (string, bool) {
	data, err := os.ReadFile(filepath.Join(dir, "VERSION"))
	if err != nil {
		res.addError("VERSION", "file is missing", "missing_version")
		return "", false
	}
	version := strings.TrimSpace(string(data))
	if version == "" {
		res.addError("VERSION", "file is empty", "empty_version")
		return "", false
	}
	version = strings.TrimPrefix(version, "v")
	if stableVersionPattern.MatchString(version) {
		return version, true
	}
	if prereleaseVersionPattern.MatchString(version) {
		if v.options.AllowPrerelease {
			return version, true
		}
		res.addError("VERSION", fmt.Sprintf("%q is a prerelease version; use MAJOR.MINOR.PATCH for stable releases (or --allow-prerelease)", version), "prerelease_version")
		return version, false
	}
	res.addError("VERSION", fmt.Sprintf("%q is not valid stable SemVer; expected MAJOR.MINOR.PATCH", version), "invalid_semver")
	return version, false
}

var knownSkillYAMLFields = map[string]bool{
	"name": true, "version": true, "description": true,
	"namespace": true, "owners": true, "license": true,
	"entrypoint": true, "tags": true, "compatible_with": true,
	"security": true, "metadata": true,
}

func (v *StructuredValidator) validateSkillYAML(dir string, res *ValidationResult) (*SkillYAML, bool) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		res.addError("skill.yaml", "file is missing", "missing_skill_yaml")
		return nil, false
	}
	var sy SkillYAML
	if err := yaml.Unmarshal(data, &sy); err != nil {
		res.addError("skill.yaml", fmt.Sprintf("parse error: %v", err), "yaml_parse_error")
		return nil, false
	}
	if strings.TrimSpace(sy.Name) == "" {
		res.addError("skill.yaml", "name is required", "missing_name")
	}
	if strings.TrimSpace(sy.Version) == "" {
		res.addError("skill.yaml", "version is required", "missing_version_field")
	}
	if strings.TrimSpace(sy.Description) == "" {
		res.addError("skill.yaml", "description is required", "missing_description")
	}
	if len(sy.CompatibleWith) == 0 {
		res.addError("skill.yaml", "compatible_with must include at least one platform", "missing_compatible_with")
	}
	return &sy, true
}

func (v *StructuredValidator) validatePlatforms(sy *SkillYAML, res *ValidationResult) {
	for _, p := range sy.CompatibleWith {
		if !KnownPlatforms[NormalizePlatform(p)] {
			res.addError("skill.yaml", fmt.Sprintf("unknown platform %q in compatible_with", p), "unknown_platform")
		}
	}
}

func (v *StructuredValidator) validateDuplicates(sy *SkillYAML, res *ValidationResult) {
	seen := map[Platform]bool{}
	for _, p := range sy.CompatibleWith {
		canon := NormalizePlatform(p)
		if seen[canon] {
			res.addError("skill.yaml", fmt.Sprintf("duplicate platform %q in compatible_with", canon), "duplicate_platform")
		}
		seen[canon] = true
	}
	seenTags := map[string]bool{}
	for _, t := range sy.Tags {
		if seenTags[t] {
			res.addError("skill.yaml", fmt.Sprintf("duplicate tag %q in tags", t), "duplicate_tag")
		}
		seenTags[t] = true
	}
}

func (v *StructuredValidator) validateUnknownFields(dir string, res *ValidationResult) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		return
	}
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return
	}
	for key := range raw {
		if !knownSkillYAMLFields[key] {
			res.addError("skill.yaml", fmt.Sprintf("unknown field %q (use the metadata: section for custom fields)", key), "unknown_field")
		}
	}
}

var changelogVersionPattern = regexp.MustCompile(`(?m)^##\s+v?(\S+)`)

func (v *StructuredValidator) validateChangelog(dir, version string, res *ValidationResult) {
	data, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		v.strictWarn(res, "CHANGELOG.md", "file is missing", "missing_changelog")
		return
	}
	for _, m := range changelogVersionPattern.FindAllStringSubmatch(string(data), -1) {
		if m[1] == version || m[1] == "v"+version {
			return
		}
	}
	v.strictWarn(res, "CHANGELOG.md", fmt.Sprintf("no entry found for version %q", version), "missing_changelog_entry")
}

func (v *StructuredValidator) validateOptionalFiles(dir string, res *ValidationResult) {
	if _, err := os.Stat(filepath.Join(dir, "README.md")); os.IsNotExist(err) {
		v.strictWarn(res, "README.md", "file is missing", "missing_readme")
	}
	if _, err := os.Stat(filepath.Join(dir, "tests")); os.IsNotExist(err) {
		switch {
		case v.strictOrPublish() && !v.options.AllowMissingTests:
			res.addError("tests/", "tests directory is missing", "missing_tests")
		case v.strictOrPublish() && v.options.AllowMissingTests:
			res.addInfo("tests/", "tests directory is missing (allowed by --allow-missing-tests)", "missing_tests")
		default:
			res.addWarning("tests/", "tests directory is missing", "missing_tests")
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "LICENSE")); os.IsNotExist(err) {
		switch {
		case v.strictOrPublish() && !v.options.AllowMissingLicense:
			res.addError("LICENSE", "file is missing", "missing_license")
		case v.strictOrPublish() && v.options.AllowMissingLicense:
			res.addInfo("LICENSE", "file is missing (allowed by --allow-missing-license)", "missing_license")
		default:
			res.addWarning("LICENSE", "file is missing — recommended for production skills", "missing_license")
		}
	}
}

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$|^[a-z][a-z0-9]?$`)

func (v *StructuredValidator) validatePublishRequirements(dir string, sy *SkillYAML, res *ValidationResult) {
	entrypoint := sy.Entrypoint
	if entrypoint == "" {
		entrypoint = "SKILL.md"
	}
	if _, err := os.Stat(filepath.Join(dir, entrypoint)); os.IsNotExist(err) {
		res.addError("skill.yaml", fmt.Sprintf("entrypoint %q does not exist", entrypoint), "missing_entrypoint")
	}
	if sy.Name != "" && !namePattern.MatchString(sy.Name) {
		res.addError("skill.yaml", fmt.Sprintf("name %q must be lowercase alphanumeric with hyphens, no leading/trailing hyphens", sy.Name), "invalid_name")
	}
}

var (
	secretPattern     = regexp.MustCompile(`(?i)(api[_-]?key|secret|token|password)\s*[:=]\s*['"]?[A-Za-z0-9_\-]{16,}`)
	privateKeyPattern = regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)
	localPathPattern  = regexp.MustCompile(`(/Users/[^\s"']+|/home/[^\s"']+|/var/folders/[^\s"']+|C:\\[^\s"']+)`)
)

var forbiddenFileNames = map[string]bool{
	".env":        true,
	"id_rsa":      true,
	"id_ed25519":  true,
}

var forbiddenFilePatterns = []string{"*.pem", "*.key"}

var cloudCredentialNames = map[string]bool{
	"credentials":          true,
	"credentials.json":     true,
	"service_account.json": true,
	"gcloud.json":          true,
}

var generatedArtifactPatterns = []string{"manifest.json", "checksums.txt", "*.zip", "*.tgz"}

var buildDirNames = map[string]bool{
	"node_modules": true,
	".venv":        true,
	"dist":         true,
	"target":       true,
	".cache":       true,
}

func isForbiddenFile(name string) bool {
	if forbiddenFileNames[name] {
		return true
	}
	for _, pat := range forbiddenFilePatterns {
		if matched, _ := filepath.Match(pat, name); matched {
			return true
		}
	}
	return false
}

func isGeneratedArtifact(name string) bool {
	for _, pat := range generatedArtifactPatterns {
		if matched, _ := filepath.Match(pat, name); matched {
			return true
		}
	}
	return false
}

func (v *StructuredValidator) validateHygiene(dir string, res *ValidationResult) {
	maxSize := v.options.MaxFileSizeBytes
	if maxSize == 0 {
		maxSize = 5 * 1024 * 1024
	}

	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)

		if d.IsDir() {
			if name == ".git" {
				return filepath.SkipDir
			}
			if v.strictOrPublish() && buildDirNames[name] {
				res.addError(rel, fmt.Sprintf("build/cache directory %q should not be present in a skill package", name), "build_dir_present")
				return filepath.SkipDir
			}
			return nil
		}

		if isForbiddenFile(name) {
			res.addError(rel, "forbidden file is not allowed in a skill package", "forbidden_file")
			return nil
		}
		if cloudCredentialNames[name] {
			res.addError(rel, "credential file is not allowed in a skill package", "forbidden_file")
			return nil
		}
		if v.strictOrPublish() && isGeneratedArtifact(name) {
			res.addError(rel, fmt.Sprintf("generated artifact %q should not be checked in; add it to .gitignore", name), "generated_artifact")
			return nil
		}

		info, statErr := d.Info()
		if statErr == nil && info.Size() > maxSize {
			res.addWarning(rel, fmt.Sprintf("large file exceeds %d MiB", maxSize/(1024*1024)), "large_file")
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}

		if privateKeyPattern.Match(data) {
			res.addError(rel, "private key block detected", "possible_secret")
			return nil
		}
		if name == ".npmrc" && secretPattern.Match(data) {
			res.addError(rel, "forbidden file contains credentials", "forbidden_file")
			return nil
		}
		if secretPattern.Match(data) {
			res.addError(rel, "possible secret detected", "possible_secret")
		}
		if localPathPattern.Match(data) {
			v.strictWarn(res, rel, "contains an absolute local path", "absolute_path")
		}
		return nil
	})
}

func (v *StructuredValidator) promoteWarnings(res *ValidationResult) {
	for _, w := range res.Warnings {
		res.Errors = append(res.Errors, ValidationFinding{
			Field:    w.Field,
			Path:     w.Path,
			Message:  w.Message,
			Severity: SeverityError,
			Code:     w.Code,
		})
	}
	res.Warnings = nil
}

func (res *ValidationResult) addError(field, msg, code string) {
	res.Errors = append(res.Errors, ValidationFinding{Field: field, Message: msg, Severity: SeverityError, Code: code})
}

func (res *ValidationResult) addWarning(field, msg, code string) {
	res.Warnings = append(res.Warnings, ValidationFinding{Field: field, Message: msg, Severity: SeverityWarning, Code: code})
}

func (res *ValidationResult) addInfo(field, msg, code string) {
	res.Infos = append(res.Infos, ValidationFinding{Field: field, Message: msg, Severity: SeverityInfo, Code: code})
}

func (sy *SkillYAML) supportsPlatform(platform Platform) bool {
	normalized := NormalizePlatform(platform)
	for _, p := range sy.CompatibleWith {
		np := NormalizePlatform(p)
		if np == normalized || np == PlatformAll {
			return true
		}
	}
	return false
}
