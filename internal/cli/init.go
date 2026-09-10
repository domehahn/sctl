package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const gitignoreBlock = `
# skpm — installed skill directories are generated artifacts
# restore with: skpm install
.claude/skills/
skills/
.agents/skills/
.github/skills/

# skpm — these files must be committed
# !agent-skills.yaml
# !agent-skills.lock
`

func newInitCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize skpm: set up config and project files interactively",
		Long: `Interactive setup wizard for skpm.

Runs two steps in sequence:

  1. Config  — creates ~/.config/skpm/config.yaml (skipped if already exists)
  2. Project — creates agent-skills.yaml, agent-skills.lock, updates .gitignore

To scaffold a new skill instead, use:
  skpm init skill <name>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := newPrompter(cmd)
			p.header("skpm init")

			// ── Step 1: Config ──────────────────────────────────────────
			cfgPath, err := configFilePath()
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}

			cfgExists := fileExists(cfgPath)
			if cfgExists && !force {
				p.print("✓ Config already exists at %s\n", cfgPath)
			} else {
				p.print("Step 1/2 — Config (%s)\n\n", cfgPath)

				registryType := p.choose("Registry type", []string{"gitlab", "github", "artifactory", "local"})
				registryName := p.ask("Registry name", "company-"+registryType)

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

				token := p.secret("Token (leave empty to set via SKPM_REGISTRY_TOKEN later)")
				cacheDir := p.ask("Cache directory", "~/.cache/skpm")

				cfg := buildConfigYAML(registryName, registryType, registryURL, registryProject, token, cacheDir)
				if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
					return &InternalError{Message: "create config dir", Cause: err}
				}
				if err := writeAtomic(cfgPath, cfg); err != nil {
					return &InternalError{Message: "write config", Cause: err}
				}
				p.print("\n✓ Created %s\n", cfgPath)
				if token == "" {
					p.print("  Token not set — export when needed: export SKPM_REGISTRY_TOKEN=<token>\n")
				}
			}

			// ── Step 2: Project ─────────────────────────────────────────
			p.print("\n")
			manifestExists := fileExists(manifest.DefaultFilename)
			if manifestExists && !force {
				p.print("✓ Project already initialized (%s exists)\n", manifest.DefaultFilename)
			} else {
				p.print("Step 2/2 — Project\n\n")

				if err := lockfile.New().Write(lockfile.DefaultFilename); err != nil {
					return &InternalError{Message: "write lockfile", Cause: err}
				}
				p.print("✓ Created %s\n", lockfile.DefaultFilename)

				if err := manifest.New().Write(manifest.DefaultFilename); err != nil {
					return &InternalError{Message: "write manifest", Cause: err}
				}
				p.print("✓ Created %s\n", manifest.DefaultFilename)

				updated, err := updateGitignore(".gitignore")
				if err != nil {
					p.print("  warning: could not update .gitignore: %v\n", err)
				} else if updated {
					p.print("✓ Updated .gitignore\n")
				} else {
					p.print("✓ .gitignore already up to date\n")
				}
			}

			// ── Summary ─────────────────────────────────────────────────
			p.print("\nAll done. Next steps:\n")
			p.print("  skpm add <skill>[@version] --source <registry>\n")
			p.print("  git add agent-skills.yaml agent-skills.lock\n")
			p.print("  git commit -m \"add agent skills\"\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Re-run even if config/project files already exist")
	cmd.AddCommand(newInitSkillCmd())
	return cmd
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// updateGitignore appends the skpm block to .gitignore if not already present.
// Returns true if the file was modified.
func updateGitignore(path string) (bool, error) {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if strings.Contains(string(existing), "skpm — installed skill directories") {
		return false, nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// ensure we start on a new line
	if len(existing) > 0 && existing[len(existing)-1] != '\n' {
		f.WriteString("\n")
	}
	_, err = f.WriteString(gitignoreBlock)
	return err == nil, err
}

// ── skpm init skill ───────────────────────────────────────────────────────

func newInitSkillCmd() *cobra.Command {
	var outputDir string
	cmd := &cobra.Command{
		Use:   "skill <name>",
		Short: "Scaffold a new skill directory (compatibility wrapper)",
		Long: `Scaffold a new skill directory.

Compatibility note:
For new workflows, prefer ` + "`skcr scaffold skill <name>`" + `.
Use ` + "`skpm`" + ` for validation, versioning, packaging, publishing, and installation.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"scaffold", "skill"}, args...), func() error {
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
				p.header("skpm skill init")
				p.print("Note: `skpm init skill` is kept for compatibility.\n")
				p.print("For new workflows, prefer `skcr scaffold skill <name>`.\n\n")

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
					"SKILL.md":        skillMDTemplate,
					"skill.yaml":      skillYAMLTemplate,
					"VERSION":         version + "\n",
					"CHANGELOG.md":    skillChangelogTemplate,
					"README.md":       skillReadmeTemplate,
					"LICENSE":         skillLicenseTemplate,
					"tests/README.md": skillTestsReadmeTemplate,
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
				p.print("  2. skpm validate %s\n", skillDir)
				p.print("  3. skpm version bump patch %s\n", skillDir)
				p.print("  4. skpm package %s\n", skillDir)
				return nil
			})
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

func (p *prompter) secret(label string) string {
	p.print("  %s: ", label)
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		p.print("\n")
		if err != nil || len(b) == 0 {
			return ""
		}
		return string(b)
	}
	// non-interactive fallback
	line, _ := p.reader.ReadString('\n')
	return strings.TrimSpace(line)
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
		return filepath.Join(xdg, "skpm", "config.yaml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "skpm", "config.yaml"), nil
}

func writeAtomic(path, content string) error {
	return writeAtomicWithMode(path, content, 0o644)
}

func writeAtomicWithMode(path, content string, mode os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(tmp, []byte(content), mode); err != nil {
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

func buildConfigYAML(name, regType, url, project, token, cacheDir string) string {
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
	if token != "" {
		sb.WriteString("    token: " + token + "\n")
	} else {
		sb.WriteString("    token: \"\"  # set via SKPM_REGISTRY_TOKEN\n")
	}
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

var skillReadmeTemplate = `# {{.Name}}

{{.Description}}

## Lifecycle

Use skpm to validate, version, package, and publish this skill.
`

var skillLicenseTemplate = `Copyright (c) {{.Year}} {{if .Owner}}{{.Owner}}{{else}}{{.Name}} maintainers{{end}}

All rights reserved.
`

var skillTestsReadmeTemplate = `# Tests

Add skill fixtures, examples, and validation notes here.
`
