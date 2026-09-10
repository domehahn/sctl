package skill

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/sklib/validate"
	"gopkg.in/yaml.v3"
)

// decodeSkillMetadata adapts the versioned skil wire identity for packaging.
// It leaves source bytes untouched; capability/assurance semantics belong to skil.
func decodeSkillMetadata(data []byte) (*SkillYAML, bool, error) {
	var root map[string]yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, false, err
	}
	if _, native := root["skill"]; !native {
		var legacy SkillYAML
		err := yaml.Unmarshal(data, &legacy)
		return &legacy, false, err
	}
	var doc struct {
		Version int `yaml:"version"`
		Skill   struct {
			Name        string `yaml:"name"`
			Version     string `yaml:"version"`
			Description string `yaml:"description"`
		} `yaml:"skill"`
		Entrypoint    string `yaml:"entrypoint"`
		Compatibility struct {
			Platforms []Platform `yaml:"platforms"`
		} `yaml:"compatibility"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, true, err
	}
	if doc.Version != 1 {
		return nil, true, fmt.Errorf("unsupported skil schema version %d", doc.Version)
	}
	if _, ambiguous := root["name"]; ambiguous {
		return nil, true, fmt.Errorf("ambiguous native and legacy metadata")
	}
	return &SkillYAML{Name: doc.Skill.Name, Version: doc.Skill.Version, Description: doc.Skill.Description, Entrypoint: doc.Entrypoint, CompatibleWith: doc.Compatibility.Platforms}, true, nil
}

func (v *StructuredValidator) validateNative(dir string, sy *SkillYAML, res *ValidationResult) {
	result := validate.ValidateSkillMetadata(skillYAMLToSpec(sy), validate.Options{Strict: true})
	for _, f := range result.Findings {
		if f.Severity == validate.SeverityError {
			res.addError("skill.yaml", f.Message, f.Code)
		}
	}
	version, ok := v.validateVersion(dir, res)
	if ok && version != sy.Version {
		res.addError("VERSION", "native skill version differs from VERSION", "version_mismatch")
	}
	v.validatePlatforms(sy, res)
	if v.options.Platform != "" && !sy.supportsPlatform(v.options.Platform) {
		res.addError("skill.yaml", "unsupported platform", "platform_not_compatible")
	}
	entry := sy.Entrypoint
	if entry == "" {
		entry = "SKILL.md"
	}
	if !filepath.IsLocal(entry) {
		res.addError("skill.yaml", "entrypoint must stay inside package", "invalid_entrypoint")
	} else if info, err := os.Stat(filepath.Join(dir, entry)); err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
		res.addError(entry, "entrypoint missing or empty", "missing_entrypoint")
	}
	// Native compiler output legitimately includes its content checksums.
	// Keep all hygiene checks; filter only that exact compiler-generated file.
	before := len(res.Errors)
	v.validateHygiene(dir, res)
	kept := res.Errors[:before]
	for _, f := range res.Errors[before:] {
		if f.Field == "checksums.txt" && f.Code == "generated_artifact" {
			continue
		}
		kept = append(kept, f)
	}
	res.Errors = kept
	res.Valid = len(res.Errors) == 0
}
