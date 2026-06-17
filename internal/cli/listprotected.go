package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newListProtectedCmd() *cobra.Command {
	var lockDir string

	cmd := &cobra.Command{
		Use:   "list-protected",
		Short: "List all skills marked as protected from upgrade",
		Long:  `Reads ` + protectFile + ` and prints every protected skill name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			storePath := filepath.Join(lockDir, protectFile)
			names := ProtectedList(storePath)

			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No protected skills.")
				return nil
			}

			if outputFormat() == OutputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(names)
			}

			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), " ", n)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d protected skill(s).\n", len(names))
			return nil
		},
	}

	cmd.Flags().StringVar(&lockDir, "lock-dir", ".", "Directory containing "+protectFile)
	return cmd
}
