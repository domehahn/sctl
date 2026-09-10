package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var stableSemverPattern = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)

type skillVersionState struct {
	Path        string `json:"path"`
	Version     string `json:"version"`
	SkillYAML   string `json:"skill_yaml_version"`
	ChangelogOK bool   `json:"changelog_entry"`
}

func newVersionShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <path>",
		Short: "Show a skill's current version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := readSkillVersionState(args[0])
			if err != nil {
				return err
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "version show", Data: state})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", state.Version)
			return nil
		},
	}
	return cmd
}

func newVersionBumpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump a skill version",
	}
	cmd.AddCommand(newVersionBumpPartCmd("patch"))
	cmd.AddCommand(newVersionBumpPartCmd("minor"))
	cmd.AddCommand(newVersionBumpPartCmd("major"))
	return cmd
}

func newVersionBumpPartCmd(part string) *cobra.Command {
	return &cobra.Command{
		Use:   part + " <path>",
		Short: "Bump the " + part + " version component",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"version", "bump", part}, args...), func() error {
				state, err := readSkillVersionState(args[0])
				if err != nil {
					return err
				}
				next, err := bumpStableVersion(state.Version, part)
				if err != nil {
					return err
				}
				return writeSkillVersion(cmd, args[0], state.Version, next, "version bump "+part)
			})
		},
	}
}

func newVersionSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <version> <path>",
		Short: "Set a skill version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"version", "set"}, args...), func() error {
				next, err := normalizeStableVersion(args[0])
				if err != nil {
					return err
				}
				state, err := readSkillVersionState(args[1])
				if err != nil {
					return err
				}
				return writeSkillVersion(cmd, args[1], state.Version, next, "version set")
			})
		},
	}
}

func readSkillVersionState(dir string) (*skillVersionState, error) {
	version, err := readVersionFile(dir)
	if err != nil {
		return nil, err
	}
	skillVersion, err := readSkillYAMLVersion(dir)
	if err != nil {
		return nil, err
	}
	if version != skillVersion {
		return nil, &UserError{Message: fmt.Sprintf("VERSION %q does not match skill.yaml.version %q", version, skillVersion)}
	}
	return &skillVersionState{
		Path:        dir,
		Version:     version,
		SkillYAML:   skillVersion,
		ChangelogOK: changelogHasVersion(dir, version),
	}, nil
}

func writeSkillVersion(cmd *cobra.Command, dir, previous, next, command string) error {
	if globalDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would update %s from %s to %s\n", dir, previous, next)
		return nil
	}
	if err := writeAtomic(filepath.Join(dir, "VERSION"), next+"\n"); err != nil {
		return &InternalError{Message: "write VERSION", Cause: err}
	}
	if err := updateSkillYAMLVersion(dir, next); err != nil {
		return err
	}
	if err := ensureChangelogVersion(dir, next); err != nil {
		return &InternalError{Message: "update CHANGELOG.md", Cause: err}
	}
	data := map[string]string{"path": dir, "previous": previous, "version": next}
	if outputFormat() == OutputJSON {
		PrintResult(OutputJSON, CommandResult{Success: true, Command: command, Data: data})
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Updated %s: %s -> %s\n", dir, previous, next)
	return nil
}

func readVersionFile(dir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "VERSION"))
	if err != nil {
		return "", &UserError{Message: fmt.Sprintf("read VERSION: %v", err)}
	}
	return normalizeStoredStableVersion("VERSION", strings.TrimSpace(string(data)))
}

func readSkillYAMLVersion(dir string) (string, error) {
	node, err := readSkillYAMLNode(dir)
	if err != nil {
		return "", err
	}
	version, ok := findYAMLScalar(node, "version")
	if !ok || strings.TrimSpace(version) == "" {
		return "", &UserError{Message: "skill.yaml.version is required"}
	}
	return normalizeStoredStableVersion("skill.yaml.version", version)
}

func updateSkillYAMLVersion(dir, version string) error {
	node, err := readSkillYAMLNode(dir)
	if err != nil {
		return err
	}
	if !setYAMLScalar(node, "version", version) {
		return &UserError{Message: "skill.yaml.version is required"}
	}
	data, err := yaml.Marshal(node)
	if err != nil {
		return &InternalError{Message: "marshal skill.yaml", Cause: err}
	}
	if err := writeAtomic(filepath.Join(dir, "skill.yaml"), string(data)); err != nil {
		return &InternalError{Message: "write skill.yaml", Cause: err}
	}
	return nil
}

func readSkillYAMLNode(dir string) (*yaml.Node, error) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("read skill.yaml: %v", err)}
	}
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse skill.yaml: %v", err)}
	}
	return &node, nil
}

func findYAMLScalar(node *yaml.Node, key string) (string, bool) {
	mapping := yamlMapping(node)
	if mapping == nil {
		return "", false
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1].Value, true
		}
	}
	return "", false
}

func setYAMLScalar(node *yaml.Node, key, value string) bool {
	mapping := yamlMapping(node)
	if mapping == nil {
		return false
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1].Kind = yaml.ScalarNode
			mapping.Content[i+1].Tag = "!!str"
			mapping.Content[i+1].Value = value
			return true
		}
	}
	return false
}

func yamlMapping(node *yaml.Node) *yaml.Node {
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = node.Content[0]
	}
	if node.Kind != yaml.MappingNode {
		return nil
	}
	return node
}

func normalizeStableVersion(version string) (string, error) {
	matches := stableSemverPattern.FindStringSubmatch(strings.TrimSpace(version))
	if matches == nil {
		return "", &UserError{Message: fmt.Sprintf("%q is not a stable SemVer version (expected 1.2.3)", version)}
	}
	return fmt.Sprintf("%d.%d.%d", mustAtoi(matches[1]), mustAtoi(matches[2]), mustAtoi(matches[3])), nil
}

func normalizeStoredStableVersion(field, version string) (string, error) {
	normalized, err := normalizeStableVersion(version)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(version) != normalized {
		return "", &UserError{Message: fmt.Sprintf("%s must be stored without a leading v: %s", field, normalized)}
	}
	return normalized, nil
}

func bumpStableVersion(version, part string) (string, error) {
	normalized, err := normalizeStableVersion(version)
	if err != nil {
		return "", err
	}
	parts := strings.Split(normalized, ".")
	major, minor, patch := mustAtoi(parts[0]), mustAtoi(parts[1]), mustAtoi(parts[2])
	switch part {
	case "patch":
		patch++
	case "minor":
		minor++
		patch = 0
	case "major":
		major++
		minor = 0
		patch = 0
	default:
		return "", &UserError{Message: fmt.Sprintf("unknown version part %q", part)}
	}
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func changelogHasVersion(dir, version string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		return false
	}
	re := regexp.MustCompile(`(?m)^##\s+v?` + regexp.QuoteMeta(version) + `\s*$`)
	return re.Match(data)
}

func ensureChangelogVersion(dir, version string) error {
	path := filepath.Join(dir, "CHANGELOG.md")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return writeAtomic(path, "# Changelog\n\n## "+version+"\n\n- TODO: describe changes.\n")
		}
		return err
	}
	content := string(data)
	if changelogHasVersion(dir, version) {
		return nil
	}
	section := "## " + version + "\n\n- TODO: describe changes.\n\n"
	if strings.HasPrefix(content, "# Changelog") {
		lines := strings.SplitN(content, "\n", 2)
		rest := ""
		if len(lines) == 2 {
			rest = strings.TrimLeft(lines[1], "\n")
		}
		content = lines[0] + "\n\n" + section + rest
	} else {
		content = "# Changelog\n\n" + section + strings.TrimLeft(content, "\n")
	}
	return writeAtomic(path, content)
}
