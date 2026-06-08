package cli

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List locked or installed skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockfile.DefaultFilename)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockfile.DefaultFilename, err)}
			}
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "list", Data: lf.Skills})
				return nil
			}
			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills locked.")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%-28s %-12s %-16s %-28s %s\n", "Name", "Version", "Source", "Platforms", "Installed to")
			for _, sl := range lf.Skills {
				platforms := strings.Join(platformsToStrings(sl.CompatibleWith), ",")
				if platforms == "" {
					platforms = "-"
				}
				paths := strings.Join(sl.InstalledTo, ",")
				if paths == "" {
					paths = "-"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-28s %-12s %-16s %-28s %s\n", sl.Name, sl.Version, sl.Source, platforms, paths)
			}
			return nil
		},
	}
	return cmd
}
