package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type QueryMatch struct {
	Source      string `json:"source"` // "lockfile" or "installed"
	Skill       string `json:"skill"`
	Version     string `json:"version"`
	Path        string `json:"path,omitempty"`
	Digest      string `json:"digest,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

func newQueryCmd() *cobra.Command {
	var (
		digest   string
		yanked   bool
		lockPath string
		dir      string
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query installed skills and lockfiles for incident response",
		Long: `Performs audit and incident response queries across installed skills and lockfiles.

Examples:
  skpm query --digest sha256:abc123...
  skpm query --yanked
  skpm query --digest abc123... --dir ./skills`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			if digest == "" && !yanked {
				return &UserError{Message: "must specify at least one query filter (--digest or --yanked)"}
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			if dir == "" {
				dir = "."
			}

			var matches []QueryMatch

			// Search lockfile
			if lf, err := lockfile.Read(lockPath); err == nil {
				for _, sl := range lf.Skills {
					match := false
					desc := ""

					if digest != "" {
						cleanSearch := strings.TrimPrefix(digest, "sha256:")
						cleanLock := strings.TrimPrefix(sl.SHA256, "sha256:")
						if cleanSearch == cleanLock || strings.HasPrefix(cleanLock, cleanSearch) {
							match = true
							desc = fmt.Sprintf("lockfile entry SHA256 matches %s", digest)
						}
					}

					if yanked {
						if sl.Metadata != nil && (sl.Metadata["yanked"] == "true" || sl.Metadata["deprecated"] == "true") {
							match = true
							if desc != "" {
								desc += "; "
							}
							desc += "lockfile version is marked yanked/deprecated"
						}
					}

					if match {
						matches = append(matches, QueryMatch{
							Source:      "lockfile",
							Skill:       sl.Name,
							Version:     sl.Version,
							Path:        lockPath,
							Digest:      sl.SHA256,
							Status:      sl.Metadata["status"],
							Description: desc,
						})
					}
				}
			}

			// Search installed directory for .skpm-installed.json
			_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() || d.Name() != ".skpm-installed.json" {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				var id installer.MaterializedIdentity
				if err := json.Unmarshal(data, &id); err != nil {
					return nil
				}

				match := false
				desc := ""

				if digest != "" {
					cleanSearch := strings.TrimPrefix(digest, "sha256:")
					cleanPkg := strings.TrimPrefix(id.PackageDigest, "sha256:")
					cleanComp := strings.TrimPrefix(id.CompiledDigest, "sha256:")
					cleanArt := strings.TrimPrefix(id.ArtifactDigest, "sha256:")

					if cleanSearch == cleanPkg || cleanSearch == cleanComp || cleanSearch == cleanArt ||
						strings.HasPrefix(cleanPkg, cleanSearch) || strings.HasPrefix(cleanComp, cleanSearch) {
						match = true
						desc = fmt.Sprintf("installed artifact digest matches %s", digest)
					}
				}

				if match {
					matches = append(matches, QueryMatch{
						Source:      "installed",
						Skill:       id.Package,
						Version:     id.Version,
						Path:        filepath.Dir(path),
						Digest:      id.PackageDigest,
						Description: desc,
					})
				}
				return nil
			})

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "query",
					Data: map[string]interface{}{
						"query_digest": digest,
						"query_yanked": yanked,
						"matches":      matches,
						"total":        len(matches),
					},
				})
				return nil
			}

			if len(matches) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No matching skills or lockfile entries found.")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Found %d match(es):\n\n", len(matches))
			for _, m := range matches {
				fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s (Version: %s)\n", m.Source, m.Skill, m.Version)
				if m.Path != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "    Path:   %s\n", m.Path)
				}
				if m.Digest != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "    Digest: %s\n", m.Digest)
				}
				if m.Description != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "    Note:   %s\n", m.Description)
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&digest, "digest", "", "Artifact or package SHA256 digest to search for")
	cmd.Flags().BoolVar(&yanked, "yanked", false, "Search for yanked or deprecated versions")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().StringVar(&dir, "dir", "", "Root directory to search installed skills (default: current dir)")
	return cmd
}
