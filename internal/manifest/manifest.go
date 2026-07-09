package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/domehahn/sklib/spec"
	"gopkg.in/yaml.v3"
)

const DefaultFilename = "agent-skills.yaml"
const OverrideFilename = ".skill-override.yaml"

// RegistryConfig is aliased from sklib/spec for use in agent-skills.yaml.
type RegistryConfig = spec.RegistryConfig

// SkillOverride holds user customizations applied at install time without
// modifying the original SKILL.md. The agent reads .skill-override.yaml
// at runtime (e.g., via a Pi extension) to inject these values.
type SkillOverride struct {
	// Prepend is prepended to the skill's system prompt.
	Prepend string `yaml:"prepend,omitempty"`
	// Append is appended to the skill's system prompt.
	Append string `yaml:"append,omitempty"`
	// Env are environment variables to set when the skill is active.
	// Values support "${VAR}" or "env:VAR" syntax for referencing
	// the host environment at runtime.
	Env map[string]string `yaml:"env,omitempty"`
}

// ManifestFile is the Go model for agent-skills.yaml.
// Field layout matches schemas/agent-skills.schema.json from skillspec.
type ManifestFile struct {
	Version         int                       `yaml:"version"`
	DefaultRegistry string                    `yaml:"default_registry,omitempty"`
	Registries      map[string]RegistryConfig `yaml:"registries,omitempty"`
	Skills          []SkillEntry              `yaml:"skills"`
	Overrides       map[string]SkillOverride  `yaml:"overrides,omitempty"`
	Hooks           map[string]string         `yaml:"hooks,omitempty"`
	Metadata        map[string]string         `yaml:"metadata,omitempty"`
}

// SkillEntry is aliased from sklib/spec so the canonical agent-skills.yaml schema
// is shared across skpm, skcr, and SkillForge.
type SkillEntry = spec.ManifestSkill

func New() *ManifestFile {
	return &ManifestFile{Version: 1}
}

func Read(path string) (*ManifestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", path, err)
	}
	var m ManifestFile
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	if m.Version == 0 {
		m.Version = 1
	}
	return &m, nil
}

func (m *ManifestFile) Write(path string) error {
	m.Sort()
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create manifest dir: %w", err)
	}
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write manifest tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("atomic rename manifest: %w", err)
	}
	return nil
}

func (m *ManifestFile) Sort() {
	sort.SliceStable(m.Skills, func(i, j int) bool {
		return m.Skills[i].Name < m.Skills[j].Name
	})
}

func (m *ManifestFile) Upsert(entry SkillEntry) {
	for i, s := range m.Skills {
		if s.Name == entry.Name {
			m.Skills[i] = entry
			return
		}
	}
	m.Skills = append(m.Skills, entry)
}

func (m *ManifestFile) Remove(name string) bool {
	for i, s := range m.Skills {
		if s.Name == name {
			m.Skills = append(m.Skills[:i], m.Skills[i+1:]...)
			return true
		}
	}
	return false
}

// OverrideFor returns the override configuration for a named skill, or nil if none exists.
func (m *ManifestFile) OverrideFor(skillName string) *SkillOverride {
	if m.Overrides == nil {
		return nil
	}
	ov, ok := m.Overrides[skillName]
	if !ok {
		return nil
	}
	return &ov
}

// WriteOverride writes the .skill-override.yaml file to the given directory.
// If the override is empty (no prepend, append, env), the file is removed
// instead to avoid clutter.
func WriteOverride(dir string, ov *SkillOverride) error {
	path := filepath.Join(dir, OverrideFilename)
	if ov == nil || (ov.Prepend == "" && ov.Append == "" && len(ov.Env) == 0) {
		os.Remove(path) // ignore error — stale file is harmless
		return nil
	}
	data, err := yaml.Marshal(ov)
	if err != nil {
		return fmt.Errorf("marshal override: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write override %s: %w", path, err)
	}
	return nil
}
