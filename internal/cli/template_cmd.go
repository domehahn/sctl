package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// templateRegistryFile is stored alongside the skpm config.
const templateRegistryFile = "templates.yaml"

type skillTemplate struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Source      string            `yaml:"source,omitempty"` // URL or local path (informational)
	AddedAt     string            `yaml:"added_at,omitempty"`
	Files       map[string]string `yaml:"files"`
}

type templateRegistry struct {
	Templates map[string]skillTemplate `yaml:"templates"`
}

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage reusable skill scaffolding templates",
		Long: `Templates extend 'skpm create' with reusable scaffolding for common skill patterns.

  skpm template list                  — list saved templates
  skpm template add <name> <path>     — save a skill directory as a template
  skpm template remove <name>         — delete a saved template
  skpm template show <name>           — print the files in a template
  skpm template use <name> <skill>    — scaffold a new skill from a template`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"template"}, args...), func() error { return nil })
		},
	}
	cmd.AddCommand(newTemplateListCmd())
	cmd.AddCommand(newTemplateAddCmd())
	cmd.AddCommand(newTemplateRemoveCmd())
	cmd.AddCommand(newTemplateShowCmd())
	cmd.AddCommand(newTemplateUseCmd())
	return cmd
}

func newTemplateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List saved templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := loadTemplateRegistry()
			if err != nil {
				return err
			}
			if len(reg.Templates) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No templates saved. Use 'skpm template add <name> <path>' to add one.")
				return nil
			}
			names := sortedTemplateNames(reg)
			for _, name := range names {
				t := reg.Templates[name]
				desc := t.Description
				if desc == "" {
					desc = "(no description)"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %-24s %s\n", name, desc)
				if t.Source != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  source: %s\n", strings.Repeat(" ", 24), t.Source)
				}
			}
			return nil
		},
	}
}

func newTemplateAddCmd() *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "add <name> <path>",
		Short: "Save a skill directory as a named template",
		Long: `Reads the skill files from <path> (SKILL.md, skill.yaml, VERSION, CHANGELOG.md
and any additional files) and saves them as a reusable template.

The template replaces occurrences of the skill name in file content with the
placeholder {{.Name}} so 'skpm template use' can substitute it.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, srcDir := args[0], args[1]

			if _, err := os.Stat(srcDir); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("directory not found: %s", srcDir)}
			}

			// Read all text files from the skill dir.
			files := map[string]string{}
			entries, err := os.ReadDir(srcDir)
			if err != nil {
				return &InternalError{Message: "read directory", Cause: err}
			}

			// Detect the skill name from skill.yaml to parameterise content.
			skillName := ""
			if sy, err := readSkillYAML(srcDir); err == nil {
				skillName = sy.Name
			}

			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				fname := e.Name()
				if !isTemplateFile(fname) {
					continue
				}
				content, err := os.ReadFile(filepath.Join(srcDir, fname))
				if err != nil {
					continue
				}
				text := string(content)
				if skillName != "" {
					text = strings.ReplaceAll(text, skillName, "{{.Name}}")
				}
				files[fname] = text
			}

			if len(files) == 0 {
				return &UserError{Message: fmt.Sprintf("no skill files found in %s", srcDir)}
			}

			reg, err := loadTemplateRegistry()
			if err != nil {
				return err
			}

			reg.Templates[name] = skillTemplate{
				Name:        name,
				Description: description,
				Source:      srcDir,
				AddedAt:     time.Now().UTC().Format("2006-01-02"),
				Files:       files,
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would save template %q (%d files)\n", name, len(files))
				return nil
			}

			if err := saveTemplateRegistry(reg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved template %q (%d files)\n", name, len(files))
			return nil
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Short description of the template")
	return cmd
}

func newTemplateRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a saved template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			reg, err := loadTemplateRegistry()
			if err != nil {
				return err
			}
			if _, ok := reg.Templates[name]; !ok {
				return &UserError{Message: fmt.Sprintf("template %q not found", name)}
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove template %q\n", name)
				return nil
			}
			delete(reg.Templates, name)
			if err := saveTemplateRegistry(reg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed template %q\n", name)
			return nil
		},
	}
}

func newTemplateShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Print the file list (and content) of a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			reg, err := loadTemplateRegistry()
			if err != nil {
				return err
			}
			t, ok := reg.Templates[name]
			if !ok {
				return &UserError{Message: fmt.Sprintf("template %q not found", name)}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Template: %s\n", t.Name)
			if t.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", t.Description)
			}
			if t.Source != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Source: %s\n", t.Source)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Files (%d):\n", len(t.Files))
			for _, fname := range sortedKeys(t.Files) {
				lines := strings.Count(t.Files[fname], "\n")
				fmt.Fprintf(cmd.OutOrStdout(), "  %-24s (%d lines)\n", fname, lines)
			}
			return nil
		},
	}
}

func newTemplateUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <template> <skill-name>",
		Short: "Scaffold a new skill directory from a saved template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName, skillName := args[0], args[1]

			reg, err := loadTemplateRegistry()
			if err != nil {
				return err
			}
			t, ok := reg.Templates[templateName]
			if !ok {
				return &UserError{Message: fmt.Sprintf("template %q not found; run 'skpm template list' to see available templates", templateName)}
			}

			destDir := skillName
			if _, err := os.Stat(destDir); err == nil {
				return &UserError{Message: fmt.Sprintf("directory %q already exists", destDir)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would scaffold %q from template %q (%d files)\n", skillName, templateName, len(t.Files))
				return nil
			}

			if err := os.MkdirAll(destDir, 0o755); err != nil {
				return &InternalError{Message: "create directory", Cause: err}
			}

			for fname, content := range t.Files {
				// Substitute {{.Name}} placeholder.
				rendered := strings.ReplaceAll(content, "{{.Name}}", skillName)
				path := filepath.Join(destDir, fname)
				if err := os.WriteFile(path, []byte(rendered), 0o644); err != nil {
					os.RemoveAll(destDir)
					return &InternalError{Message: fmt.Sprintf("write %s", fname), Cause: err}
				}
			}

			abs, _ := filepath.Abs(destDir)
			fmt.Fprintf(cmd.OutOrStdout(), "Created skill %q from template %q in %s\n\n", skillName, templateName, abs)
			fmt.Fprintf(cmd.OutOrStdout(), "Next steps:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  1. Edit %s/SKILL.md\n", destDir)
			fmt.Fprintf(cmd.OutOrStdout(), "  2. skpm validate %s\n", destDir)
			fmt.Fprintf(cmd.OutOrStdout(), "  3. skpm publish %s\n", destDir)
			return nil
		},
	}
}

// ── Storage helpers ───────────────────────────────────────────────────────────

func templateRegistryPath() (string, error) {
	cfgPath, err := configFilePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(cfgPath), templateRegistryFile), nil
}

func loadTemplateRegistry() (*templateRegistry, error) {
	path, err := templateRegistryPath()
	if err != nil {
		return nil, &InternalError{Message: "resolve template registry path", Cause: err}
	}
	reg := &templateRegistry{Templates: map[string]skillTemplate{}}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return reg, nil
	}
	if err != nil {
		return nil, &InternalError{Message: "read template registry", Cause: err}
	}
	if err := yaml.Unmarshal(data, reg); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse template registry: %v", err)}
	}
	if reg.Templates == nil {
		reg.Templates = map[string]skillTemplate{}
	}
	return reg, nil
}

func saveTemplateRegistry(reg *templateRegistry) error {
	path, err := templateRegistryPath()
	if err != nil {
		return &InternalError{Message: "resolve template registry path", Cause: err}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return &InternalError{Message: "create config dir", Cause: err}
	}
	data, err := yaml.Marshal(reg)
	if err != nil {
		return &InternalError{Message: "marshal template registry", Cause: err}
	}
	return writeAtomic(path, string(data))
}

func isTemplateFile(name string) bool {
	switch strings.ToLower(name) {
	case "skill.md", "skill.yaml", "version", "changelog.md", "readme.md":
		return true
	}
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".yaml" || ext == ".yml" || ext == ".txt"
}

func sortedTemplateNames(reg *templateRegistry) []string {
	names := make([]string, 0, len(reg.Templates))
	for k := range reg.Templates {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
