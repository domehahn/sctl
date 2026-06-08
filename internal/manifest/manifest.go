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

// RegistryConfig is aliased from sklib/spec for use in agent-skills.yaml.
type RegistryConfig = spec.RegistryConfig

// ManifestFile is the Go model for agent-skills.yaml.
// Field layout matches schemas/agent-skills.schema.json from skillspec.
type ManifestFile struct {
	Version         int                       `yaml:"version"`
	DefaultRegistry string                    `yaml:"default_registry,omitempty"`
	Registries      map[string]RegistryConfig `yaml:"registries,omitempty"`
	Skills          []SkillEntry              `yaml:"skills"`
	Metadata        map[string]string         `yaml:"metadata,omitempty"`
}

// SkillEntry is a single skill declared in agent-skills.yaml.
type SkillEntry struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace,omitempty"`
	Version   string            `yaml:"version,omitempty"`  // version constraint; empty = latest
	Source    string            `yaml:"source,omitempty"`   // registry alias; empty = default_registry
	Target    string            `yaml:"target,omitempty"`   // install path override
	Platforms []spec.Platform   `yaml:"platforms,omitempty"`
	Path      string            `yaml:"path,omitempty"`     // local path (local registry)
	Ref       string            `yaml:"ref,omitempty"`      // git ref (local registry)
	Metadata  map[string]string `yaml:"metadata,omitempty"`
}

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
