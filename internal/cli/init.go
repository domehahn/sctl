package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/domehahn/sctl/internal/lockfile"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize sctl config, a project lockfile, or a new skill",
		Long: `Scaffold sctl configuration and project files interactively.

  sctl init config          — create ~/.config/sctl/config.yaml
  sctl init project         — create agent-skills.lock in the current directory
  sctl init skill <name>    — scaffold a new skill directory`,
	}
	cmd.AddCommand(newInitConfigCmd())
	cmd.AddCommand(newInitProjectCmd())
	cmd.AddCommand(newInitSkillCmd())
	return cmd
}

// ── sctl init config ──────────────────────────────────────────────────────

func newInitConfigCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Create ~/.config/sctl/config.yaml interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := configFilePath()
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}

			if !force {
				if _, err := os.Stat(cfgPath); err == nil {
					return &UserError{Message: fmt.Sprintf(
						"config already exists at %s\nUse --force to overwrite.", cfgPath,
					)}
				}
			}

			p := newPrompter(cmd)
			p.header("sctl config init")
			p.print("This creates %s\n\n", cfgPath)

			registryType := p.choose("Registry type", []string{"gitlab", "github", "artifactory", "local"})
			registryName := p.ask("Registry name (used as key in config)", "company-"+registryType)

			var registryURL, registryProject string
			switch registryType {
			case "gitlab":
				registryURL = p.ask("GitLab base URL", "https://gitlab.company.com")
				registryProject = p.ask("Project path (namespace/project)", "platform/agent-skills")
			case "github":
				registryURL = p.ask("GitHub owner/repo", "myorg/agent-skills")
			case "artifactory":
				base := p.ask("Artifactory base URL", "https://artifactory.company.com/artifactory")
				repo := p.ask("Repository name", "agent-skills")
				registryURL = base + "#" + repo
			case "local":
				registryURL = p.ask("Base directory path", "./agent-skills-local")
			}

			cacheDir := p.ask("Cache directory", "~/.cache/sctl")

			cfg := buildConfigYAML(registryName, registryType, registryURL, registryProject, cacheDir)

			if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
				return &InternalError{Message: "create config dir", Cause: err}
			}
			if err := writeAtomic(cfgPath, cfg); err != nil {
				return &InternalError{Message: "write config", Cause: err}
			}

			p.print("\n✓ Created %s\n", cfgPath)
			p.print("  Default registry: %s\n", registryName)
			if registryType == "gitlab" || registryType == "github" {
				p.print("\n  Set your token:\n")
				p.print("  export SCTL_REGISTRY_TOKEN=<your-token>\n")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing config")
	return cmd
}

// ── sctl init project ─────────────────────────────────────────────────────

func newInitProjectCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create agent-skills.lock in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := lockfile.DefaultFilename
			if !force {
				if _, err := os.Stat(path); err == nil {
					return &UserError{Message: fmt.Sprintf(
						"%s already exists\nUse --force to overwrite.", path,
					)}
				}
			}
			lf := lockfile.New()
			if err := lf.Write(path); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created %s\n", path)
			fmt.Fprintf(cmd.OutOrStdout(), "\n  Add skills with:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  sctl add <skill>[@version] --source <registry>\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing lockfile")
	return cmd
}

// ── sctl init skill ───────────────────────────────────────────────────────

func newInitSkillCmd() *cobra.Command {
	var outputDir string
	cmd := &cobra.Command{
		Use:   "skill <name>",
		Short: "Scaffold a new skill directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !isValidSkillName(name) {
				return &UserError{Message: fmt.Sprintf(
					"invalid skill name %q — use lowercase letters, digits, and hyphens only", name,
				)}
			}

			base := outputDir
			if base == "" {
				base = "."
			}
			skillDir := filepath.Join(base, name)

			if _, err := os.Stat(skillDir); err == nil {
				return &UserError{Message: fmt.Sprintf("directory %s already exists", skillDir)}
			}

			p := newPrompter(cmd)
			p.header("sctl skill init")

			description := p.ask("Description", "A specialized skill for "+name)
			version := p.ask("Initial version", "0.1.0")
			owner := p.ask("Owner (team or username)", "")
			platforms := p.multiChoose(
				"Compatible platforms (space-separated: claude-code gitlab-duo github-copilot codex all)",
				[]string{"claude-code", "gitlab-duo", "github-copilot", "codex"},
				[]string{"claude-code", "gitlab-duo"},
			)

			if err := os.MkdirAll(skillDir, 0o755); err != nil {
				return &InternalError{Message: "create skill dir", Cause: err}
			}

			data := skillScaffoldData{
				Name:        name,
				Description: description,
				Version:     version,
				Owner:       owner,
				Platforms:   platforms,
				Year:        time.Now().Year(),
			}

			files := map[string]string{
				"SKILL.md":    skillMDTemplate,
				"skill.yaml":  skillYAMLTemplate,
				"VERSION":     version,
				"CHANGELOG.md": skillChangelogTemplate,
			}

			for filename, tmplStr := range files {
				content, err := renderTemplate(tmplStr, data)
				if err != nil {
					return &InternalError{Message: "render " + filename, Cause: err}
				}
				if err := writeAtomic(filepath.Join(skillDir, filename), content); err != nil {
					return &InternalError{Message: "write " + filename, Cause: err}
				}
			}

			p.print("\n✓ Created %s/\n", skillDir)
			for f := range files {
				p.print("    %s/%s\n", name, f)
			}
			p.print("\nNext steps:\n")
			p.print("  1. Edit %s/SKILL.md with your skill's instructions\n", name)
			p.print("  2. sctl validate %s\n", skillDir)
			p.print("  3. sctl package  %s\n", skillDir)
			return nil
		},
	}
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Parent directory for the skill (default: current directory)")
	return cmd
}

// ── Prompter ─────────────────────────────────────────────────────────────

type prompter struct {
	cmd    *cobra.Command
	reader *bufio.Reader
}

func newPrompter(cmd *cobra.Command) *prompter {
	return &prompter{cmd: cmd, reader: bufio.NewReader(os.Stdin)}
}

func (p *prompter) print(format string, args ...interface{}) {
	fmt.Fprintf(p.cmd.OutOrStdout(), format, args...)
}

func (p *prompter) header(title string) {
	p.print("%s\n%s\n\n", title, strings.Repeat("─", len(title)))
}

func (p *prompter) ask(label, defaultVal string) string {
	if defaultVal != "" {
		p.print("  %s [%s]: ", label, defaultVal)
	} else {
		p.print("  %s: ", label)
	}
	line, _ := p.reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

func (p *prompter) choose(label string, options []string) string {
	p.print("  %s\n", label)
	for i, o := range options {
		p.print("    [%d] %s\n", i+1, o)
	}
	p.print("  Choice [1]: ")
	line, _ := p.reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" || line == "1" {
		return options[0]
	}
	for i, o := range options {
		if line == fmt.Sprintf("%d", i+1) || line == o {
			return o
		}
	}
	return options[0]
}

func (p *prompter) multiChoose(label string, _ []string, defaults []string) []string {
	p.print("  %s\n", label)
	p.print("  Default [%s]: ", strings.Join(defaults, " "))
	line, _ := p.reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaults
	}
	return strings.Fields(line)
}

// ── Helpers ───────────────────────────────────────────────────────────────

func configFilePath() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "sctl", "config.yaml"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sctl", "config.yaml"), nil
}

func writeAtomic(path, content string) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func isValidSkillName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

func buildConfigYAML(name, regType, url, project, cacheDir string) string {
	var sb strings.Builder
	sb.WriteString("default_registry: " + name + "\n\n")
	sb.WriteString("cache_dir: " + cacheDir + "\n")
	sb.WriteString("log_level: info\n\n")
	sb.WriteString("registries:\n")
	sb.WriteString("  " + name + ":\n")
	sb.WriteString("    type: " + regType + "\n")
	sb.WriteString("    url: " + url + "\n")
	if project != "" {
		sb.WriteString("    project: " + project + "\n")
	}
	sb.WriteString("    token: \"\"  # set via SCTL_REGISTRY_TOKEN\n")
	return sb.String()
}

func renderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

type skillScaffoldData struct {
	Name        string
	Description string
	Version     string
	Owner       string
	Platforms   []string
	Year        int
}

// ── Skill file templates ──────────────────────────────────────────────────

var skillMDTemplate = `---
name: {{.Name}}
description: {{.Description}}
---

# {{.Name}}

{{.Description}}

## When This Skill Activates

Describe the situations where this skill should be applied.

## Instructions

Provide step-by-step instructions for the agent here.

## Output Format

Describe what output the agent should produce.

## Examples

Provide a concrete input/output example here.
`

var skillYAMLTemplate = `name: {{.Name}}
version: "{{.Version}}"
description: {{.Description}}
{{- if .Owner}}
owners:
  - {{.Owner}}
{{- end}}
compatible_with:
{{- range .Platforms}}
  - {{.}}
{{- end}}
`

var skillChangelogTemplate = `# Changelog

## {{.Version}}

### Added
- Initial release of {{.Name}}.
`
