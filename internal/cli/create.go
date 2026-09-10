package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/sklib/spec"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var description string
	var license string
	var platforms []string
	var namespace string
	var noInteractive bool

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Scaffold a new publishable skill directory",
		Long: `Creates a new skill directory with all files required for publishing:

  SKILL.md       — human-readable capability description (template)
  skill.yaml     — machine-readable metadata (name, version, platforms)
  VERSION        — current version (0.1.0)
  CHANGELOG.md   — initial changelog entry

Run 'skpm publish' inside the directory when ready to release.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"scaffold", "skill"}, args...), func() error {
				name := args[0]
				if err := spec.ValidateSkillName(name); err != nil {
					return &UserError{Message: fmt.Sprintf("invalid skill name %q: use lowercase letters, digits, and hyphens", name)}
				}

				destDir := name
				if _, err := os.Stat(destDir); err == nil {
					return &UserError{Message: fmt.Sprintf("directory %q already exists", destDir)}
				}

				// Interactive prompts when running in a terminal and not suppressed.
				if !noInteractive && isInteractiveTerminal() {
					description = promptIfEmpty(cmd, "Description", description)
					license = promptIfEmpty(cmd, "License (SPDX)", license)
					if len(platforms) == 0 {
						raw := promptIfEmpty(cmd, "Platforms (comma-separated, default: all)", "")
						if raw != "" {
							for _, p := range strings.Split(raw, ",") {
								p = strings.TrimSpace(p)
								if p != "" {
									platforms = append(platforms, p)
								}
							}
						}
					}
				}

				if license == "" {
					license = "MIT"
				}
				if len(platforms) == 0 {
					platforms = []string{"all"}
				}
				if namespace == "" {
					namespace = "default"
				}

				if globalDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would create skill %q in %s/\n", name, destDir)
					fmt.Fprintf(cmd.OutOrStdout(), "  description: %s\n  license: %s\n  platforms: %s\n",
						description, license, strings.Join(platforms, ", "))
					return nil
				}

				if err := os.MkdirAll(destDir, 0o755); err != nil {
					return &InternalError{Message: "create directory", Cause: err}
				}

				files := scaffoldFiles(name, description, license, namespace, platforms)
				for filename, content := range files {
					path := filepath.Join(destDir, filename)
					if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
						os.RemoveAll(destDir)
						return &InternalError{Message: fmt.Sprintf("write %s", filename), Cause: err}
					}
				}

				abs, _ := filepath.Abs(destDir)
				fmt.Fprintf(cmd.OutOrStdout(), "Created skill %q in %s\n\n", name, abs)
				fmt.Fprintf(cmd.OutOrStdout(), "Next steps:\n")
				fmt.Fprintf(cmd.OutOrStdout(), "  1. Edit %s/SKILL.md with your skill's instructions\n", destDir)
				fmt.Fprintf(cmd.OutOrStdout(), "  2. skpm validate %s\n", destDir)
				fmt.Fprintf(cmd.OutOrStdout(), "  3. skpm publish %s\n", destDir)
				return nil
			})
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "One-line skill description")
	cmd.Flags().StringVar(&license, "license", "", "SPDX license identifier (default: MIT)")
	cmd.Flags().StringArrayVar(&platforms, "platform", nil, "Compatible platforms (default: all)")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Registry namespace (default: default)")
	cmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Skip interactive prompts, use flag values only")
	return cmd
}

func scaffoldFiles(name, description, license, namespace string, platforms []string) map[string]string {
	if description == "" {
		description = "TODO: describe what this skill enables the agent to do."
	}

	platformList := make([]string, len(platforms))
	for i, p := range platforms {
		platformList[i] = "  - " + p
	}
	platformYAML := strings.Join(platformList, "\n")

	skillYAML := fmt.Sprintf(`name: %s
version: 0.1.0
description: %s
namespace: %s
license: %s
compatible_with:
%s
`, name, description, namespace, license, platformYAML)

	skillMD := fmt.Sprintf(`---
name: %s
description: %s
license: %s
compatibility:
  platforms:
%s
metadata:
  version: 0.1.0
---

# %s

TODO: describe what this skill does and how to use it.

## Usage

TODO: provide usage examples.

## Examples

TODO: show concrete examples.
`, name, description, license, platformYAML, name)

	changelog := fmt.Sprintf(`# Changelog

## 0.1.0

- Initial release.
`)

	return map[string]string{
		"SKILL.md":     skillMD,
		"skill.yaml":   skillYAML,
		"VERSION":      "0.1.0\n",
		"CHANGELOG.md": changelog,
	}
}

func promptIfEmpty(cmd *cobra.Command, label, current string) string {
	if current != "" {
		return current
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s: ", label)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
