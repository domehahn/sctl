package cli

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type lockSummary struct {
	TotalSkills     int            `json:"total_skills"`
	UniqueSources   int            `json:"unique_sources"`
	ProtectedCount  int            `json:"protected_count"`
	SignedCount     int            `json:"signed_count"`
	VersionDist     map[string]int `json:"version_distribution"`
	OldestVersion   string         `json:"oldest_version,omitempty"`
	NewestVersion   string         `json:"newest_version,omitempty"`
}

func newSummarizeCmd() *cobra.Command {
	var lockPath string
	var protectPath string

	cmd := &cobra.Command{
		Use:   "summarize",
		Short: "Print a concise stats overview of the lockfile",
		Long: `Shows: total skill count, unique source registries, protected and signed
skill counts, and the version distribution (how many skills are at each version).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Lockfile is empty.")
				return nil
			}

			sources := map[string]struct{}{}
			signed := 0
			versions := map[string]int{}
			for _, sl := range lf.Skills {
				if sl.Source != "" {
					sources[sl.Source] = struct{}{}
				}
				if sl.Signature != "" {
					signed++
				}
				versions[sl.Version]++
			}

			protected := 0
			if protectPath != "" {
				for _, sl := range lf.Skills {
					if IsProtected(protectPath, sl.Name) {
						protected++
					}
				}
			}

			// Collect sorted versions for oldest/newest.
			versionKeys := make([]string, 0, len(versions))
			for v := range versions {
				versionKeys = append(versionKeys, v)
			}
			sort.Strings(versionKeys)

			oldest, newest := "", ""
			if len(versionKeys) > 0 {
				oldest = versionKeys[0]
				newest = versionKeys[len(versionKeys)-1]
			}

			summary := lockSummary{
				TotalSkills:   len(lf.Skills),
				UniqueSources: len(sources),
				ProtectedCount: protected,
				SignedCount:   signed,
				VersionDist:   versions,
				OldestVersion: oldest,
				NewestVersion: newest,
			}

			if outputFormat() == OutputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(summary)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Lockfile: %s\n\n", lockPath)
			fmt.Fprintf(cmd.OutOrStdout(), "  Skills:           %d\n", summary.TotalSkills)
			fmt.Fprintf(cmd.OutOrStdout(), "  Unique sources:   %d\n", summary.UniqueSources)
			fmt.Fprintf(cmd.OutOrStdout(), "  Protected:        %d\n", summary.ProtectedCount)
			fmt.Fprintf(cmd.OutOrStdout(), "  Signed:           %d\n", summary.SignedCount)
			if oldest != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Version range:    %s – %s\n", oldest, newest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n  Version distribution:\n")
			for _, v := range versionKeys {
				fmt.Fprintf(cmd.OutOrStdout(), "    %-20s  %d\n", v, versions[v])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&protectPath, "protect-file", protectFile, "Path to .skpm-protect for protected count")
	return cmd
}
