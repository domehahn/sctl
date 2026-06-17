package cli

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type namespaceSummary struct {
	Namespace string   `json:"namespace"`
	Count     int      `json:"count"`
	Skills    []string `json:"skills,omitempty"`
}

func newNamespaceCmd() *cobra.Command {
	var lockPath string
	var filterNS string

	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "List unique namespaces across installed skills",
		Long: `Groups installed skills by their namespace field and prints each namespace
with its skill count. Pass --namespace <ns> to list only skills in that namespace.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			// Group skills by namespace.
			nsMap := map[string][]string{}
			for _, sl := range lf.Skills {
				ns := sl.Namespace
				if ns == "" {
					ns = "default"
				}
				nsMap[ns] = append(nsMap[ns], sl.Name)
			}

			for ns := range nsMap {
				sort.Strings(nsMap[ns])
			}

			if filterNS != "" {
				skills, ok := nsMap[filterNS]
				if !ok {
					fmt.Fprintf(cmd.OutOrStdout(), "No skills in namespace %q.\n", filterNS)
					return nil
				}
				if outputFormat() == OutputJSON {
					return json.NewEncoder(cmd.OutOrStdout()).Encode(namespaceSummary{
						Namespace: filterNS,
						Count:     len(skills),
						Skills:    skills,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Namespace: %s (%d skill(s))\n\n", filterNS, len(skills))
				for _, s := range skills {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", s)
				}
				return nil
			}

			if len(nsMap) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills installed.")
				return nil
			}

			// Sort namespace names for stable output.
			nsKeys := make([]string, 0, len(nsMap))
			for ns := range nsMap {
				nsKeys = append(nsKeys, ns)
			}
			sort.Strings(nsKeys)

			if outputFormat() == OutputJSON {
				summaries := make([]namespaceSummary, 0, len(nsKeys))
				for _, ns := range nsKeys {
					summaries = append(summaries, namespaceSummary{
						Namespace: ns,
						Count:     len(nsMap[ns]),
					})
				}
				return json.NewEncoder(cmd.OutOrStdout()).Encode(summaries)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "  %-24s  count\n", "namespace")
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", "──────────────────────────────────")
			for _, ns := range nsKeys {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-24s  %d\n", ns, len(nsMap[ns]))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d namespace(s)\n", len(nsKeys))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&filterNS, "namespace", "", "Filter to a specific namespace and list its skills")
	return cmd
}
