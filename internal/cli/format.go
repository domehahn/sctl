package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type formatChange struct {
	File   string
	Before string
	After  string
	Reason string
}

func newFormatCmd() *cobra.Command {
	var check bool
	var write bool

	cmd := &cobra.Command{
		Use:   "format [path]",
		Short: "Normalize skill metadata without changing semantics",
		Long: `Normalizes a skill's metadata files:

  - Normalizes platform aliases to canonical names (e.g. gitlab → gitlab-duo)
  - Removes duplicate platforms and tags from skill.yaml
  - Adds entrypoint: SKILL.md to skill.yaml if absent
  - Removes leading 'v' from VERSION (e.g. v1.2.3 → 1.2.3)
  - Ensures a trailing newline in text files

Defaults to --check (dry-run). Use --write to apply changes.

Note: re-marshaling skill.yaml removes inline comments.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"format"}, args...), func() error {
				dir := "."
				if len(args) == 1 {
					dir = args[0]
				}
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					return &UserError{Message: fmt.Sprintf("path not found: %s", dir)}
				}
				if !check && !write {
					check = true
				}

				changes, err := computeFormatChanges(dir)
				if err != nil {
					return &InternalError{Message: "format", Cause: err}
				}

				format := outputFormat()
				if format == OutputJSON {
					type jsonChange struct {
						File   string `json:"file"`
						Reason string `json:"reason"`
					}
					list := make([]jsonChange, len(changes))
					for i, c := range changes {
						list[i] = jsonChange{File: c.File, Reason: c.Reason}
					}
					PrintResult(format, CommandResult{
						Success: len(changes) == 0 || write,
						Command: "format",
						Data: map[string]interface{}{
							"path":    dir,
							"changes": list,
							"written": write && len(changes) > 0,
						},
					})
					if check && len(changes) > 0 {
						os.Exit(1)
					}
					return nil
				}

				if len(changes) == 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "No formatting changes needed: %s\n", dir)
					return nil
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Files that would be reformatted in %s:\n", dir)
				for _, c := range changes {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  (%s)\n", c.File, c.Reason)
				}

				if write {
					for _, c := range changes {
						path := filepath.Join(dir, c.File)
						if err := os.WriteFile(path, []byte(c.After), 0o644); err != nil {
							return &InternalError{Message: fmt.Sprintf("write %s", c.File), Cause: err}
						}
					}
					fmt.Fprintf(cmd.OutOrStdout(), "\nFormatted %d file(s) in %s\n", len(changes), dir)
					return nil
				}

				fmt.Fprintf(cmd.OutOrStdout(), "\n%d file(s) would be reformatted. Run with --write to apply.\n", len(changes))
				os.Exit(1)
				return nil
			})
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "Report files that need formatting without writing (exit 1 if any)")
	cmd.Flags().BoolVar(&write, "write", false, "Apply formatting changes to disk")
	return cmd
}

func computeFormatChanges(dir string) ([]formatChange, error) {
	var changes []formatChange

	// VERSION: strip leading 'v', ensure trailing newline
	versionPath := filepath.Join(dir, "VERSION")
	if raw, err := os.ReadFile(versionPath); err == nil {
		original := string(raw)
		normalized := strings.TrimPrefix(strings.TrimSpace(original), "v")
		if !strings.HasSuffix(normalized, "\n") {
			normalized += "\n"
		}
		if original != normalized {
			changes = append(changes, formatChange{
				File:   "VERSION",
				Before: original,
				After:  normalized,
				Reason: describeVersionChange(original, normalized),
			})
		}
	}

	// skill.yaml: normalize platforms, deduplicate, add entrypoint
	yamlPath := filepath.Join(dir, "skill.yaml")
	if raw, err := os.ReadFile(yamlPath); err == nil {
		var sy skill.SkillYAML
		if parseErr := yaml.Unmarshal(raw, &sy); parseErr == nil {
			changed := false

			// Normalize platform aliases and deduplicate
			seen := map[skill.Platform]bool{}
			var normalized []skill.Platform
			for _, p := range sy.CompatibleWith {
				canon := skill.NormalizePlatform(p)
				if canon != p {
					changed = true
				}
				if !seen[canon] {
					seen[canon] = true
					normalized = append(normalized, canon)
				} else {
					changed = true
				}
			}
			if changed {
				sy.CompatibleWith = normalized
			}

			// Deduplicate tags
			seenTags := map[string]bool{}
			var dedupTags []string
			for _, t := range sy.Tags {
				if !seenTags[t] {
					seenTags[t] = true
					dedupTags = append(dedupTags, t)
				}
			}
			if len(dedupTags) != len(sy.Tags) {
				sy.Tags = dedupTags
				changed = true
			}

			// Add default entrypoint if absent
			if sy.Entrypoint == "" {
				sy.Entrypoint = "SKILL.md"
				changed = true
			}

			if changed {
				remarshaled, marshalErr := yaml.Marshal(&sy)
				if marshalErr == nil && !bytes.Equal(raw, remarshaled) {
					changes = append(changes, formatChange{
						File:   "skill.yaml",
						Before: string(raw),
						After:  string(remarshaled),
						Reason: "normalized platforms, tags, and/or added default entrypoint",
					})
				}
			}
		}
	}

	// Text files: ensure trailing newline
	for _, name := range []string{"SKILL.md", "CHANGELOG.md", "README.md", "LICENSE"} {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		if len(raw) > 0 && raw[len(raw)-1] != '\n' {
			changes = append(changes, formatChange{
				File:   name,
				Before: string(raw),
				After:  string(raw) + "\n",
				Reason: "added trailing newline",
			})
		}
	}

	return changes, nil
}

func describeVersionChange(original, normalized string) string {
	reasons := []string{}
	if strings.HasPrefix(strings.TrimSpace(original), "v") {
		reasons = append(reasons, "removed leading 'v'")
	}
	if !strings.HasSuffix(strings.TrimSpace(original), "\n") {
		reasons = append(reasons, "added trailing newline")
	}
	if len(reasons) == 0 {
		return "normalized"
	}
	return strings.Join(reasons, ", ")
}
