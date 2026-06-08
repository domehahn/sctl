package cli

import (
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
)

func newLintCmd() *cobra.Command {
	var platform string
	var allowPrerelease bool
	var allowMissingTests bool
	var allowMissingLicense bool

	cmd := &cobra.Command{
		Use:   "lint [path]",
		Short: "Lint a skill directory (alias for validate --strict)",
		Long: `skpm lint is a convenience alias for skpm validate --strict.

Strict validation is suitable for CI: warnings become errors, CHANGELOG.md and
README.md are required, absolute local paths are disallowed, and generated
artifacts must not be checked in.

All flags and output formats work identically to skpm validate --strict.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("path not found: %s", dir)}
			}

			v := skill.NewValidatorWithOptions(skill.ValidationOptions{
				Strict:              true,
				Platform:            skill.Platform(platform),
				AllowPrerelease:     allowPrerelease,
				AllowMissingTests:   allowMissingTests,
				AllowMissingLicense: allowMissingLicense,
			})
			result, err := v.Validate(cmd.Context(), dir)
			if err != nil {
				return &InternalError{Message: "validate", Cause: err}
			}

			printValidateResult(cmd, outputFormat(), dir, result)
			if !result.Valid {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&platform, "platform", "", "Validate compatibility with a specific target platform")
	cmd.Flags().BoolVar(&allowPrerelease, "allow-prerelease", false, "Accept prerelease SemVer versions")
	cmd.Flags().BoolVar(&allowMissingTests, "allow-missing-tests", false, "Downgrade missing tests/ from error to warning")
	cmd.Flags().BoolVar(&allowMissingLicense, "allow-missing-license", false, "Downgrade missing LICENSE from error to warning")
	return cmd
}
