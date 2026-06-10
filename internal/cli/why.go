package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

func newWhyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "why <skill>",
		Short: "Explain why a skill is in the lockfile",
		Long: `Shows the manifest constraint, resolved version, registry source,
and installation paths for a locked skill.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			format := outputFormat()

			lf, err := lockfile.Read(lockfile.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockfile.DefaultFilename, err)}
			}

			sl, ok := lf.Find(name)
			if !ok {
				return &UserError{Message: fmt.Sprintf("%q is not in %s", name, lockfile.DefaultFilename)}
			}

			// Find manifest entry for constraint.
			mf, _ := manifest.Read(manifest.DefaultFilename)
			constraint := manifestConstraint(mf, name)
			direct := constraint != ""

			// Check which install paths actually exist on disk.
			workDir, _ := os.Getwd()
			type pathStatus struct {
				Path      string `json:"path"`
				Installed bool   `json:"installed"`
			}
			statuses := make([]pathStatus, 0, len(sl.InstalledTo))
			for _, p := range sl.InstalledTo {
				_, statErr := os.Stat(filepath.Join(workDir, p, "SKILL.md"))
				statuses = append(statuses, pathStatus{Path: p, Installed: statErr == nil})
			}

			if format == OutputJSON {
				type jsonResult struct {
					Name       string       `json:"name"`
					Version    string       `json:"version"`
					Constraint string       `json:"constraint,omitempty"`
					Direct     bool         `json:"direct"`
					Source     string       `json:"source,omitempty"`
					SHA256     string       `json:"sha256,omitempty"`
					Paths      []pathStatus `json:"paths"`
				}
				PrintResult(format, CommandResult{
					Success: true,
					Command: "why",
					Data: jsonResult{
						Name:       sl.Name,
						Version:    sl.Version,
						Constraint: constraint,
						Direct:     direct,
						Source:     sl.Source,
						SHA256:     sl.SHA256,
						Paths:      statuses,
					},
				})
				return nil
			}

			origin := "transitive / manually added"
			if direct {
				orig := constraint
				if orig == "" {
					orig = "(any)"
				}
				origin = fmt.Sprintf("direct — declared in %s with constraint %q", manifest.DefaultFilename, orig)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", sl.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "  version:    %s\n", sl.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  origin:     %s\n", origin)
			if sl.Source != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  registry:   %s\n", sl.Source)
			}
			if sl.SHA256 != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  sha256:     %s\n", sl.SHA256[:12]+"…")
			}

			if len(statuses) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  installed:\n")
				for _, ps := range statuses {
					icon := "✓"
					note := ""
					if !ps.Installed {
						icon = "✗"
						note = "  (missing — run: skpm install)"
					}
					fmt.Fprintf(cmd.OutOrStdout(), "    %s %s%s\n", icon, ps.Path, note)
				}
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "  installed:  (no paths recorded — run: skpm install)\n")
			}

			linked, linkedPath := isLinkedSkill(name)
			if linked {
				fmt.Fprintf(cmd.OutOrStdout(), "  linked:     %s (local development link)\n", linkedPath)
			}

			// Platforms
			if len(sl.CompatibleWith) > 0 {
				platforms := make([]string, len(sl.CompatibleWith))
				for i, p := range sl.CompatibleWith {
					platforms[i] = string(p)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  platforms:  %s\n", strings.Join(platforms, ", "))
			}

			return nil
		},
	}
	return cmd
}
