package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// currentSpecVersion is the spec version skpm currently targets.
const currentSpecVersion = 1

func newMigrateCmd() *cobra.Command {
	var (
		dryRun   bool
		skillDir string
	)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate project or skill files to the current spec version",
		Long: `Detects the spec version of local files and applies incremental migrations
to bring them up to date.

  Project migration (default) — upgrades agent-skills.yaml and agent-skills.lock
  Skill migration (--skill-dir) — upgrades skill.yaml and SKILL.md frontmatter

Migrations are idempotent: running migrate on already-current files is safe.
Use --dry-run to preview changes without writing any files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dry := dryRun || globalDryRun
			if skillDir != "" {
				return migrateSkillDir(cmd, skillDir, dry)
			}
			return migrateProject(cmd, dry)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing files")
	cmd.Flags().StringVar(&skillDir, "skill-dir", "", "Migrate a single skill directory instead of the project")
	return cmd
}

// ── Project migration ─────────────────────────────────────────────────────────

func migrateProject(cmd *cobra.Command, dry bool) error {
	changed := false

	// agent-skills.yaml
	if _, err := os.Stat("agent-skills.yaml"); err == nil {
		ok, err := migrateManifestFile(cmd, "agent-skills.yaml", dry)
		if err != nil {
			return err
		}
		changed = changed || ok
	}

	// agent-skills.lock
	if _, err := os.Stat("agent-skills.lock"); err == nil {
		ok, err := migrateLockFile(cmd, "agent-skills.lock", dry)
		if err != nil {
			return err
		}
		changed = changed || ok
	}

	if !changed {
		fmt.Fprintln(cmd.OutOrStdout(), "Already up to date.")
	}
	return nil
}

func migrateManifestFile(cmd *cobra.Command, path string, dry bool) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, &InternalError{Message: "read " + path, Cause: err}
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return false, &UserError{Message: fmt.Sprintf("parse %s: %v", path, err)}
	}

	mapping := yamlMapping(&doc)
	if mapping == nil {
		return false, &UserError{Message: path + ": unexpected YAML structure"}
	}

	version := 0
	if v, ok := findYAMLScalar(&doc, "version"); ok {
		fmt.Sscanf(v, "%d", &version) //nolint:errcheck
	}

	if version >= currentSpecVersion {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: version %d — no migration needed\n", path, version)
		return false, nil
	}

	migrations := manifestMigrations()
	applied := []string{}
	for from := version; from < currentSpecVersion; from++ {
		m, ok := migrations[from]
		if !ok {
			continue
		}
		if err := m(&doc); err != nil {
			return false, &InternalError{Message: fmt.Sprintf("migration v%d→v%d", from, from+1), Cause: err}
		}
		applied = append(applied, fmt.Sprintf("v%d→v%d", from, from+1))
	}

	if len(applied) == 0 {
		return false, nil
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return false, &InternalError{Message: "marshal " + path, Cause: err}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s: applied %s\n", path, strings.Join(applied, ", "))
	if dry {
		fmt.Fprintf(cmd.OutOrStdout(), "  (dry run — not written)\n")
		return true, nil
	}
	if err := writeAtomic(path, string(out)); err != nil {
		return false, &InternalError{Message: "write " + path, Cause: err}
	}
	return true, nil
}

func migrateLockFile(cmd *cobra.Command, path string, dry bool) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, &InternalError{Message: "read " + path, Cause: err}
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return false, &UserError{Message: fmt.Sprintf("parse %s: %v", path, err)}
	}

	version := 0
	if v, ok := findYAMLScalar(&doc, "version"); ok {
		fmt.Sscanf(v, "%d", &version) //nolint:errcheck
	}

	if version >= currentSpecVersion {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: version %d — no migration needed\n", path, version)
		return false, nil
	}

	migrations := lockfileMigrations()
	applied := []string{}
	for from := version; from < currentSpecVersion; from++ {
		m, ok := migrations[from]
		if !ok {
			continue
		}
		if err := m(&doc); err != nil {
			return false, &InternalError{Message: fmt.Sprintf("migration v%d→v%d", from, from+1), Cause: err}
		}
		applied = append(applied, fmt.Sprintf("v%d→v%d", from, from+1))
	}

	if len(applied) == 0 {
		return false, nil
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return false, &InternalError{Message: "marshal " + path, Cause: err}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s: applied %s\n", path, strings.Join(applied, ", "))
	if dry {
		fmt.Fprintf(cmd.OutOrStdout(), "  (dry run — not written)\n")
		return true, nil
	}
	if err := writeAtomic(path, string(out)); err != nil {
		return false, &InternalError{Message: "write " + path, Cause: err}
	}
	return true, nil
}

// ── Skill directory migration ─────────────────────────────────────────────────

func migrateSkillDir(cmd *cobra.Command, dir string, dry bool) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return &UserError{Message: fmt.Sprintf("directory not found: %s", dir)}
	}

	changed := false

	// skill.yaml version field
	syPath := filepath.Join(dir, "skill.yaml")
	if _, err := os.Stat(syPath); err == nil {
		ok, err := migrateSkillYAMLFile(cmd, syPath, dry)
		if err != nil {
			return err
		}
		changed = changed || ok
	}

	// SKILL.md frontmatter compatibility field
	mdPath := filepath.Join(dir, "SKILL.md")
	if _, err := os.Stat(mdPath); err == nil {
		ok, err := migrateSkillMD(cmd, mdPath, dry)
		if err != nil {
			return err
		}
		changed = changed || ok
	}

	if !changed {
		fmt.Fprintln(cmd.OutOrStdout(), "Already up to date.")
	}
	return nil
}

func migrateSkillYAMLFile(cmd *cobra.Command, path string, dry bool) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, &InternalError{Message: "read " + path, Cause: err}
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return false, &UserError{Message: fmt.Sprintf("parse %s: %v", path, err)}
	}

	changes := []string{}

	// Ensure 'namespace' field exists (added in spec v1).
	if _, ok := findYAMLScalar(&doc, "namespace"); !ok {
		mapping := yamlMapping(&doc)
		if mapping != nil {
			mapping.Content = append(
				[]*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "namespace"},
					{Kind: yaml.ScalarNode, Value: "default"},
				},
				mapping.Content...,
			)
			changes = append(changes, "added namespace: default")
		}
	}

	// Rename 'platforms' → 'compatible_with' if present (pre-v1 field name).
	if _, ok := findYAMLScalar(&doc, "platforms"); ok {
		if _, alreadyHas := findYAMLScalar(&doc, "compatible_with"); !alreadyHas {
			renameYAMLKey(&doc, "platforms", "compatible_with")
			changes = append(changes, "renamed platforms → compatible_with")
		}
	}

	if len(changes) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: no changes needed\n", path)
		return false, nil
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return false, &InternalError{Message: "marshal " + path, Cause: err}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", path, strings.Join(changes, "; "))
	if dry {
		fmt.Fprintf(cmd.OutOrStdout(), "  (dry run — not written)\n")
		return true, nil
	}
	return true, writeAtomic(path, string(out))
}

func migrateSkillMD(cmd *cobra.Command, path string, dry bool) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, &InternalError{Message: "read " + path, Cause: err}
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: no frontmatter — skipping\n", path)
		return false, nil
	}

	// Find closing ---
	rest := content[4:]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: malformed frontmatter — skipping\n", path)
		return false, nil
	}

	frontmatter := rest[:end]
	body := rest[end+5:]

	changes := []string{}

	// Ensure 'compatibility.spec_version' is present (added in spec v1).
	if !strings.Contains(frontmatter, "spec_version:") {
		if strings.Contains(frontmatter, "compatibility:") {
			// Inject under existing compatibility block.
			frontmatter = strings.Replace(frontmatter,
				"compatibility:",
				"compatibility:\n  spec_version: 1",
				1,
			)
		} else {
			frontmatter += "\ncompatibility:\n  spec_version: 1"
		}
		changes = append(changes, "added compatibility.spec_version: 1")
	}

	if len(changes) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: no changes needed\n", path)
		return false, nil
	}

	newContent := "---\n" + frontmatter + "\n---\n" + body
	fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", path, strings.Join(changes, "; "))
	if dry {
		fmt.Fprintf(cmd.OutOrStdout(), "  (dry run — not written)\n")
		return true, nil
	}
	return true, writeAtomic(path, newContent)
}

// ── Migration tables ──────────────────────────────────────────────────────────

type yamlMigrationFn func(doc *yaml.Node) error

// manifestMigrations returns a map of from-version → migration function.
func manifestMigrations() map[int]yamlMigrationFn {
	return map[int]yamlMigrationFn{
		// v0 → v1: set version: 1
		0: func(doc *yaml.Node) error {
			setYAMLScalar(doc, "version", "1")
			return nil
		},
	}
}

// lockfileMigrations returns a map of from-version → migration function.
func lockfileMigrations() map[int]yamlMigrationFn {
	return map[int]yamlMigrationFn{
		// v0 → v1: set version: 1
		0: func(doc *yaml.Node) error {
			setYAMLScalar(doc, "version", "1")
			return nil
		},
	}
}

// renameYAMLKey renames a mapping key in-place without changing its value node.
func renameYAMLKey(doc *yaml.Node, oldKey, newKey string) {
	mapping := yamlMapping(doc)
	if mapping == nil {
		return
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == oldKey {
			mapping.Content[i].Value = newKey
			return
		}
	}
}
