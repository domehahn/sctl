package cli

import (
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func newFetchCmd() *cobra.Command {
	var lockPath string
	var skillFilter string

	cmd := &cobra.Command{
		Use:   "fetch [skill...]",
		Short: "Download skill ZIPs into the local cache without installing",
		Long: `Downloads the ZIP artifact for every locked skill into the local cache
so that subsequent 'skpm install --frozen-lockfile' runs require no network.

Typical CI usage:

  # Cache layer (network allowed):
  skpm fetch

  # Build layer (no network needed):
  skpm install --frozen-lockfile

Pass one or more skill names to fetch only those skills.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			// Filter skills.
			skills := lf.Skills
			if len(args) > 0 || skillFilter != "" {
				want := make(map[string]bool)
				for _, a := range args {
					want[a] = true
				}
				if skillFilter != "" {
					want[skillFilter] = true
				}
				filtered := skills[:0]
				for _, sl := range skills {
					if want[sl.Name] {
						filtered = append(filtered, sl)
					}
				}
				if len(filtered) == 0 {
					return &UserError{Message: "no matching skills found in lockfile"}
				}
				skills = filtered
			}

			if len(skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to fetch.")
				return nil
			}

			c := cache.New(cfg.CacheDir)

			type fetchResult struct {
				name      string
				fromCache bool
				skipped   bool
			}

			results := make([]fetchResult, len(skills))
			g, ctx := errgroup.WithContext(cmd.Context())
			sem := make(chan struct{}, globalConcurrency)

			for i, sl := range skills {
				i, sl := i, sl
				g.Go(func() error {
					sem <- struct{}{}
					defer func() { <-sem }()

					if sl.SHA256 != "" && c.Has(sl.SHA256) {
						log.Debug().Str("skill", sl.Name).Msg("already cached")
						results[i] = fetchResult{name: sl.Name, fromCache: true}
						return nil
					}

					if sl.DownloadURL == "" {
						results[i] = fetchResult{name: sl.Name, skipped: true}
						return nil
					}

					src := sl.Source
					if src == "" {
						src = cfg.DefaultRegistry
					}

					reg, err := registry.New(src, cfg)
					if err != nil {
						return fmt.Errorf("registry for %s: %w", sl.Name, err)
					}

					artifact := &registry.ResolvedArtifact{
						Name:        sl.Name,
						Version:     sl.Version,
						DownloadURL: sl.DownloadURL,
						SHA256:      sl.SHA256,
					}

					tmpKey := sl.SHA256
					if tmpKey == "" {
						tmpKey = "fetch-" + sl.Name + "-" + sl.Version
					}
					tmpPath := c.Path(tmpKey) + ".fetch"
					f, err := os.Create(tmpPath)
					if err != nil {
						return fmt.Errorf("create temp for %s: %w", sl.Name, err)
					}

					if err := reg.Download(ctx, artifact, f); err != nil {
						f.Close()
						os.Remove(tmpPath)
						return fmt.Errorf("download %s: %w", sl.Name, err)
					}
					f.Close()

					// Verify and rename to final cache path.
					actual, err := sha256File(tmpPath)
					if err != nil {
						os.Remove(tmpPath)
						return fmt.Errorf("hash %s: %w", sl.Name, err)
					}
					if sl.SHA256 != "" && actual != sl.SHA256 {
						os.Remove(tmpPath)
						return fmt.Errorf("SHA256 mismatch for %s: expected %s got %s", sl.Name, sl.SHA256, actual)
					}

					finalPath := c.Path(actual)
					if err := os.Rename(tmpPath, finalPath); err != nil {
						os.Remove(tmpPath)
						return fmt.Errorf("cache %s: %w", sl.Name, err)
					}

					results[i] = fetchResult{name: sl.Name}
					return nil
				})
			}

			if err := g.Wait(); err != nil {
				return &InternalError{Message: "fetch", Cause: err}
			}

			downloaded, fromCache, skipped := 0, 0, 0
			for _, r := range results {
				switch {
				case r.skipped:
					skipped++
				case r.fromCache:
					fromCache++
				default:
					downloaded++
				}
			}

			if format == OutputJSON {
				type row struct {
					Name      string `json:"name"`
					FromCache bool   `json:"from_cache"`
					Skipped   bool   `json:"skipped"`
				}
				rows := make([]row, len(results))
				for i, r := range results {
					rows[i] = row{Name: r.name, FromCache: r.fromCache, Skipped: r.skipped}
				}
				PrintResult(format, CommandResult{Success: true, Command: "fetch", Data: map[string]interface{}{
					"results":    rows,
					"downloaded": downloaded,
					"from_cache": fromCache,
					"skipped":    skipped,
				}})
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Fetched %d skill(s)", downloaded+fromCache)
			if fromCache > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), " (%d already cached)", fromCache)
			}
			if skipped > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), ", %d skipped (no download URL)", skipped)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().StringVar(&skillFilter, "skill", "", "Fetch only this skill")
	return cmd
}
