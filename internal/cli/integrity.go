package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type integrityStatus int

const (
	integrityOK integrityStatus = iota
	integrityMissing
	integrityModified
	integrityUnverifiable
)

func (s integrityStatus) String() string {
	switch s {
	case integrityOK:
		return "ok"
	case integrityMissing:
		return "missing"
	case integrityModified:
		return "modified"
	case integrityUnverifiable:
		return "unverifiable"
	default:
		return "unknown"
	}
}

type integrityResult struct {
	Name    string          `json:"name"`
	Version string          `json:"version"`
	Path    string          `json:"path"`
	Status  integrityStatus `json:"status"`
	Note    string          `json:"note,omitempty"`
}

func newIntegrityCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "integrity",
		Short: "Verify installed skills match their cached ZIP artifacts",
		Long: `Compares installed SKILL.md files against the content extracted from
the locally cached ZIP artifact (identified by the SHA256 in the lockfile).

Status meanings:
  ok            — installed content matches the cached artifact
  modified      — installed SKILL.md differs from the artifact (manual edits or tampering)
  missing       — the skill directory or SKILL.md is absent (run: skpm install)
  unverifiable  — no cached artifact and no SHA256 in lockfile (run: skpm install to cache)

Exits with code 1 if any skill is modified or missing.`,
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
			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills in lockfile.")
				return nil
			}

			c := cache.New(cfg.CacheDir)
			workDir, _ := os.Getwd()

			var results []integrityResult
			for _, sl := range lf.Skills {
				// Extract SKILL.md from cached ZIP (if available).
				var canonicalSKILL []byte
				if sl.SHA256 != "" && c.Has(sl.SHA256) {
					zipPath := c.Path(sl.SHA256)
					data, err := readZipFile(zipPath, "SKILL.md")
					if err != nil {
						// Try with a common prefix (e.g. name-version/SKILL.md).
						data, _ = readZipFile(zipPath, sl.Name+"-"+sl.Version+"/SKILL.md")
					}
					canonicalSKILL = data
				}

				for _, p := range sl.InstalledTo {
					installPath := filepath.Join(workDir, p, "SKILL.md")
					r := integrityResult{
						Name:    sl.Name,
						Version: sl.Version,
						Path:    filepath.Join(p, "SKILL.md"),
					}

					installed, readErr := os.ReadFile(installPath)
					if readErr != nil {
						r.Status = integrityMissing
						r.Note = "SKILL.md not found"
						results = append(results, r)
						continue
					}

					if canonicalSKILL == nil {
						r.Status = integrityUnverifiable
						r.Note = "artifact not in cache (run: skpm install)"
						results = append(results, r)
						continue
					}

					if string(installed) == string(canonicalSKILL) {
						r.Status = integrityOK
					} else {
						r.Status = integrityModified
						r.Note = "content differs from artifact"
					}
					results = append(results, r)
				}

				// Skill with no recorded install paths.
				if len(sl.InstalledTo) == 0 {
					r := integrityResult{
						Name:    sl.Name,
						Version: sl.Version,
						Path:    "(no install paths recorded)",
					}
					if canonicalSKILL != nil {
						r.Status = integrityUnverifiable
						r.Note = "no install paths in lockfile"
					} else {
						r.Status = integrityUnverifiable
						r.Note = "no install paths and artifact not cached"
					}
					results = append(results, r)
				}
			}

			counts := map[integrityStatus]int{}
			for _, r := range results {
				counts[r.Status]++
			}

			if format == OutputJSON {
				type jsonResult struct {
					Name    string `json:"name"`
					Version string `json:"version"`
					Path    string `json:"path"`
					Status  string `json:"status"`
					Note    string `json:"note,omitempty"`
				}
				out := make([]jsonResult, len(results))
				for i, r := range results {
					out[i] = jsonResult{Name: r.Name, Version: r.Version, Path: r.Path, Status: r.Status.String(), Note: r.Note}
				}
				PrintResult(format, CommandResult{Success: true, Command: "integrity", Data: out})
				return nil
			}

			for _, r := range results {
				icon := integrityIcon(r.Status)
				note := ""
				if r.Note != "" {
					note = "  (" + r.Note + ")"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s %-28s  %s%s\n", icon, r.Name+"@"+r.Version, r.Path, note)
			}

			fmt.Fprintln(cmd.OutOrStdout())
			parts := []string{}
			if n := counts[integrityOK]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d ok", n))
			}
			if n := counts[integrityModified]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d modified", n))
			}
			if n := counts[integrityMissing]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d missing", n))
			}
			if n := counts[integrityUnverifiable]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d unverifiable", n))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Integrity: %d path(s) — %s\n",
				len(results), strings.Join(parts, ", "))

			if counts[integrityModified] > 0 || counts[integrityMissing] > 0 {
				return fmt.Errorf("integrity check failed")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	return cmd
}

func integrityIcon(s integrityStatus) string {
	switch s {
	case integrityOK:
		return "✓"
	case integrityModified:
		return "✗"
	case integrityMissing:
		return "✗"
	case integrityUnverifiable:
		return "?"
	default:
		return " "
	}
}
