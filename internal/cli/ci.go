package cli

import (
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/admission"
	"github.com/domehahn/skpm/v2/internal/cache"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/installer"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newCICmd() *cobra.Command {
	var lockPath string
	var platform string
	var target string
	var noVerify bool
	var production bool
	var requireAdmission bool
	var environment string

	cmd := &cobra.Command{
		Use:   "ci",
		Short: "Frozen install for CI: strict lockfile check, prune, and verify in one step",
		Long: `Like 'skpm install --frozen-lockfile --prune', but stricter and optimized
for CI pipelines:

  • Requires agent-skills.lock to already exist (never generates one).
  • Fails immediately if agent-skills.yaml and agent-skills.lock are inconsistent.
  • Always prunes skills no longer in the lockfile.
  • Runs 'skpm verify' after installation to confirm integrity.
  • Support --production mode enforcing frozen lockfile + skgate admission + production environment.
  • No interactive prompts. Exits non-zero on any issue.

This mirrors the behavior of 'npm ci' / 'cargo fetch --locked'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			if production {
				requireAdmission = true
				environment = "production"
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}

			// Step 1: lockfile must exist — never generate from manifest.
			if _, err := os.Stat(lockPath); err != nil {
				if os.IsNotExist(err) {
					return &UserError{Message: fmt.Sprintf(
						"%s not found.\n\nCI requires a committed lockfile. Run 'skpm install' locally, commit %s, and retry.",
						lockPath, lockPath,
					)}
				}
				return &InternalError{Message: "stat lockfile", Cause: err}
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			// Step 2: manifest-lockfile consistency check.
			if _, manifestErr := os.Stat(manifest.DefaultFilename); manifestErr == nil {
				outdated, checkErr := lockfileOutdated(cmd.Context(), cfg, manifest.DefaultFilename, lockPath)
				if checkErr != nil {
					return &UserError{Message: fmt.Sprintf("consistency check: %v", checkErr)}
				}
				if outdated {
					return &UserError{Message: fmt.Sprintf(
						"%s is out of sync with %s.\n\nRun 'skpm install' locally, commit the updated %s, and retry.",
						lockPath, manifest.DefaultFilename, lockPath,
					)}
				}
			}

			// Step 3: load lockfile and install.
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}
			if platform != "" {
				filterLockfileForPlatform(lf, platform)
			}
			fillMissingInstallPaths(lf)

			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to install.")
				return nil
			}

			if err := runProjectHook(cmd, "pre_install"); err != nil {
				return err
			}

			admClient := admission.NewClientFromEnv()
			if requireAdmission || environment == "production" {
				if admClient == nil || admClient.URL == "" {
					return &AdmissionError{Message: "admission evidence required in production or --require-admission mode, but SKPM_ADMISSION_URL is not set"}
				}
			}

			if admClient != nil {
				for _, sl := range lf.Skills {
					req := admission.AdmissionRequest{
						Name:           sl.Name,
						Version:        sl.Version,
						PackageDigest:  sl.SHA256,
						ArtifactDigest: sl.Artifact,
						Source:         sl.Source,
						Registry:       sl.Source,
						Action:         "install",
						Environment:    environment,
					}
					dec, err := admClient.Evaluate(cmd.Context(), req)
					if err != nil {
						if requireAdmission || environment == "production" || admClient.Enforce {
							return &AdmissionError{Message: fmt.Sprintf("admission evaluation failed for %s@%s: %v", sl.Name, sl.Version, err)}
						}
					}
					if dec != nil && dec.Decision != admission.DecisionAllow {
						if requireAdmission || environment == "production" || admClient.Enforce {
							return &AdmissionError{Message: fmt.Sprintf("admission check rejected action install for %s@%s: decision=%s (%s)", sl.Name, sl.Version, dec.Decision, dec.Reason)}
						}
						log.Warn().Str("skill", sl.Name).Str("decision", string(dec.Decision)).Str("reason", dec.Reason).Msg("admission policy advisory")
					}
				}
			}

			c := cache.New(cfg.CacheDir)
			ins := installer.New(c)

			workDir, _ := os.Getwd()
			if target != "" {
				workDir = target
			}
			opts := installer.Options{
				DryRun:      globalDryRun,
				Concurrency: globalConcurrency,
				WorkDir:     workDir,
			}

			log.Debug().Str("lockfile", lockPath).Int("skills", len(lf.Skills)).Msg("ci install")

			result, err := ins.Install(cmd.Context(), lf, opts)
			if err != nil {
				return &InternalError{Message: "install failed", Cause: err}
			}

			// Step 4: prune skills no longer in lockfile.
			if !globalDryRun {
				if err := pruneInstallPaths(workDir, lf); err != nil {
					return &InternalError{Message: "prune installed skills", Cause: err}
				}
			}

			_ = runProjectHook(cmd, "post_install")

			total := len(result.Installed) + len(result.FromCache)
			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "ci",
					Data: map[string]interface{}{
						"installed":  result.Installed,
						"from_cache": result.FromCache,
						"skipped":    result.Skipped,
						"dry_run":    globalDryRun,
					},
				})
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Installed %d skill(s)", total)
				if len(result.FromCache) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), " (%d from cache)", len(result.FromCache))
				}
				fmt.Fprintln(cmd.OutOrStdout())
				for _, name := range append(result.Installed, result.FromCache...) {
					sl, _ := lf.Find(name)
					if sl != nil {
						fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s@%s\n", name, sl.Version)
					}
				}
			}

			// Step 5: integrity verification (unless --no-verify).
			if !noVerify && !globalDryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "\nVerifying...")
				var verifyErrors []string
				for _, sl := range lf.Skills {
					for _, p := range sl.InstalledTo {
						skillMD := p + "/SKILL.md"
						if _, statErr := os.Stat(skillMD); statErr != nil {
							verifyErrors = append(verifyErrors, fmt.Sprintf("%s: %s missing", sl.Name, skillMD))
						}
					}
				}
				if len(verifyErrors) > 0 {
					return &UserError{Message: fmt.Sprintf("post-install verification failed:\n  %s", stringsJoin(verifyErrors, "\n  "))}
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Verification passed.")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	cmd.Flags().StringVar(&platform, "platform", "", "Install only skills compatible with a platform")
	cmd.Flags().StringVar(&target, "target", "", "Install into a target directory")
	cmd.Flags().BoolVar(&noVerify, "no-verify", false, "Skip post-install integrity check")
	cmd.Flags().BoolVar(&production, "production", false, "Enforce production mode (frozen lockfile + skgate admission + production environment)")
	cmd.Flags().BoolVar(&requireAdmission, "require-admission", false, "Require skgate admission decision ALLOW before installing")
	cmd.Flags().StringVar(&environment, "environment", "development", "Target environment for install (e.g. development, production)")
	return cmd
}
