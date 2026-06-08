package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/internal/config"
	"github.com/domehahn/skpm/internal/lockfile"
	"github.com/domehahn/skpm/internal/skill"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newVerifyCmd() *cobra.Command {
	var frozen bool
	var platform string
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify lockfile and installed skill integrity",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			if frozen {
				outdated, err := lockfileOutdated(cmd.Context(), cfg, "agent-skills.yaml", lockfile.DefaultFilename)
				if err != nil {
					return &UserError{Message: fmt.Sprintf("check lockfile: %v", err)}
				}
				if outdated {
					return &UserError{Message: "agent-skills.lock is inconsistent with agent-skills.yaml"}
				}
			}
			lf, err := lockfile.Read(lockfile.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockfile.DefaultFilename, err)}
			}
			var errors []string
			workDir, _ := os.Getwd()
			for _, sl := range lf.Skills {
				if sl.SHA256 == "" {
					errors = append(errors, fmt.Sprintf("%s: missing sha256", sl.Name))
				}
				for _, p := range sl.CompatibleWith {
					if !skill.KnownPlatforms[skill.Platform(p)] {
						errors = append(errors, fmt.Sprintf("%s: unknown platform %s", sl.Name, p))
					}
				}
				if platform != "" && len(sl.CompatibleWith) > 0 {
					matches := false
					for _, p := range sl.CompatibleWith {
						matches = matches || p == platform || p == string(skill.PlatformAll)
					}
					if !matches {
						errors = append(errors, fmt.Sprintf("%s: not compatible with %s", sl.Name, platform))
					}
				}
				for _, installPath := range sl.InstalledTo {
					dir := filepath.Join(workDir, installPath)
					if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err != nil {
						errors = append(errors, fmt.Sprintf("%s: missing %s/SKILL.md", sl.Name, installPath))
					}
					data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
					if err != nil {
						errors = append(errors, fmt.Sprintf("%s: missing %s/skill.yaml", sl.Name, installPath))
						continue
					}
					var sy skill.SkillYAML
					if err := yaml.Unmarshal(data, &sy); err != nil {
						errors = append(errors, fmt.Sprintf("%s: invalid skill.yaml: %v", sl.Name, err))
						continue
					}
					if sy.Name != sl.Name || sy.Version != sl.Version {
						errors = append(errors, fmt.Sprintf("%s: installed metadata %s@%s does not match lockfile %s@%s", sl.Name, sy.Name, sy.Version, sl.Name, sl.Version))
					}
				}
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: len(errors) == 0, Command: "verify", Errors: errors})
			}
			if len(errors) > 0 {
				return &UserError{Message: fmt.Sprintf("verification failed:\n  %s", stringsJoin(errors, "\n  "))}
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Verification passed")
			return nil
		},
	}
	cmd.Flags().BoolVar(&frozen, "frozen-lockfile", false, "Fail if manifest and lockfile are inconsistent")
	cmd.Flags().StringVar(&platform, "platform", "", "Verify compatibility for a platform")
	return cmd
}

func stringsJoin(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}
