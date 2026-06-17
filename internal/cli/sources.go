package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type sourceEntry struct {
	Source string   `json:"source"`
	Count  int      `json:"count"`
	Skills []string `json:"skills"`
}

func newSourcesCmd() *cobra.Command {
	var lockPath string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "sources",
		Short: "List all unique source URLs in the lockfile with skill counts",
		Long: `Groups lockfile entries by their source: field and prints each unique
source with a count of how many skills come from it.

Use --list-skills to list the individual skill names per source.
Use --output json for machine-readable output.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			sourceMap := map[string][]string{}
			for _, sl := range lf.Skills {
				src := sl.Source
				if src == "" {
					src = "(unknown)"
				}
				sourceMap[src] = append(sourceMap[src], sl.Name)
			}

			entries := make([]sourceEntry, 0, len(sourceMap))
			for src, skills := range sourceMap {
				sort.Strings(skills)
				entries = append(entries, sourceEntry{Source: src, Count: len(skills), Skills: skills})
			}
			sort.Slice(entries, func(i, j int) bool {
				if entries[i].Count != entries[j].Count {
					return entries[i].Count > entries[j].Count
				}
				return entries[i].Source < entries[j].Source
			})

			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills in lockfile.")
				return nil
			}

			if strings.ToLower(globalOutput) == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(entries)
			}

			for _, e := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-4d  %s\n", e.Count, e.Source)
				if verbose {
					for _, s := range e.Skills {
						fmt.Fprintf(cmd.OutOrStdout(), "         - %s\n", s)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().BoolVar(&verbose, "list-skills", false, "List skill names per source")
	return cmd
}
