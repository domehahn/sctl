package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show cache and installation statistics",
		Long: `Prints a summary of the skpm cache and installed skills:

  - Number and total size of cached ZIP artifacts
  - Number of installed skills per platform directory
  - Total disk usage of installed skills`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			// ── Cache stats ──────────────────────────────────────────
			cacheCount, cacheBytes := dirStats(cfg.CacheDir, func(name string) bool {
				return strings.HasSuffix(name, ".zip")
			})

			// ── Lockfile ─────────────────────────────────────────────
			lf, _ := lockfile.Read(lockfile.DefaultFilename)
			lockedCount := 0
			if lf != nil {
				lockedCount = len(lf.Skills)
			}

			// ── Per-platform install stats ───────────────────────────
			type platformStat struct {
				Root   string `json:"root"`
				Skills int    `json:"skills"`
				Bytes  int64  `json:"bytes"`
			}
			workDir, _ := os.Getwd()
			var platforms []platformStat
			totalInstallBytes := int64(0)
			totalInstallSkills := 0

			for _, root := range defaultSkillRoots() {
				entries, err := os.ReadDir(filepath.Join(workDir, root))
				if err != nil {
					continue
				}
				count := 0
				size := int64(0)
				for _, e := range entries {
					if !e.IsDir() {
						continue
					}
					if _, err := os.Stat(filepath.Join(workDir, root, e.Name(), "SKILL.md")); err != nil {
						continue
					}
					count++
					size += dirSize(filepath.Join(workDir, root, e.Name()))
				}
				if count == 0 {
					continue
				}
				platforms = append(platforms, platformStat{Root: root, Skills: count, Bytes: size})
				totalInstallSkills += count
				totalInstallBytes += size
			}

			sort.Slice(platforms, func(i, j int) bool { return platforms[i].Root < platforms[j].Root })

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "stats",
					Data: map[string]interface{}{
						"cache": map[string]interface{}{
							"artifacts": cacheCount,
							"bytes":     cacheBytes,
						},
						"lockfile": map[string]interface{}{
							"skills": lockedCount,
						},
						"installed": map[string]interface{}{
							"total_skills": totalInstallSkills,
							"total_bytes":  totalInstallBytes,
							"platforms":    platforms,
						},
					},
				})
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Cache\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  path:       %s\n", cfg.CacheDir)
			fmt.Fprintf(cmd.OutOrStdout(), "  artifacts:  %d\n", cacheCount)
			fmt.Fprintf(cmd.OutOrStdout(), "  size:       %s\n", humanBytes(cacheBytes))

			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintf(cmd.OutOrStdout(), "Lockfile\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  skills:     %d\n", lockedCount)

			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintf(cmd.OutOrStdout(), "Installed\n")
			if len(platforms) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  (no skills installed)\n")
			} else {
				for _, p := range platforms {
					fmt.Fprintf(cmd.OutOrStdout(), "  %-34s %3d skill(s)   %s\n",
						p.Root, p.Skills, humanBytes(p.Bytes))
				}
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprintf(cmd.OutOrStdout(), "  total:      %d skill(s), %s\n",
					totalInstallSkills, humanBytes(totalInstallBytes))
			}

			return nil
		},
	}
	return cmd
}

// dirStats counts files and their total size in dir, filtered by keep.
func dirStats(dir string, keep func(name string) bool) (int, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}
	count := 0
	total := int64(0)
	for _, e := range entries {
		if e.IsDir() || (keep != nil && !keep(e.Name())) {
			continue
		}
		info, _ := e.Info()
		if info != nil {
			count++
			total += info.Size()
		}
	}
	return count, total
}

// dirSize returns the total byte size of all files under root.
func dirSize(root string) int64 {
	total := int64(0)
	_ = filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total
}

// humanBytes formats a byte count as a human-readable string.
func humanBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
