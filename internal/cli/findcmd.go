package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newFindCmd() *cobra.Command {
	var lockPath string
	var namespace string

	cmd := &cobra.Command{
		Use:   "find <pattern>",
		Short: "Filter lockfile entries by name glob pattern",
		Long: `Prints lockfile entries whose name matches <pattern> (glob syntax: * and ?).
Optionally filter further with --namespace.

Useful for quickly locating a group of skills in large lockfiles.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]

			if pattern == "" {
				return &UserError{Message: "pattern must not be empty"}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			var matches []lockfile.SkillLock
			for _, sl := range lf.Skills {
				if !matchesGlob(pattern, sl.Name) {
					continue
				}
				if namespace != "" {
					ns := sl.Namespace
					if ns == "" {
						ns = "default"
					}
					if !strings.EqualFold(ns, namespace) {
						continue
					}
				}
				matches = append(matches, sl)
			}

			if len(matches) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No skills match %q.\n", pattern)
				return nil
			}

			if strings.ToLower(globalOutput) == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(matches)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%-32s  %-10s  %-12s\n", "NAME", "VERSION", "NAMESPACE")
			fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 58))
			for _, sl := range matches {
				ns := sl.Namespace
				if ns == "" {
					ns = "default"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-32s  %-10s  %-12s\n", sl.Name, sl.Version, ns)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Filter by namespace")
	return cmd
}
