package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/domehahn/skpm/internal/config"
	"github.com/domehahn/skpm/internal/publisher"
	"github.com/domehahn/skpm/internal/skill"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newPublishCmd() *cobra.Command {
	var (
		source    string
		tagFormat string
		noTag     bool
		noPush    bool
		outputDir string
	)

	cmd := &cobra.Command{
		Use:   "publish [path]",
		Short: "Validate, package, tag, and upload a skill to a registry",
		Long: `Runs the full release pipeline for a skill:

  1. Validate   — checks SKILL.md, VERSION, skill.yaml, CHANGELOG.md
  2. Package    — builds <name>-<version>.zip
  3. Git tag    — creates <name>/v<version> (or v<version> with --tag-format plain)
  4. Git push   — pushes the tag to origin
  5. Upload     — uploads the ZIP to the configured registry

Use --dry-run to preview all steps without making changes.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := outputFormat()

			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("path not found: %s", dir)}
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}

			src := source
			if src == "" {
				src = cfg.DefaultRegistry
			}
			if src == "" {
				return &UserError{Message: "no registry specified — use --source or set default_registry in config"}
			}

			// ── Step 1: Validate ───────────────────────────────────────
			printStep(cmd, "1/5", "Validating", dir)
			v := skill.NewValidatorWithOptions(skill.ValidationOptions{Publish: true})
			vResult, err := v.Validate(cmd.Context(), dir)
			if err != nil {
				return &InternalError{Message: "validate", Cause: err}
			}
			if !vResult.Valid {
				msgs := make([]string, len(vResult.Errors))
				for i, e := range vResult.Errors {
					msgs[i] = fmt.Sprintf("  [%s] %s", e.Field, e.Message)
				}
				return &UserError{Message: "validation failed:\n" + strings.Join(msgs, "\n")}
			}
			printOK(cmd, "Valid")

			// ── Step 2: Package ────────────────────────────────────────
			printStep(cmd, "2/5", "Packaging", "")
			if outputDir == "" {
				outputDir = os.TempDir()
			}
			p := skill.NewPackager()
			pkgResult, err := p.Package(cmd.Context(), dir, outputDir)
			if err != nil {
				return &InternalError{Message: "package", Cause: err}
			}
			printOK(cmd, fmt.Sprintf("%s (SHA256: %s)", pkgResult.OutputPath, pkgResult.SHA256[:12]+"…"))

			tag := formatGitTag(pkgResult.Name, pkgResult.Version, tagFormat)

			// ── Step 3: Git tag ────────────────────────────────────────
			if !noTag {
				printStep(cmd, "3/5", "Git tag", tag)
				if !globalDryRun {
					if err := gitTag(cmd.Context(), tag); err != nil {
						return &UserError{Message: fmt.Sprintf("git tag: %v\nIf the tag already exists, delete it first: git tag -d %s", err, tag)}
					}
				}
				printOK(cmd, tag)
			} else {
				printSkipped(cmd, "3/5", "Git tag (--no-tag)")
			}

			// ── Step 4: Git push ───────────────────────────────────────
			if !noTag && !noPush {
				printStep(cmd, "4/5", "Git push", "origin "+tag)
				if !globalDryRun {
					if err := gitPushTag(cmd.Context(), tag); err != nil {
						return &UserError{Message: fmt.Sprintf("git push tag: %v", err)}
					}
				}
				printOK(cmd, "pushed")
			} else {
				printSkipped(cmd, "4/5", "Git push (--no-push or --no-tag)")
			}

			// ── Step 5: Upload ─────────────────────────────────────────
			printStep(cmd, "5/5", "Uploading to", src)
			pub, err := publisher.New(src, tagFormat, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("publisher: %v", err)}
			}

			var pubResult *publisher.Result
			if !globalDryRun {
				pubResult, err = pub.Publish(cmd.Context(), pkgResult.Name, pkgResult.Version, pkgResult.OutputPath, pkgResult.SHA256)
				if err != nil {
					return &InternalError{Message: "upload", Cause: err}
				}
			} else {
				pubResult = &publisher.Result{
					Name:        pkgResult.Name,
					Version:     pkgResult.Version,
					SHA256:      pkgResult.SHA256,
					DownloadURL: "(dry run)",
					Tag:         tag,
				}
			}
			printOK(cmd, pubResult.DownloadURL)

			log.Debug().
				Str("name", pubResult.Name).
				Str("version", pubResult.Version).
				Str("url", pubResult.DownloadURL).
				Msg("published")

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "publish",
					Data: map[string]interface{}{
						"name":         pubResult.Name,
						"version":      pubResult.Version,
						"sha256":       pubResult.SHA256,
						"download_url": pubResult.DownloadURL,
						"tag":          tag,
						"registry":     src,
						"dry_run":      globalDryRun,
					},
				})
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout())
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run — no changes made.\n")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Published %s@%s\n", pubResult.Name, pubResult.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  Registry:     %s\n", src)
			fmt.Fprintf(cmd.OutOrStdout(), "  Download URL: %s\n", pubResult.DownloadURL)
			fmt.Fprintf(cmd.OutOrStdout(), "  SHA256:       %s\n", pubResult.SHA256)
			if !noTag {
				fmt.Fprintf(cmd.OutOrStdout(), "  Tag:          %s\n", tag)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintf(cmd.OutOrStdout(), "Add to agent-skills.lock:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  skpm add %s@%s --source %s\n", pubResult.Name, pubResult.Version, src)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry to publish to (uses default_registry if not set)")
	cmd.Flags().StringVar(&tagFormat, "tag-format", "prefixed", "Git tag format: prefixed (<name>/v<ver>) or plain (v<ver>)")
	cmd.Flags().BoolVar(&noTag, "no-tag", false, "Skip creating a git tag")
	cmd.Flags().BoolVar(&noPush, "no-push", false, "Skip pushing the git tag")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory for the intermediate ZIP (default: temp dir)")
	return cmd
}

// ── Git helpers ───────────────────────────────────────────────────────────

func gitTag(ctx context.Context, tag string) error {
	out, err := exec.CommandContext(ctx, "git", "tag", tag).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func gitPushTag(ctx context.Context, tag string) error {
	out, err := exec.CommandContext(ctx, "git", "push", "origin", tag).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func formatGitTag(name, version, tagFormat string) string {
	if tagFormat == "plain" {
		return "v" + version
	}
	return name + "/v" + version
}

// ── Output helpers ────────────────────────────────────────────────────────

func printStep(cmd *cobra.Command, step, action, detail string) {
	if detail != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s %s\n", step, action, detail)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s\n", step, action)
	}
}

func printOK(cmd *cobra.Command, detail string) {
	fmt.Fprintf(cmd.OutOrStdout(), "        ✓ %s\n", detail)
}

func printSkipped(cmd *cobra.Command, step, reason string) {
	fmt.Fprintf(cmd.OutOrStdout(), "  [%s] — %s\n", step, reason)
}
