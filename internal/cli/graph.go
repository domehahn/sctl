package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

type graphNode struct {
	name     string
	version  string
	requires []string
}

func newGraphCmd() *cobra.Command {
	var (
		lockPath     string
		manifestPath string
		format       string
		skillDir     string
	)

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Show the skill dependency and installation graph",
		Long: `Visualises two complementary views of installed skills:

  Project view (default) — reads agent-skills.yaml + agent-skills.lock and shows
  which skills are installed and into which platform directories.

  Skill view (--skill-dir) — reads a single skill.yaml 'requires' list and
  renders that skill's declared dependencies.

Output formats:
  text (default)  human-readable tree
  dot             Graphviz DOT for rendering with: dot -Tsvg graph.dot > graph.svg`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if skillDir != "" {
				return runSkillGraph(cmd, skillDir, format)
			}
			return runProjectGraph(cmd, manifestPath, lockPath, format)
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", lockfile.DefaultFilename, "Path to agent-skills.lock")
	cmd.Flags().StringVar(&manifestPath, "manifest", manifest.DefaultFilename, "Path to agent-skills.yaml")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or dot")
	cmd.Flags().StringVar(&skillDir, "skill-dir", "", "Show the dependency graph for a single skill.yaml")
	return cmd
}

// runProjectGraph shows all locked skills and their installation paths.
func runProjectGraph(cmd *cobra.Command, manifestPath, lockPath, format string) error {
	mf, mfErr := manifest.Read(manifestPath)
	lf, lfErr := lockfile.Read(lockPath)

	if mfErr != nil && lfErr != nil {
		return &UserError{Message: "no agent-skills.yaml or agent-skills.lock found in current directory\nRun 'skpm init' or 'skpm add <skill>' first."}
	}

	type projectNode struct {
		name       string
		version    string
		constraint string
		paths      []string
		linked     bool
	}

	constraints := map[string]string{}
	if mfErr == nil {
		for _, s := range mf.Skills {
			constraints[s.Name] = s.Version
		}
	}

	var nodes []projectNode
	if lfErr == nil {
		for _, s := range lf.Skills {
			linkedOK, _ := isLinkedSkill(s.Name)
			nodes = append(nodes, projectNode{
				name:       s.Name,
				version:    s.Version,
				constraint: constraints[s.Name],
				paths:      s.InstalledTo,
				linked:     linkedOK,
			})
		}
	} else {
		for _, s := range mf.Skills {
			nodes = append(nodes, projectNode{
				name:       s.Name,
				constraint: s.Version,
			})
		}
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].name < nodes[j].name })

	switch strings.ToLower(format) {
	case "dot":
		fmt.Fprintln(cmd.OutOrStdout(), "digraph skills {")
		fmt.Fprintln(cmd.OutOrStdout(), `  graph [rankdir=LR];`)
		fmt.Fprintln(cmd.OutOrStdout(), `  node  [shape=box];`)
		fmt.Fprintln(cmd.OutOrStdout(), `  "project" [shape=ellipse];`)
		for _, n := range nodes {
			label := n.name
			if n.version != "" {
				label += "@" + n.version
			}
			if n.linked {
				label += " (linked)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %q [label=%q];\n", n.name, label)
			fmt.Fprintf(cmd.OutOrStdout(), "  \"project\" -> %q;\n", n.name)
			for _, p := range n.paths {
				fmt.Fprintf(cmd.OutOrStdout(), "  %q -> %q [style=dashed, label=\"installed\"];\n", n.name, p)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), "}")

	default:
		if len(nodes) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No skills found.")
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), "project")
		for i, n := range nodes {
			isLast := i == len(nodes)-1
			branch, child := "├── ", "│   "
			if isLast {
				branch, child = "└── ", "    "
			}

			versionStr := ""
			if n.version != "" {
				versionStr = "@" + n.version
			} else if n.constraint != "" {
				versionStr = " (" + n.constraint + ")"
			}
			linkedStr := ""
			if n.linked {
				linkedStr = " [linked]"
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s%s%s%s\n", branch, n.name, versionStr, linkedStr)

			for j, p := range n.paths {
				isLastPath := j == len(n.paths)-1
				pathBranch := child + "├── "
				if isLastPath {
					pathBranch = child + "└── "
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", pathBranch, p)
			}
		}
	}
	return nil
}

// runSkillGraph shows the requires tree for a single skill.yaml.
func runSkillGraph(cmd *cobra.Command, dir, format string) error {
	root, err := readGraphNode(dir)
	if err != nil {
		return err
	}

	// BFS to resolve requires recursively from installed skill dirs (max depth 5).
	all := map[string]graphNode{root.name: root}
	type entry struct {
		node  graphNode
		depth int
	}
	queue := []entry{{node: root, depth: 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.depth >= 5 {
			continue
		}
		for _, dep := range cur.node.requires {
			if _, seen := all[dep]; seen {
				continue
			}
			depNode, depErr := resolveInstalledGraphNode(dep)
			if depErr != nil {
				depNode = graphNode{name: dep, version: "?"}
			}
			all[dep] = depNode
			queue = append(queue, entry{node: depNode, depth: cur.depth + 1})
		}
	}

	switch strings.ToLower(format) {
	case "dot":
		fmt.Fprintln(cmd.OutOrStdout(), "digraph skill_deps {")
		fmt.Fprintln(cmd.OutOrStdout(), `  graph [rankdir=LR];`)
		fmt.Fprintln(cmd.OutOrStdout(), `  node [shape=box];`)
		for _, n := range all {
			label := n.name
			if n.version != "" && n.version != "?" {
				label += "@" + n.version
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %q [label=%q];\n", n.name, label)
			for _, dep := range n.requires {
				fmt.Fprintf(cmd.OutOrStdout(), "  %q -> %q;\n", n.name, dep)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), "}")

	default:
		printGraphTree(cmd, root, all, "", true)
	}
	return nil
}

func printGraphTree(cmd *cobra.Command, n graphNode, all map[string]graphNode, prefix string, isRoot bool) {
	label := n.name
	if n.version != "" {
		label += "@" + n.version
	}
	if isRoot {
		fmt.Fprintln(cmd.OutOrStdout(), label)
	}
	for i, dep := range n.requires {
		isLast := i == len(n.requires)-1
		branch := prefix + "├── "
		childPrefix := prefix + "│   "
		if isLast {
			branch = prefix + "└── "
			childPrefix = prefix + "    "
		}
		depNode := all[dep]
		depLabel := dep
		if depNode.version != "" {
			depLabel = dep + "@" + depNode.version
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", branch, depLabel)
		if len(depNode.requires) > 0 {
			printGraphTree(cmd, depNode, all, childPrefix, false)
		}
	}
}

func readGraphNode(dir string) (graphNode, error) {
	sy, err := readSkillYAML(dir)
	if err != nil {
		return graphNode{}, err
	}
	return graphNode{
		name:     sy.Name,
		version:  sy.Version,
		requires: sy.Requires,
	}, nil
}

func resolveInstalledGraphNode(name string) (graphNode, error) {
	for _, root := range defaultSkillRoots() {
		dir := filepath.Join(root, name)
		if _, err := os.Stat(filepath.Join(dir, "skill.yaml")); err == nil {
			return readGraphNode(dir)
		}
	}
	return graphNode{name: name, version: "?"}, nil
}
