package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const DefaultFilename = "agent-skills.yaml"

type ManifestFile struct {
	Version int             `yaml:"version"`
	Skills  []SkillEntry    `yaml:"skills"`
}

type SkillEntry struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version,omitempty"` // empty = latest
	Source  string `yaml:"source,omitempty"`  // empty = default_registry
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
