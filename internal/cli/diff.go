package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	var source string
	var fromVersion string
	var toVersion string

	cmd := &cobra.Command{
		Use:   "diff <name>",
		Short: "Show SKILL.md diff between two versions of a skill",
		Long: `Downloads SKILL.md for two versions of a skill and shows a unified diff.

--from defaults to the locked version; --to defaults to the latest available version.
Useful for reviewing changes before running skpm update.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			src := sourceOrDefault(source, cfg)

			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: err.Error()}
			}

			ref := registry.ParseSkillRef(name, "default")

			// Default --from to the locked version.
			if fromVersion == "" {
				if lf, lerr := lockfile.Read(lockfile.DefaultFilename); lerr == nil {
					for _, sl := range lf.Skills {
						if sl.Name == name {
							fromVersion = sl.Version
							break
						}
					}
				}
			}
			if fromVersion == "" {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile; specify --from <version>", name)}
			}

			// Default --to to the latest available version.
			if toVersion == "" {
				discovery, ok := reg.(registry.DiscoveryRegistry)
				if !ok {
					return &UserError{Message: "registry does not support version discovery; use --to <version>"}
				}
				versions, err := discovery.ListVersions(cmd.Context(), ref)
				if err != nil {
					return &UserError{Message: fmt.Sprintf("list versions for %q: %v", name, err)}
				}
				toVersion = latestVersion(versions)
				if toVersion == "" {
					return &UserError{Message: fmt.Sprintf("no versions found for %q", name)}
				}
			}

			if fromVersion == toVersion {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: --from and --to are both %s; nothing to diff.\n", name, fromVersion)
				return nil
			}

			oldContent, err := fetchSkillMD(cmd.Context(), reg, ref, fromVersion)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("fetch %s@%s: %v", name, fromVersion, err)}
			}
			newContent, err := fetchSkillMD(cmd.Context(), reg, ref, toVersion)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("fetch %s@%s: %v", name, toVersion, err)}
			}

			if oldContent == newContent {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: SKILL.md is identical between %s and %s\n", name, fromVersion, toVersion)
				return nil
			}

			d := difflib.UnifiedDiff{
				A:        difflib.SplitLines(oldContent),
				B:        difflib.SplitLines(newContent),
				FromFile: fmt.Sprintf("%s@%s", name, fromVersion),
				ToFile:   fmt.Sprintf("%s@%s", name, toVersion),
				Context:  3,
			}
			text, err := difflib.GetUnifiedDiffString(d)
			if err != nil {
				return &InternalError{Message: "generate diff", Cause: err}
			}
			fmt.Fprint(cmd.OutOrStdout(), text)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source")
	cmd.Flags().StringVar(&fromVersion, "from", "", "Old version (default: locked version)")
	cmd.Flags().StringVar(&toVersion, "to", "", "New version (default: latest available)")
	return cmd
}

// fetchSkillMD downloads the skill package for the given version and extracts SKILL.md.
func fetchSkillMD(ctx context.Context, reg registry.Registry, ref registry.SkillRef, version string) (string, error) {
	artifact, err := reg.Resolve(ctx, registry.ResolveRequest{Ref: ref, Constraint: version})
	if err != nil {
		return "", fmt.Errorf("resolve: %w", err)
	}

	tmp, err := os.CreateTemp("", "skpm-diff-*.zip")
	if err != nil {
		return "", fmt.Errorf("temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if err := reg.Download(ctx, artifact, tmp); err != nil {
		return "", fmt.Errorf("download: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}

	// Try SKILL.md at the top level and inside a skill-name/ subdirectory.
	candidates := []string{"SKILL.md", ref.Name + "/SKILL.md"}
	for _, candidate := range candidates {
		data, err := readZipFile(tmp.Name(), candidate)
		if err == nil {
			return strings.TrimSpace(string(data)), nil
		}
	}
	return "", fmt.Errorf("SKILL.md not found in package for %s@%s", ref.Name, version)
}
