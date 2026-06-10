package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const linksFile = ".skpm-links.yaml"

type skillLink struct {
	Path        string   `yaml:"path"`
	InstalledTo []string `yaml:"installed_to"`
}

type linksManifest struct {
	Links map[string]skillLink `yaml:"links"`
}

func newLinkCmd() *cobra.Command {
	var platforms []string

	cmd := &cobra.Command{
		Use:   "link [path]",
		Short: "Link a local skill directory into the project for development",
		Long: `Creates symlinks from the project's skill install paths to a local skill
directory so that edits are reflected immediately without republishing.

The skill name and compatible platforms are read from skill.yaml in the
target directory. Use --platform to override the platform list.

Run 'skpm unlink <name>' to remove the link.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcDir := "."
			if len(args) == 1 {
				srcDir = args[0]
			}
			absSrc, err := filepath.Abs(srcDir)
			if err != nil {
				return err
			}

			sy, err := readSkillYAML(absSrc)
			if err != nil {
				return err
			}
			if sy.Name == "" {
				return &UserError{Message: "skill.yaml must declare a name"}
			}

			plats := sy.CompatibleWith
			if len(platforms) > 0 {
				plats = make([]skill.Platform, len(platforms))
				for i, p := range platforms {
					plats[i] = skill.Platform(p)
				}
			}
			if len(plats) == 0 {
				plats = []skill.Platform{skill.PlatformAll}
			}

			paths, err := installer.ResolvePaths(sy.Name, plats)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve paths: %v", err)}
			}

			workDir, _ := os.Getwd()

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would link %q → %s\n", sy.Name, absSrc)
				for _, p := range paths {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", p)
				}
				return nil
			}

			for _, p := range paths {
				dest := filepath.Join(workDir, p)
				// Remove any existing directory or symlink at the destination.
				if _, err := os.Lstat(dest); err == nil {
					if err := os.RemoveAll(dest); err != nil {
						return &InternalError{Message: fmt.Sprintf("remove existing %s", dest), Cause: err}
					}
				}
				if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
					return &InternalError{Message: "create parent dir", Cause: err}
				}
				if err := os.Symlink(absSrc, dest); err != nil {
					return &InternalError{Message: fmt.Sprintf("symlink %s", dest), Cause: err}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "linked  %s  →  %s\n", p, absSrc)
			}

			if err := upsertLink(sy.Name, absSrc, paths); err != nil {
				return &InternalError{Message: "update " + linksFile, Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n%q linked. Run 'skpm unlink %s' to remove.\n", sy.Name, sy.Name)
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&platforms, "platform", nil, "Override compatible platforms (e.g. claude-code)")
	return cmd
}

func newUnlinkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlink <name>",
		Short: "Remove a local development link created by 'skpm link'",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			lm, err := readLinksManifest()
			if err != nil {
				return err
			}

			lk, ok := lm.Links[name]
			if !ok {
				return &UserError{Message: fmt.Sprintf("%q is not linked (see %s)", name, linksFile)}
			}

			workDir, _ := os.Getwd()

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would unlink %q (%s)\n", name, lk.Path)
				return nil
			}

			for _, p := range lk.InstalledTo {
				dest := filepath.Join(workDir, p)
				info, err := os.Lstat(dest)
				if err != nil {
					continue
				}
				if info.Mode()&os.ModeSymlink == 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "skip  %s  (not a symlink — remove manually if needed)\n", p)
					continue
				}
				if err := os.Remove(dest); err != nil {
					return &InternalError{Message: fmt.Sprintf("remove symlink %s", dest), Cause: err}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "unlinked  %s\n", p)
			}

			delete(lm.Links, name)
			if err := writeLinksManifest(lm); err != nil {
				return &InternalError{Message: "update " + linksFile, Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n%q unlinked. Run 'skpm install' to restore from registry.\n", name)
			return nil
		},
	}
}

// ── links manifest helpers ────────────────────────────────────────────────

func readLinksManifest() (*linksManifest, error) {
	data, err := os.ReadFile(linksFile)
	if os.IsNotExist(err) {
		return &linksManifest{Links: map[string]skillLink{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var lm linksManifest
	if err := yaml.Unmarshal(data, &lm); err != nil {
		return nil, err
	}
	if lm.Links == nil {
		lm.Links = map[string]skillLink{}
	}
	return &lm, nil
}

func writeLinksManifest(lm *linksManifest) error {
	if len(lm.Links) == 0 {
		return os.Remove(linksFile)
	}
	data, err := yaml.Marshal(lm)
	if err != nil {
		return err
	}
	return os.WriteFile(linksFile, data, 0o644)
}

func upsertLink(name, absPath string, installedTo []string) error {
	lm, err := readLinksManifest()
	if err != nil {
		return err
	}
	lm.Links[name] = skillLink{Path: absPath, InstalledTo: installedTo}
	return writeLinksManifest(lm)
}

// isLinkedSkill returns (true, path) if name is recorded in .skpm-links.yaml.
func isLinkedSkill(name string) (bool, string) {
	lm, err := readLinksManifest()
	if err != nil {
		return false, ""
	}
	lk, ok := lm.Links[name]
	return ok, lk.Path
}

// readSkillYAML reads and unmarshals skill.yaml from a directory.
func readSkillYAML(dir string) (*skill.SkillYAML, error) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("read %s/skill.yaml: %v", dir, err)}
	}
	var sy skill.SkillYAML
	if err := yaml.Unmarshal(data, &sy); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse %s/skill.yaml: %v", dir, err)}
	}
	return &sy, nil
}
