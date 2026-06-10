package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type auditSeverity int

const (
	auditOK auditSeverity = iota
	auditInfo
	auditWarn
	auditError
)

func (s auditSeverity) String() string {
	switch s {
	case auditOK:
		return "ok"
	case auditInfo:
		return "info"
	case auditWarn:
		return "warn"
	case auditError:
		return "error"
	default:
		return "unknown"
	}
}

type auditFinding struct {
	Name     string        `json:"name"`
	Version  string        `json:"version"`
	Severity auditSeverity `json:"severity"`
	Issues   []string      `json:"issues,omitempty"`
}

func newAuditCmd() *cobra.Command {
	var lockPath string
	var source string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit installed skills for issues, deprecations, and outdated versions",
		Long: `Checks each locked skill against the registry and the local filesystem.

Findings:
  error  — version was yanked by the publisher
  warn   — version is deprecated, or a major update is available
  info   — a minor or patch update is available
  ok     — up to date and installed correctly

Registry discovery is attempted for each skill; skills whose registry does not
support discovery are checked only for local installation.`,
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

			workDir, _ := os.Getwd()
			findings := make([]auditFinding, 0, len(lf.Skills))

			for _, sl := range lf.Skills {
				f := auditFinding{Name: sl.Name, Version: sl.Version}

				// ── Registry checks ─────────────────────────────────
				src := sl.Source
				if source != "" {
					src = source
				}
				if src == "" {
					src = cfg.DefaultRegistry
				}
				if src != "" {
					reg, discovery, err := registryDiscovery(cmd.Context(), src, cfg)
					if err == nil {
						versions, err := discovery.ListVersions(cmd.Context(), registry.ParseSkillRef(sl.Name, registryDefaultNamespace(reg)))
						if err == nil {
							// Check if current version is yanked or deprecated.
							for _, v := range versions {
								if v.Version != sl.Version {
									continue
								}
								if v.Yanked {
									msg := "version was yanked"
									if v.YankReason != "" {
										msg += ": " + v.YankReason
									}
									f.Issues = append(f.Issues, msg)
									if f.Severity < auditError {
										f.Severity = auditError
									}
								} else if v.Deprecated {
									msg := "version is deprecated"
									if v.DeprecationReason != "" {
										msg += ": " + v.DeprecationReason
									}
									f.Issues = append(f.Issues, msg)
									if f.Severity < auditWarn {
										f.Severity = auditWarn
									}
								}
							}

							// Check for newer versions.
							latest := latestVersion(versions)
							if latest != "" && latest != sl.Version {
								delta := semverDelta(sl.Version, latest)
								switch delta {
								case "major":
									f.Issues = append(f.Issues, fmt.Sprintf("major update available: %s → %s", sl.Version, latest))
									if f.Severity < auditWarn {
										f.Severity = auditWarn
									}
								case "minor":
									f.Issues = append(f.Issues, fmt.Sprintf("minor update available: %s → %s", sl.Version, latest))
									if f.Severity < auditInfo {
										f.Severity = auditInfo
									}
								case "patch":
									f.Issues = append(f.Issues, fmt.Sprintf("patch update available: %s → %s", sl.Version, latest))
									if f.Severity < auditInfo {
										f.Severity = auditInfo
									}
								}
							}
						}
					}
				}

				// ── Local installation checks ────────────────────────
				missing := missingSkillMDs(workDir, sl.InstalledTo)
				if len(missing) > 0 {
					f.Issues = append(f.Issues, fmt.Sprintf("not installed in: %s", strings.Join(missing, ", ")))
					if f.Severity < auditWarn {
						f.Severity = auditWarn
					}
				}

				findings = append(findings, f)
			}

			sort.Slice(findings, func(i, j int) bool {
				if findings[i].Severity != findings[j].Severity {
					return findings[i].Severity > findings[j].Severity
				}
				return findings[i].Name < findings[j].Name
			})

			if format == OutputJSON {
				type jsonFinding struct {
					Name     string   `json:"name"`
					Version  string   `json:"version"`
					Severity string   `json:"severity"`
					Issues   []string `json:"issues,omitempty"`
				}
				out := make([]jsonFinding, len(findings))
				for i, f := range findings {
					out[i] = jsonFinding{Name: f.Name, Version: f.Version, Severity: f.Severity.String(), Issues: f.Issues}
				}
				PrintResult(format, CommandResult{Success: true, Command: "audit", Data: out})
				return nil
			}

			counts := map[auditSeverity]int{}
			for _, f := range findings {
				counts[f.Severity]++
			}

			for _, f := range findings {
				icon := auditIcon(f.Severity)
				fmt.Fprintf(cmd.OutOrStdout(), "  %s %-28s %s\n", icon, f.Name, f.Version)
				for _, issue := range f.Issues {
					fmt.Fprintf(cmd.OutOrStdout(), "      %s\n", issue)
				}
			}

			fmt.Fprintln(cmd.OutOrStdout())
			parts := []string{}
			if n := counts[auditError]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d error(s)", n))
			}
			if n := counts[auditWarn]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d warning(s)", n))
			}
			if n := counts[auditInfo]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d info", n))
			}
			if n := counts[auditOK]; n > 0 {
				parts = append(parts, fmt.Sprintf("%d ok", n))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Audit: %d skill(s) — %s\n", len(findings), strings.Join(parts, ", "))

			if counts[auditError] > 0 {
				return fmt.Errorf("audit found %d error(s)", counts[auditError])
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().StringVar(&source, "source", "", "Override registry source for all skills")
	return cmd
}

func auditIcon(s auditSeverity) string {
	switch s {
	case auditError:
		return "✗"
	case auditWarn:
		return "~"
	case auditInfo:
		return "↑"
	default:
		return "✓"
	}
}

// semverDelta returns the magnitude of the difference between two semver strings
// ("major", "minor", "patch", or "" if indeterminate).
func semverDelta(current, latest string) string {
	cv := "v" + current
	lv := "v" + latest
	if semver.Major(cv) != semver.Major(lv) {
		return "major"
	}
	if semver.MajorMinor(cv) != semver.MajorMinor(lv) {
		return "minor"
	}
	return "patch"
}

// missingSkillMDs returns the paths where SKILL.md does not exist.
func missingSkillMDs(workDir string, installedTo []string) []string {
	var missing []string
	for _, p := range installedTo {
		if _, err := os.Stat(filepath.Join(workDir, p, "SKILL.md")); err != nil {
			missing = append(missing, p)
		}
	}
	return missing
}
