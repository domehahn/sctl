package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type manifestEntry struct {
	Name        string   `json:"name" yaml:"name"`
	Version     string   `json:"version" yaml:"version"`
	SHA256      string   `json:"sha256" yaml:"sha256"`
	Source      string   `json:"source" yaml:"source"`
	InstalledTo []string `json:"installed_to,omitempty" yaml:"installed_to,omitempty"`
	Platforms   []string `json:"platforms,omitempty" yaml:"platforms,omitempty"`
	Signature   string   `json:"signature,omitempty" yaml:"signature,omitempty"`
}

type deployManifest struct {
	GeneratedAt string          `json:"generated_at" yaml:"generated_at"`
	LockFile    string          `json:"lockfile" yaml:"lockfile"`
	Skills      []manifestEntry `json:"skills" yaml:"skills"`
}

func newManifestCmd() *cobra.Command {
	var lockPath string
	var format string
	var outPath string

	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Generate a deployment manifest from the lockfile",
		Long: `Produces a structured summary of every installed skill — name, version,
SHA-256 digest, install paths, and compatible platforms — as JSON or YAML.

Suitable for feeding into deployment pipelines, audit logs, or release notes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			skills := make([]lockfile.SkillLock, len(lf.Skills))
			copy(skills, lf.Skills)
			sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })

			entries := make([]manifestEntry, 0, len(skills))
			for _, sl := range skills {
				platforms := make([]string, len(sl.CompatibleWith))
				for i, p := range sl.CompatibleWith {
					platforms[i] = string(p)
				}
				entries = append(entries, manifestEntry{
					Name:        sl.Name,
					Version:     sl.Version,
					SHA256:      sl.SHA256,
					Source:      sl.Source,
					InstalledTo: sl.InstalledTo,
					Platforms:   platforms,
					Signature:   sl.Signature,
				})
			}

			absLock, _ := filepath.Abs(lockPath)
			m := deployManifest{
				GeneratedAt: time.Now().UTC().Format(time.RFC3339),
				LockFile:    absLock,
				Skills:      entries,
			}

			var data []byte
			switch strings.ToLower(format) {
			case "yaml", "yml":
				data, err = yaml.Marshal(m)
			default:
				enc := json.NewEncoder(nil)
				var buf strings.Builder
				e := json.NewEncoder(&buf)
				e.SetIndent("", "  ")
				if encErr := e.Encode(m); encErr != nil {
					return &InternalError{Message: "marshal manifest", Cause: encErr}
				}
				_ = enc
				data = []byte(buf.String())
			}
			if err != nil {
				return &InternalError{Message: "marshal manifest", Cause: err}
			}

			if outPath != "" {
				if writeErr := os.WriteFile(outPath, data, 0o644); writeErr != nil {
					return &InternalError{Message: "write manifest", Cause: writeErr}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Manifest written to %s (%d skill(s))\n", outPath, len(entries))
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), string(data))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&format, "format", "json", "Output format: json or yaml")
	cmd.Flags().StringVar(&outPath, "out", "", "Write manifest to file instead of stdout")
	return cmd
}
