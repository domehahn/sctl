package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var strict bool
	var publish bool
	var platform string
	var allowPrerelease bool
	var allowMissingTests bool
	var allowMissingLicense bool

	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate a skill directory",
		Long: `Validates a canonical skill source directory.

Profiles:
  default  Useful during local development. Warns on missing optional files.
  --strict  CI-grade. Warnings become errors; stricter structural checks.
  --publish  Release-grade. All strict checks plus publish-readiness requirements.

Minimum required files (default):   SKILL.md  VERSION  skill.yaml
Minimum required files (--publish):  above + CHANGELOG.md

Exits 0 if valid, 1 if errors are found.`,
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
				Strict:              strict,
				Publish:             publish,
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

	cmd.Flags().BoolVar(&strict, "strict", false, "CI-grade: treat warnings as errors and apply stricter checks")
	cmd.Flags().BoolVar(&publish, "publish", false, "Release-grade: all strict checks plus publish-readiness requirements")
	cmd.Flags().StringVar(&platform, "platform", "", "Validate compatibility with a specific target platform")
	cmd.Flags().BoolVar(&allowPrerelease, "allow-prerelease", false, "Accept prerelease SemVer versions")
	cmd.Flags().BoolVar(&allowMissingTests, "allow-missing-tests", false, "Downgrade missing tests/ from error to warning in strict/publish")
	cmd.Flags().BoolVar(&allowMissingLicense, "allow-missing-license", false, "Downgrade missing LICENSE from error to warning in strict/publish")
	return cmd
}

// printValidateResult renders a ValidationResult to text or JSON.
// Used by both validate and lint commands.
func printValidateResult(cmd *cobra.Command, format OutputFormat, dir string, result *skill.ValidationResult) {
	if format == OutputJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return
	}

	profile := result.Profile
	if profile == "" {
		profile = "default"
	}

	if result.Valid {
		fmt.Fprintf(cmd.OutOrStdout(), "Validation passed: %s\n", dir)
		fmt.Fprintf(cmd.OutOrStdout(), "Profile: %s\n", profile)
		fmt.Fprintf(cmd.OutOrStdout(), "Warnings: %d\n", len(result.Warnings))
		fmt.Fprintf(cmd.OutOrStdout(), "Errors: %d\n", len(result.Errors))
	} else {
		fmt.Fprintf(cmd.ErrOrStderr(), "Validation failed: %s\n", dir)
		fmt.Fprintf(cmd.ErrOrStderr(), "Profile: %s\n", profile)
		fmt.Fprintln(cmd.ErrOrStderr())
	}

	if len(result.Errors) > 0 {
		fmt.Fprintln(cmd.ErrOrStderr(), "Errors:")
		for _, e := range result.Errors {
			fmt.Fprintf(cmd.ErrOrStderr(), "  - %s\n", formatFinding(e))
		}
	}
	if len(result.Warnings) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Warnings:")
		for _, w := range result.Warnings {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", formatFinding(w))
		}
	}
	if len(result.Infos) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Info:")
		for _, i := range result.Infos {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", formatFinding(i))
		}
	}
}

func formatFinding(f skill.ValidationFinding) string {
	field := f.Field
	if f.Path != "" && f.Path != f.Field {
		field = f.Path
	}
	msg := f.Message
	if f.Code != "" {
		msg = fmt.Sprintf("%s [%s]", msg, f.Code)
	}
	if field != "" {
		return fmt.Sprintf("%s: %s", field, msg)
	}
	return msg
}

func profileFromFlags(strict, publish bool) string {
	if publish {
		return "publish"
	}
	if strict {
		return "strict"
	}
	return "default"
}

// profileLabel returns a human-readable label, capitalising the first letter.
func profileLabel(p string) string {
	if p == "" {
		return "default"
	}
	return strings.ToUpper(p[:1]) + p[1:]
}
