package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/domehahn/skpm/internal/config"
	"github.com/spf13/cobra"
)

func newCacheCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "cache", Short: "Manage skpm artifact cache"}
	cmd.AddCommand(newCacheListCmd())
	cmd.AddCommand(newCacheCleanCmd())
	return cmd
}

func newCacheListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List cached artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			entries, err := os.ReadDir(cfg.CacheDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "Cache is empty.")
					return nil
				}
				return &InternalError{Message: "read cache", Cause: err}
			}
			type item struct {
				Name string `json:"name"`
				Size int64  `json:"size"`
			}
			var items []item
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				info, _ := e.Info()
				items = append(items, item{Name: e.Name(), Size: info.Size()})
			}
			sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "cache list", Data: items})
				return nil
			}
			for _, item := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "%-72s %d bytes\n", item.Name, item.Size)
			}
			if len(items) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Cache is empty.")
			}
			return nil
		},
	}
}

func newCacheCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove cached artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove %s\n", cfg.CacheDir)
				return nil
			}
			if err := os.RemoveAll(filepath.Clean(cfg.CacheDir)); err != nil {
				return &InternalError{Message: "clean cache", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Cleaned %s\n", cfg.CacheDir)
			return nil
		},
	}
}
