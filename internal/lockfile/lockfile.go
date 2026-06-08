package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFilename = "agent-skills.lock"
	currentVersion  = 1
)

type LockFile struct {
	Version     int         `yaml:"version"`
	GeneratedAt string      `yaml:"generated_at,omitempty"`
	Skills      []SkillLock `yaml:"skills"`
}

type SkillLock struct {
	Name           string            `yaml:"name"`
	Namespace      string            `yaml:"namespace,omitempty"`
	Version        string            `yaml:"version"`
	Source         string            `yaml:"source"`
	RegistryType   string            `yaml:"registry_type,omitempty"`
	SourceURL      string            `yaml:"source_url"`
	ArtifactName   string            `yaml:"artifact_name,omitempty"`
	SHA256         string            `yaml:"sha256"`
	PackageType    string            `yaml:"package_type,omitempty"`
	CompatibleWith []string          `yaml:"compatible_with,omitempty"`
	InstalledTo    []string          `yaml:"installed_to,omitempty"`
	Metadata       map[string]string `yaml:"metadata,omitempty"`
}

func New() *LockFile {
	return &LockFile{Version: currentVersion}
}

func Read(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lockfile %s: %w", path, err)
	}
	lf := &LockFile{}
	if err := yaml.Unmarshal(data, lf); err != nil {
		return nil, fmt.Errorf("parse lockfile %s: %w", path, err)
	}
	if lf.Version == 0 {
		lf.Version = currentVersion
	}
	return lf, nil
}

func (lf *LockFile) Write(path string) error {
	lf.Sort()
	data, err := yaml.Marshal(lf)
	if err != nil {
		return fmt.Errorf("marshal lockfile: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lockfile dir: %w", err)
	}
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write lockfile tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("atomic rename lockfile: %w", err)
	}
	return nil
}

func (lf *LockFile) Sort() {
	sort.SliceStable(lf.Skills, func(i, j int) bool {
		return lf.Skills[i].Name < lf.Skills[j].Name
	})
}

func (lf *LockFile) Upsert(lock SkillLock) {
	for i, s := range lf.Skills {
		if s.Name == lock.Name {
			lf.Skills[i] = lock
			return
		}
	}
	lf.Skills = append(lf.Skills, lock)
}

func (lf *LockFile) Find(name string) (*SkillLock, bool) {
	for i := range lf.Skills {
		if lf.Skills[i].Name == name {
			return &lf.Skills[i], true
		}
	}
	return nil, false
}

func (lf *LockFile) Remove(name string) bool {
	for i, s := range lf.Skills {
		if s.Name == name {
			lf.Skills = append(lf.Skills[:i], lf.Skills[i+1:]...)
			return true
		}
	}
	return false
}
