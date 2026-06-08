package cli

import (
	"fmt"
	"os"

	"github.com/domehahn/skpm/internal/skill"
	"github.com/spf13/cobra"
)

func newPackageCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "package [path]",
		Short: "Package a skill directory into a versioned ZIP artifact",
		Long: `Validates the skill, then creates <name>-<version>.zip containing all skill
files plus a manifest.json with build metadata.

The SHA256 of the final ZIP is printed — use it in agent-skills.lock.`,
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

			if outputDir == "" {
				outputDir = dir
			}

			p := skill.NewPackager()
			result, err := p.Package(cmd.Context(), dir, outputDir)
			if err != nil {
				var vfe *skill.ValidationFailedError
				if isValidationFailedError(err, &vfe) {
					if format == OutputJSON {
						PrintResult(format, CommandResult{
							Success: false,
							Command: "package",
							Errors:  vfe.Errors,
						})
						os.Exit(1)
						return nil
					}
					fmt.Fprintln(cmd.ErrOrStderr(), "Skill validation failed:")
					for _, e := range vfe.Errors {
						fmt.Fprintf(cmd.ErrOrStderr(), "  - %s\n", e)
					}
					os.Exit(1)
					return nil
				}
				return &InternalError{Message: "package", Cause: err}
			}

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "package",
					Data: map[string]interface{}{
						"name":    result.Name,
						"version": result.Version,
						"path":    result.OutputPath,
						"sha256":  result.SHA256,
					},
				})
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Packaged %s@%s\n", result.Name, result.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  File:   %s\n", result.OutputPath)
			fmt.Fprintf(cmd.OutOrStdout(), "  SHA256: %s\n", result.SHA256)
			return nil
		},
	}

	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory for the output ZIP (default: skill directory)")
	return cmd
}

func isValidationFailedError(err error, target **skill.ValidationFailedError) bool {
	if e, ok := err.(*skill.ValidationFailedError); ok {
		*target = e
		return true
	}
	return false
}
