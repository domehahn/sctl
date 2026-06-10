package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/publisher"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	var (
		bump        string
		message     string
		source      string
		tagFormat   string
		noTag       bool
		noPush      bool
		noPublish   bool
		outputDir   string
	)

	cmd := &cobra.Command{
		Use:   "release [path]",
		Short: "Bump version, update changelog, package, and publish in one step",
		Long: `Runs the full authoring release pipeline:

  1. Validate   — checks SKILL.md, VERSION, skill.yaml, CHANGELOG.md
  2. Bump       — increments the version (patch/minor/major)
  3. Changelog  — adds a section for the new version (message optional)
  4. Package    — builds <name>-<version>.zip
  5. Git tag    — creates <name>/v<version>
  6. Git push   — pushes the tag to origin
  7. Upload     — uploads the ZIP to the registry

Flags:
  --bump patch|minor|major   version component to increment (default: patch)
  --message "text"           changelog entry for the new version
  --source <registry>        registry to publish to
  --no-publish               skip steps 5–7 (local release only)
  --no-tag                   skip git tag
  --no-push                  skip git push

Use --dry-run to preview all steps without making changes.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return &UserError{Message: fmt.Sprintf("path not found: %s", dir)}
			}

			bump = strings.ToLower(bump)
			switch bump {
			case "patch", "minor", "major":
			default:
				return &UserError{Message: fmt.Sprintf("--bump must be patch, minor, or major (got %q)", bump)}
			}

			totalSteps := 7
			if noPublish {
				totalSteps = 4
			}
			stepFmt := func(n int) string { return fmt.Sprintf("%d/%d", n, totalSteps) }

			// ── Step 1: Validate ──────────────────────────────────────────
			printStep(cmd, stepFmt(1), "Validating", dir)
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

			// ── Step 2: Bump version ──────────────────────────────────────
			printStep(cmd, stepFmt(2), "Bumping version", "("+bump+")")
			state, err := readSkillVersionState(dir)
			if err != nil {
				return err
			}
			nextVersion, err := bumpStableVersion(state.Version, bump)
			if err != nil {
				return err
			}
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "        ✓ Would bump %s → %s\n", state.Version, nextVersion)
			} else {
				if err := writeSkillVersion(cmd, dir, state.Version, nextVersion, "release bump"); err != nil {
					return err
				}
			}
			printOK(cmd, fmt.Sprintf("%s → %s", state.Version, nextVersion))

			// ── Step 3: Changelog ─────────────────────────────────────────
			printStep(cmd, stepFmt(3), "Updating changelog", nextVersion)
			if !globalDryRun {
				if message != "" {
					if err := addChangelogMessage(dir, nextVersion, message); err != nil {
						return &InternalError{Message: "update changelog", Cause: err}
					}
				} else if err := ensureChangelogVersion(dir, nextVersion); err != nil {
					return &InternalError{Message: "ensure changelog section", Cause: err}
				}
			}
			changelogDetail := nextVersion
			if message != "" {
				changelogDetail = nextVersion + ": " + message
			}
			printOK(cmd, changelogDetail)

			// ── Step 4: Package ───────────────────────────────────────────
			printStep(cmd, stepFmt(4), "Packaging", "")
			if outputDir == "" {
				outputDir = os.TempDir()
			}
			p := skill.NewPackager()
			var pkgResult *skill.PackageResult
			if !globalDryRun {
				pkgResult, err = p.Package(cmd.Context(), dir, outputDir)
				if err != nil {
					return &InternalError{Message: "package", Cause: err}
				}
				printOK(cmd, fmt.Sprintf("%s (SHA256: %s)", pkgResult.OutputPath, pkgResult.SHA256[:12]+"…"))
			} else {
				// Read name for dry-run messaging.
				if syData, err := readSkillYAMLData(dir); err == nil {
					fmt.Fprintf(cmd.OutOrStdout(), "        ✓ Would package %s@%s\n", syData.Name, nextVersion)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "        ✓ Would package skill@%s\n", nextVersion)
				}
			}

			if noPublish {
				fmt.Fprintln(cmd.OutOrStdout())
				if globalDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "Dry run — no changes made.\n")
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Released %s (local only, --no-publish)\n", nextVersion)
				}
				return nil
			}

			tag := ""
			if pkgResult != nil {
				tag = formatGitTag(pkgResult.Name, pkgResult.Version, tagFormat)
			}

			// ── Step 5: Git tag ───────────────────────────────────────────
			if !noTag {
				printStep(cmd, stepFmt(5), "Git tag", tag)
				if !globalDryRun {
					if err := gitTag(cmd.Context(), tag); err != nil {
						return &UserError{Message: fmt.Sprintf("git tag: %v", err)}
					}
				}
				printOK(cmd, tag)
			} else {
				printSkipped(cmd, stepFmt(5), "Git tag (--no-tag)")
			}

			// ── Step 6: Git push ──────────────────────────────────────────
			if !noTag && !noPush {
				printStep(cmd, stepFmt(6), "Git push", "origin "+tag)
				if !globalDryRun {
					if err := gitPushTag(cmd.Context(), tag); err != nil {
						return &UserError{Message: fmt.Sprintf("git push tag: %v", err)}
					}
				}
				printOK(cmd, "pushed")
			} else {
				printSkipped(cmd, stepFmt(6), "Git push (--no-push or --no-tag)")
			}

			// ── Step 7: Upload ────────────────────────────────────────────
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

			printStep(cmd, stepFmt(7), "Uploading to", src)
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
				pubResult = &publisher.Result{DownloadURL: "(dry run)"}
			}
			printOK(cmd, pubResult.DownloadURL)

			fmt.Fprintln(cmd.OutOrStdout())
			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run — no changes made.\n")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Released %s@%s\n", pkgResult.Name, pkgResult.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  Registry:     %s\n", src)
			fmt.Fprintf(cmd.OutOrStdout(), "  Download URL: %s\n", pubResult.DownloadURL)
			fmt.Fprintf(cmd.OutOrStdout(), "  SHA256:       %s\n", pkgResult.SHA256)
			if !noTag {
				fmt.Fprintf(cmd.OutOrStdout(), "  Tag:          %s\n", tag)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&bump, "bump", "patch", "Version component to increment: patch, minor, or major")
	cmd.Flags().StringVar(&message, "message", "", "Changelog entry for the new version")
	cmd.Flags().StringVar(&source, "source", "", "Registry to publish to (uses default_registry if not set)")
	cmd.Flags().StringVar(&tagFormat, "tag-format", "prefixed", "Git tag format: prefixed (<name>/v<ver>) or plain (v<ver>)")
	cmd.Flags().BoolVar(&noTag, "no-tag", false, "Skip creating a git tag")
	cmd.Flags().BoolVar(&noPush, "no-push", false, "Skip pushing the git tag")
	cmd.Flags().BoolVar(&noPublish, "no-publish", false, "Stop after packaging (skip git tag, push, and upload)")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory for the intermediate ZIP (default: temp dir)")
	return cmd
}

// readSkillYAMLData reads the parsed skill.yaml from a directory.
func readSkillYAMLData(dir string) (*skill.SkillYAML, error) {
	return readSkillYAML(dir)
}
