package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const workspaceFile = "skpm-workspace.yaml"

type workspaceManifest struct {
	Version  int               `yaml:"version"`
	Skills   []string          `yaml:"skills"`
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

func newWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage a monorepo workspace of multiple skills",
		Long: `Workspace commands operate across multiple skill directories defined in
skpm-workspace.yaml at the repository root.

  skpm workspace init              — create skpm-workspace.yaml
  skpm workspace list              — list all skills with their versions
  skpm workspace run <script>      — run a script in all (or selected) skills
  skpm workspace validate          — validate all skills
  skpm workspace publish           — publish all (or changed) skills
  skpm workspace graph             — show skill dependency relationships

Skills are declared as relative paths:

  version: 1
  skills:
    - ./skill-a
    - ./skill-b
    - ./shared/review-skill`,
	}
	cmd.AddCommand(newWorkspaceInitCmd())
	cmd.AddCommand(newWorkspaceListCmd())
	cmd.AddCommand(newWorkspaceRunCmd())
	cmd.AddCommand(newWorkspaceValidateCmd())
	cmd.AddCommand(newWorkspacePublishCmd())
	cmd.AddCommand(newWorkspaceGraphCmd())
	return cmd
}

// ── init ─────────────────────────────────────────────────────────────────────

func newWorkspaceInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [skill-dirs...]",
		Short: "Create skpm-workspace.yaml",
		Long: `Creates a skpm-workspace.yaml in the current directory.

Pass skill directories as arguments or let skpm discover them automatically
(looks for skill.yaml in immediate subdirectories).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(workspaceFile); err == nil {
				return &UserError{Message: workspaceFile + " already exists"}
			}

			var paths []string
			if len(args) > 0 {
				paths = args
			} else {
				entries, err := os.ReadDir(".")
				if err != nil {
					return &InternalError{Message: "read directory", Cause: err}
				}
				for _, e := range entries {
					if !e.IsDir() {
						continue
					}
					if _, err := os.Stat(filepath.Join(e.Name(), "skill.yaml")); err == nil {
						paths = append(paths, "./"+e.Name())
					}
				}
			}

			if len(paths) == 0 {
				return &UserError{Message: "no skill directories found; pass them as arguments or add skill.yaml files in subdirectories"}
			}

			wm := workspaceManifest{Version: 1, Skills: paths}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would create %s with %d skills\n", workspaceFile, len(paths))
				for _, p := range paths {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", p)
				}
				return nil
			}

			data, err := yaml.Marshal(&wm)
			if err != nil {
				return &InternalError{Message: "marshal workspace file", Cause: err}
			}
			if err := writeAtomic(workspaceFile, string(data)); err != nil {
				return &InternalError{Message: "write " + workspaceFile, Cause: err}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created %s with %d skill(s):\n", workspaceFile, len(paths))
			for _, p := range paths {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", p)
			}
			return nil
		},
	}
}

// ── list ─────────────────────────────────────────────────────────────────────

func newWorkspaceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all skills in the workspace with their versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			wm, err := readWorkspace()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Workspace (%d skills):\n", len(wm.Skills))
			v := skill.NewValidator()
			for _, path := range wm.Skills {
				sy, syErr := readSkillYAML(path)
				if syErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✗ %-30s  (could not read skill.yaml)\n", path)
					continue
				}
				res, valErr := v.Validate(context.Background(), path)
				status := "✓"
				if valErr != nil || !res.Valid {
					status = "~"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s %-30s  %s@%s\n", status, path, sy.Name, sy.Version)
			}
			return nil
		},
	}
}

// ── run ───────────────────────────────────────────────────────────────────────

func newWorkspaceRunCmd() *cobra.Command {
	var only []string
	var failFast bool

	cmd := &cobra.Command{
		Use:   "run <script>",
		Short: "Run a script in all (or selected) workspace skills",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptName := args[0]
			dirs, err := workspaceDirs(only)
			if err != nil {
				return err
			}

			var failed []string
			for _, dir := range dirs {
				scripts, err := readSkillScripts(dir)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "  [%s] skipping: %v\n", dir, err)
					continue
				}
				script, ok := scripts[scriptName]
				if !ok {
					fmt.Fprintf(cmd.OutOrStdout(), "  [%s] script %q not defined — skipping\n", dir, scriptName)
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\n  [%s] > %s\n", dir, script)
				if globalDryRun {
					continue
				}
				sh, flag := shellInterpreter()
				c := exec.Command(sh, flag, script) //nolint:gosec
				c.Dir = dir
				c.Stdout = cmd.OutOrStdout()
				c.Stderr = cmd.ErrOrStderr()
				c.Env = append(os.Environ(), "SKPM_WORKSPACE=1", "SKPM_SCRIPT="+scriptName)
				if runErr := c.Run(); runErr != nil {
					failed = append(failed, dir)
					fmt.Fprintf(cmd.ErrOrStderr(), "  [%s] ✗ script failed: %v\n", dir, runErr)
					if failFast {
						return &UserError{Message: fmt.Sprintf("script failed in %s", dir)}
					}
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "  [%s] ✓\n", dir)
				}
			}
			if len(failed) > 0 {
				return &UserError{Message: fmt.Sprintf("script %q failed in: %s", scriptName, strings.Join(failed, ", "))}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&only, "skill", nil, "Only run in these skills (by path, repeatable)")
	cmd.Flags().BoolVar(&failFast, "fail-fast", false, "Stop after the first failure")
	return cmd
}

// ── validate ─────────────────────────────────────────────────────────────────

func newWorkspaceValidateCmd() *cobra.Command {
	var only []string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate all (or selected) skills in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs, err := workspaceDirs(only)
			if err != nil {
				return err
			}
			v := skill.NewValidatorWithOptions(skill.ValidationOptions{Publish: false})
			allValid := true
			for _, dir := range dirs {
				res, err := v.Validate(cmd.Context(), dir)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "  ✗ %s: error: %v\n", dir, err)
					allValid = false
					continue
				}
				if res.Valid {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s\n", dir)
				} else {
					allValid = false
					fmt.Fprintf(cmd.OutOrStdout(), "  ✗ %s (%d error(s))\n", dir, len(res.Errors))
					for _, e := range res.Errors {
						fmt.Fprintf(cmd.OutOrStdout(), "      [%s] %s\n", e.Field, e.Message)
					}
				}
			}
			if !allValid {
				return &UserError{Message: "one or more skills failed validation"}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&only, "skill", nil, "Only validate these skills (by path, repeatable)")
	return cmd
}

// ── publish ───────────────────────────────────────────────────────────────────

func newWorkspacePublishCmd() *cobra.Command {
	var (
		only    []string
		changed bool
		source  string
		noTag   bool
		noPush  bool
	)

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish all (or changed) skills in the workspace",
		Long:  `Runs 'skpm publish' for each skill. Use --changed to publish only skills with uncommitted git changes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs, err := workspaceDirs(only)
			if err != nil {
				return err
			}
			if changed {
				dirs, err = filterChangedDirs(dirs)
				if err != nil {
					return err
				}
				if len(dirs) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "No changed skills detected.")
					return nil
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Changed skills (%d):\n", len(dirs))
				for _, d := range dirs {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", d)
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}

			skpmBin, err := os.Executable()
			if err != nil {
				skpmBin = "skpm"
			}

			var failed []string
			for _, dir := range dirs {
				fmt.Fprintf(cmd.OutOrStdout(), "Publishing %s...\n", dir)
				if globalDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "  (dry run)\n")
					continue
				}
				pubArgs := []string{"publish", dir}
				if source != "" {
					pubArgs = append(pubArgs, "--source", source)
				}
				if noTag {
					pubArgs = append(pubArgs, "--no-tag")
				}
				if noPush {
					pubArgs = append(pubArgs, "--no-push")
				}
				c := exec.Command(skpmBin, pubArgs...) //nolint:gosec
				c.Stdout = cmd.OutOrStdout()
				c.Stderr = cmd.ErrOrStderr()
				if runErr := c.Run(); runErr != nil {
					failed = append(failed, dir)
					fmt.Fprintf(cmd.ErrOrStderr(), "  ✗ publish failed for %s: %v\n", dir, runErr)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✓ published %s\n", dir)
				}
			}
			if len(failed) > 0 {
				return &UserError{Message: fmt.Sprintf("publish failed for %d skill(s): %s",
					len(failed), strings.Join(failed, ", "))}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&only, "skill", nil, "Only publish these skills (by path, repeatable)")
	cmd.Flags().BoolVar(&changed, "changed", false, "Only publish skills with uncommitted git changes")
	cmd.Flags().StringVar(&source, "source", "", "Registry to publish to")
	cmd.Flags().BoolVar(&noTag, "no-tag", false, "Skip git tag")
	cmd.Flags().BoolVar(&noPush, "no-push", false, "Skip git push")
	return cmd
}

// ── graph ─────────────────────────────────────────────────────────────────────

func newWorkspaceGraphCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Show skill dependency relationships across the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			wm, err := readWorkspace()
			if err != nil {
				return err
			}

			type wsNode struct {
				name     string
				version  string
				path     string
				requires []string
			}
			var nodes []wsNode
			nameToPath := map[string]string{}

			for _, path := range wm.Skills {
				sy, err := readSkillYAML(path)
				if err != nil {
					continue
				}
				nodes = append(nodes, wsNode{
					name:     sy.Name,
					version:  sy.Version,
					path:     path,
					requires: sy.Requires,
				})
				nameToPath[sy.Name] = path
			}
			sort.Slice(nodes, func(i, j int) bool { return nodes[i].name < nodes[j].name })

			switch strings.ToLower(format) {
			case "dot":
				fmt.Fprintln(cmd.OutOrStdout(), "digraph workspace {")
				fmt.Fprintln(cmd.OutOrStdout(), `  graph [rankdir=LR];`)
				fmt.Fprintln(cmd.OutOrStdout(), `  node [shape=box];`)
				for _, n := range nodes {
					label := n.name + "@" + n.version
					fmt.Fprintf(cmd.OutOrStdout(), "  %q [label=%q];\n", n.name, label)
					for _, dep := range n.requires {
						fmt.Fprintf(cmd.OutOrStdout(), "  %q -> %q;\n", n.name, dep)
					}
				}
				fmt.Fprintln(cmd.OutOrStdout(), "}")
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "workspace (%d skills)\n", len(nodes))
				for i, n := range nodes {
					isLast := i == len(nodes)-1
					branch, child := "├── ", "│   "
					if isLast {
						branch, child = "└── ", "    "
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s%s@%s  (%s)\n", branch, n.name, n.version, n.path)
					for j, dep := range n.requires {
						isLastDep := j == len(n.requires)-1
						depBranch := child + "├── "
						if isLastDep {
							depBranch = child + "└── "
						}
						loc := "(external)"
						if p, ok := nameToPath[dep]; ok {
							loc = p
						}
						fmt.Fprintf(cmd.OutOrStdout(), "%srequires: %s  %s\n", depBranch, dep, loc)
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or dot")
	return cmd
}

// ── helpers ───────────────────────────────────────────────────────────────────

func readWorkspace() (*workspaceManifest, error) {
	wsPath := findWorkspaceFile()
	if wsPath == "" {
		return nil, &UserError{Message: workspaceFile + " not found — run 'skpm workspace init' first"}
	}
	data, err := os.ReadFile(wsPath)
	if err != nil {
		return nil, &InternalError{Message: "read " + workspaceFile, Cause: err}
	}
	var wm workspaceManifest
	if err := yaml.Unmarshal(data, &wm); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse %s: %v", workspaceFile, err)}
	}
	if len(wm.Skills) == 0 {
		return nil, &UserError{Message: workspaceFile + " has no skills defined"}
	}
	return &wm, nil
}

// findWorkspaceFile walks up from cwd looking for skpm-workspace.yaml.
func findWorkspaceFile() string {
	dir, _ := os.Getwd()
	for {
		candidate := filepath.Join(dir, workspaceFile)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func workspaceDirs(only []string) ([]string, error) {
	wm, err := readWorkspace()
	if err != nil {
		return nil, err
	}
	if len(only) == 0 {
		return wm.Skills, nil
	}
	set := map[string]bool{}
	for _, o := range only {
		set[o] = true
	}
	var out []string
	for _, s := range wm.Skills {
		if set[s] {
			out = append(out, s)
		}
	}
	return out, nil
}

// filterChangedDirs returns only the dirs that have uncommitted changes in git.
func filterChangedDirs(dirs []string) ([]string, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return dirs, nil // not a git repo — return all
	}

	changedFiles := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) < 4 {
			continue
		}
		changedFiles[strings.TrimSpace(line[3:])] = true
	}

	cwd, _ := os.Getwd()
	var changed []string
	for _, dir := range dirs {
		absDir, _ := filepath.Abs(dir)
		rel, _ := filepath.Rel(cwd, absDir)
		for f := range changedFiles {
			if strings.HasPrefix(filepath.ToSlash(f), filepath.ToSlash(rel)+"/") {
				changed = append(changed, dir)
				break
			}
		}
	}
	return changed, nil
}
