package cli

import (
	"fmt"
	"os"

	"github.com/domehahn/skpm/internal/config"
	"github.com/domehahn/skpm/internal/lockfile"
	"github.com/domehahn/skpm/internal/manifest"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check skpm project and configuration health",
		RunE: func(cmd *cobra.Command, args []string) error {
			var issues []string
			if _, err := config.Load(); err != nil {
				issues = append(issues, fmt.Sprintf("config: %v", err))
			}
			if _, err := manifest.Read(manifest.DefaultFilename); err != nil {
				issues = append(issues, fmt.Sprintf("%s: %v", manifest.DefaultFilename, err))
			}
			if _, err := lockfile.Read(lockfile.DefaultFilename); err != nil && !os.IsNotExist(err) {
				issues = append(issues, fmt.Sprintf("%s: %v", lockfile.DefaultFilename, err))
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: len(issues) == 0, Command: "doctor", Errors: issues})
				return nil
			}
			if len(issues) > 0 {
				return &UserError{Message: fmt.Sprintf("doctor found issues:\n  %s", stringsJoin(issues, "\n  "))}
			}
			fmt.Fprintln(cmd.OutOrStdout(), "skpm doctor passed")
			return nil
		},
	}
}
