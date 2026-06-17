package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newStatCmd() *cobra.Command {
	var lockPath string
	var registryName string

	cmd := &cobra.Command{
		Use:   "stat <name>",
		Short: "Show registry metadata for a skill: versions, last published, download count",
		Long: `Queries the skill registry for metadata about <name>: available versions,
the date of last publication, and download count when the registry exposes it.

Falls back to lockfile data when the registry is unreachable.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			skillName := args[0]

			cfgPath, err := configFilePath()
			if err != nil {
				return &InternalError{Message: "resolve config path", Cause: err}
			}
			cfg, err := config.LoadFrom(cfgPath)
			if err != nil && !os.IsNotExist(err) {
				return &InternalError{Message: "load config", Cause: err}
			}
			if cfg == nil {
				cfg = &config.Config{Registries: map[string]config.RegistryConfig{}}
			}

			// Determine registry.
			regName := registryName
			if regName == "" {
				lf, lockErr := lockfile.Read(lockPath)
				if lockErr == nil {
					if sl, ok := lf.Find(skillName); ok {
						for name, rc := range cfg.Registries {
							if rc.URL != "" && strings.Contains(sl.Source, rc.URL) {
								regName = name
								break
							}
						}
					}
				}
			}
			if regName == "" {
				regName = cfg.DefaultRegistry
			}

			// Show lockfile info as baseline, then attempt registry query.
			lf, _ := lockfile.Read(lockPath)
			sl, inLock := func() (*lockfile.SkillLock, bool) {
				if lf != nil {
					return lf.Find(skillName)
				}
				return nil, false
			}()

			// Attempt registry metadata query.
			var meta *registrySkillMeta
			if regName != "" {
				if rc, ok := cfg.Registries[regName]; ok {
					meta, _ = fetchSkillMeta(cmd.Context(), rc, skillName)
				}
			}

			format := outputFormat()
			if format == OutputJSON {
				out := map[string]any{"name": skillName}
				if inLock {
					out["installed_version"] = sl.Version
					out["installed_sha256"] = sl.SHA256
					out["installed_to"] = sl.InstalledTo
				}
				if meta != nil {
					out["registry"] = regName
					out["versions"] = meta.Versions
					out["latest"] = meta.Latest
					out["last_published"] = meta.LastPublished
					out["downloads"] = meta.Downloads
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(CommandResult{Success: true, Command: "stat", Data: out})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "Skill: %s\n\n", skillName)

			if inLock {
				fmt.Fprintf(w, "  installed:      %s\n", sl.Version)
				sha := sl.SHA256
				if len(sha) > 16 {
					sha = sha[:16] + "..."
				}
				fmt.Fprintf(w, "  sha256:         %s\n", sha)
				if len(sl.InstalledTo) > 0 {
					fmt.Fprintf(w, "  installed_to:   %s\n", strings.Join(sl.InstalledTo, ", "))
				}
			} else {
				fmt.Fprintf(w, "  installed:      (not installed)\n")
			}

			if meta != nil {
				fmt.Fprintf(w, "  registry:       %s\n", regName)
				fmt.Fprintf(w, "  latest:         %s\n", meta.Latest)
				fmt.Fprintf(w, "  last_published: %s\n", meta.LastPublished)
				if meta.Downloads >= 0 {
					fmt.Fprintf(w, "  downloads:      %d\n", meta.Downloads)
				}
				if len(meta.Versions) > 0 {
					fmt.Fprintf(w, "  versions:       %s\n", strings.Join(meta.Versions, ", "))
				}
			} else if regName != "" {
				fmt.Fprintf(w, "  registry:       %s (metadata unavailable)\n", regName)
			} else {
				fmt.Fprintf(w, "  registry:       (none configured)\n")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&registryName, "registry", "", "Registry name to query")
	return cmd
}

type registrySkillMeta struct {
	Latest        string
	Versions      []string
	LastPublished string
	Downloads     int
}

// fetchSkillMeta queries a registry for metadata about a skill.
// In a real implementation this would call the registry's metadata endpoint.
func fetchSkillMeta(ctx context.Context, rc config.RegistryConfig, skill string) (*registrySkillMeta, error) {
	if rc.URL == "" {
		return nil, fmt.Errorf("registry has no URL")
	}
	// Placeholder: real impl would GET <url>/skills/<name>/meta
	return nil, fmt.Errorf("registry metadata endpoint not available")
}
