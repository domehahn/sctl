package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newReportCmd() *cobra.Command {
	var lockPath string
	var policyFile string
	var keyFile string
	var keyEnv string
	var failOnIssue bool

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Cross-cutting health report for all installed skills",
		Long: `Runs all available checks in a single pass and prints a consolidated report:

  [LOCKFILE]    file presence and parse
  [SIGNATURES]  HMAC-SHA256 verify (skipped if no key provided)
  [POLICY]      skpm-policy.yaml rule gates (skipped if no policy file)
  [LICENSES]    license field presence per skill

Use --fail to exit with code 1 when any issue is found — suitable for CI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}

			var sections []reportSection

			// ── LOCKFILE ───────────────────────────────────────────────────────
			{
				var findings []reportFinding
				lf, err := lockfile.Read(lockPath)
				if err != nil {
					findings = append(findings, reportFinding{level: "error", message: fmt.Sprintf("cannot read %s: %v", lockPath, err)})
					sections = append(sections, reportSection{"LOCKFILE", findings})
					printReport(cmd, sections, failOnIssue)
					return nil
				}
				if len(lf.Skills) == 0 {
					findings = append(findings, reportFinding{level: "warn", message: "lockfile is empty"})
				} else {
					findings = append(findings, reportFinding{level: "ok", message: fmt.Sprintf("%d skill(s) in lockfile", len(lf.Skills))})
				}
				sections = append(sections, reportSection{"LOCKFILE", findings})

				// ── SIGNATURES ────────────────────────────────────────────────
				{
					var sigFindings []reportFinding
					if keyFile != "" || keyEnv != "" {
						envName := keyEnv
						if envName == "" {
							envName = "SKPM_SIGNING_KEY"
						}
						key, keyErr := LoadSigningKey(keyFile, envName)
						if keyErr != nil {
							sigFindings = append(sigFindings, reportFinding{level: "error", message: fmt.Sprintf("load signing key: %v", keyErr)})
						} else {
							unsigned, invalid := 0, 0
							for _, sl := range lf.Skills {
								if sl.Signature == "" {
									unsigned++
								} else {
									expected := SkillSignature(key, sl.Name, sl.Version, sl.SHA256)
									if sl.Signature != expected {
										invalid++
									}
								}
							}
							if unsigned > 0 {
								sigFindings = append(sigFindings, reportFinding{level: "warn", message: fmt.Sprintf("%d skill(s) unsigned — run: skpm sign apply", unsigned)})
							}
							if invalid > 0 {
								sigFindings = append(sigFindings, reportFinding{level: "error", message: fmt.Sprintf("%d skill(s) have invalid signatures — possible tampering", invalid)})
							}
							if unsigned == 0 && invalid == 0 {
								sigFindings = append(sigFindings, reportFinding{level: "ok", message: "all signatures valid"})
							}
						}
					} else {
						sigFindings = append(sigFindings, reportFinding{level: "skip", message: "provide --key-file or --key-env to verify signatures"})
					}
					sections = append(sections, reportSection{"SIGNATURES", sigFindings})
				}

				// ── POLICY ────────────────────────────────────────────────────
				{
					var polFindings []reportFinding
					pf := policyFile
					if pf == "" {
						pf = defaultPolicyFile
					}
					if _, statErr := os.Stat(pf); os.IsNotExist(statErr) {
						polFindings = append(polFindings, reportFinding{level: "skip", message: fmt.Sprintf("no policy file found (%s) — run: skpm policy init", pf)})
					} else {
						p, loadErr := loadPolicyFile(pf)
						if loadErr != nil {
							polFindings = append(polFindings, reportFinding{level: "error", message: fmt.Sprintf("load policy: %v", loadErr)})
						} else {
							violations := CheckPolicy(p, lf)
							if len(violations) == 0 {
								polFindings = append(polFindings, reportFinding{level: "ok", message: "no policy violations"})
							} else {
								for _, v := range violations {
									polFindings = append(polFindings, reportFinding{level: "error", message: fmt.Sprintf("[%s] %s: %s", v.Rule, v.Skill, v.Message)})
								}
							}
						}
					}
					sections = append(sections, reportSection{"POLICY", polFindings})
				}

				// ── LICENSES ──────────────────────────────────────────────────
				{
					var licFindings []reportFinding
					entries := collectNoticeEntries(lf)
					missing := 0
					for _, e := range entries {
						if e.License == "" {
							missing++
							licFindings = append(licFindings, reportFinding{level: "warn", message: fmt.Sprintf("%s@%s has no license declared", e.Name, e.Version)})
						}
					}
					if missing == 0 {
						licFindings = append(licFindings, reportFinding{level: "ok", message: "all skills declare a license"})
					}
					sections = append(sections, reportSection{"LICENSES", licFindings})
				}
			}

			printReport(cmd, sections, failOnIssue)
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Lockfile path (default: agent-skills.lock)")
	cmd.Flags().StringVar(&policyFile, "policy", "", "Policy file (default: skpm-policy.yaml if present)")
	cmd.Flags().StringVar(&keyFile, "key-file", "", "Signing key file for signature verification")
	cmd.Flags().StringVar(&keyEnv, "key-env", "", "Env var with signing key (e.g. SKPM_SIGNING_KEY)")
	cmd.Flags().BoolVar(&failOnIssue, "fail", false, "Exit with code 1 when any error or warning is found")
	return cmd
}

type reportFinding struct {
	level   string // ok, warn, error, skip
	message string
}

type reportSection struct {
	name     string
	findings []reportFinding
}

func printReport(cmd *cobra.Command, sections []reportSection, failOnIssue bool) {
	icons := map[string]string{
		"ok":    "✓",
		"warn":  "!",
		"error": "✗",
		"skip":  "-",
	}

	totalErrors, totalWarns := 0, 0

	for _, sec := range sections {
		fmt.Fprintf(cmd.OutOrStdout(), "\n[%s]\n", sec.name)
		for _, f := range sec.findings {
			icon := icons[f.level]
			if icon == "" {
				icon = "?"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s\n", icon, f.message)
			if f.level == "error" {
				totalErrors++
			} else if f.level == "warn" {
				totalWarns++
			}
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\n")
	if totalErrors == 0 && totalWarns == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Report: all checks passed\n")
	} else {
		parts := []string{}
		if totalErrors > 0 {
			parts = append(parts, fmt.Sprintf("%d error(s)", totalErrors))
		}
		if totalWarns > 0 {
			parts = append(parts, fmt.Sprintf("%d warning(s)", totalWarns))
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Report: %s\n", strings.Join(parts, ", "))
	}

	if failOnIssue && (totalErrors > 0 || totalWarns > 0) {
		// Return error via the command context — caller must check.
		_ = cmd.RunE // signal to caller; actual exit is via RunE returning error
		// We rely on the caller checking the output; for programmatic use,
		// callers should use --fail with CI pipelines where exit code matters.
		os.Exit(1)
	}
}
