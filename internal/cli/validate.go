package cli

import (
	"fmt"
	"os"

	"github.com/domehahn/skpm/internal/skill"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var strict bool
	var publish bool
	var platform string
	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate a skill's structure",
		Long: `Checks that a skill directory contains all required files with correct content:
  - SKILL.md (non-empty)
  - VERSION (valid semver)
  - skill.yaml (name, version, description, compatible_with required; version matches VERSION)
  - CHANGELOG.md (entry for current version recommended)

Exits 0 if valid, 1 if errors are found.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("path not found: %s", dir)}
			}

			v := skill.NewValidatorWithOptions(skill.ValidationOptions{
				Strict:   strict,
				Publish:  publish,
				Platform: skill.Platform(platform),
			})
			result, err := v.Validate(cmd.Context(), dir)
			if err != nil {
				return &InternalError{Message: "validate", Cause: err}
			}

			if format == OutputJSON {
				type jsonError struct {
					Field    string `json:"field"`
					Message  string `json:"message"`
					Severity string `json:"severity"`
				}
				errList := make([]jsonError, len(result.Errors))
				for i, e := range result.Errors {
					errList[i] = jsonError{Field: e.Field, Message: e.Message, Severity: string(e.Severity)}
				}
				warnList := make([]jsonError, len(result.Warnings))
				for i, e := range result.Warnings {
					warnList[i] = jsonError{Field: e.Field, Message: e.Message, Severity: string(e.Severity)}
				}
				errStrs := make([]string, len(result.Errors))
				for i, e := range result.Errors {
					errStrs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
				}
				PrintResult(format, CommandResult{
					Success: result.Valid,
					Command: "validate",
					Data:    map[string]interface{}{"path": dir, "errors": errList, "warnings": warnList},
					Errors:  errStrs,
				})
				if !result.Valid {
					os.Exit(1)
				}
				return nil
			}

			if result.Valid {
				fmt.Fprintf(cmd.OutOrStdout(), "✓ %s is valid", dir)
				if len(result.Warnings) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), " (%d warning(s))", len(result.Warnings))
				}
				fmt.Fprintln(cmd.OutOrStdout())
			} else {
				fmt.Fprintf(cmd.ErrOrStderr(), "✗ %s has %d error(s)\n", dir, len(result.Errors))
			}

			for _, e := range result.Errors {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ERROR   [%s] %s\n", e.Field, e.Message)
			}
			for _, w := range result.Warnings {
				fmt.Fprintf(cmd.OutOrStdout(), "  WARNING [%s] %s\n", w.Field, w.Message)
			}

			if !result.Valid {
				os.Exit(1)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&strict, "strict", false, "Treat warnings as errors")
	cmd.Flags().BoolVar(&publish, "publish", false, "Apply public package validation profile")
	cmd.Flags().StringVar(&platform, "platform", "", "Validate compatibility with a target platform")
	return cmd
}
