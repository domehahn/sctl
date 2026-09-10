package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/spf13/cobra"
)

func newChangelogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Manage a skill's CHANGELOG.md",
	}
	cmd.AddCommand(newChangelogShowCmd())
	cmd.AddCommand(newChangelogAddCmd())
	return cmd
}

func newChangelogShowCmd() *cobra.Command {
	var count int

	cmd := &cobra.Command{
		Use:   "show [path]",
		Short: "Print changelog entries for a skill",
		Long:  `Prints the N most recent changelog entries from CHANGELOG.md. Defaults to 3.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			path := filepath.Join(dir, "CHANGELOG.md")
			data, err := os.ReadFile(path)
			if err != nil {
				if os.IsNotExist(err) {
					return &UserError{Message: fmt.Sprintf("CHANGELOG.md not found in %s", dir)}
				}
				return &InternalError{Message: "read CHANGELOG.md", Cause: err}
			}

			entries := parseChangelogEntries(string(data))
			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no changelog entries found)")
				return nil
			}

			if count > 0 && count < len(entries) {
				entries = entries[:count]
			}

			for i, e := range entries {
				if i > 0 {
					fmt.Fprintln(cmd.OutOrStdout())
				}
				fmt.Fprintln(cmd.OutOrStdout(), strings.TrimRight(e, "\n"))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&count, "count", "n", 3, "Number of entries to show")
	return cmd
}

func newChangelogAddCmd() *cobra.Command {
	var version string

	cmd := &cobra.Command{
		Use:   "add <message> [path]",
		Short: "Add a changelog entry for the current (or specified) version",
		Long: `Appends a bullet point to the changelog section for the current version.
The version is read from VERSION unless --version is given.

Unlike 'skpm version bump', this does not change the version number —
it just adds text to an existing section (or creates the section if absent).`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"changelog", "add"}, args...), func() error {
				message := args[0]
				dir := "."
				if len(args) == 2 {
					dir = args[1]
				}

				ver := version
				if ver == "" {
					v, err := skill.ReadVersion(dir)
					if err != nil {
						return &UserError{Message: fmt.Sprintf("read VERSION: %v (use --version to specify)", err)}
					}
					ver = v
				}
				// Normalise: strip leading v.
				ver = strings.TrimPrefix(ver, "v")

				if globalDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would add to %s/CHANGELOG.md [%s]: %s\n", dir, ver, message)
					return nil
				}

				if err := addChangelogMessage(dir, ver, message); err != nil {
					return err
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Added to CHANGELOG.md [%s]: %s\n", ver, message)
				return nil
			})
		},
	}

	cmd.Flags().StringVar(&version, "version", "", "Target version (defaults to VERSION file)")
	return cmd
}

// parseChangelogEntries splits CHANGELOG.md into per-version sections.
func parseChangelogEntries(content string) []string {
	re := regexp.MustCompile(`(?m)^##\s+`)
	indices := re.FindAllStringIndex(content, -1)
	if len(indices) == 0 {
		return nil
	}
	entries := make([]string, 0, len(indices))
	for i, idx := range indices {
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(content)
		}
		entries = append(entries, strings.TrimRight(content[idx[0]:end], "\n "))
	}
	return entries
}

// addChangelogMessage appends a bullet to the section for ver, creating it if absent.
func addChangelogMessage(dir, ver, message string) error {
	path := filepath.Join(dir, "CHANGELOG.md")

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return &InternalError{Message: "read CHANGELOG.md", Cause: err}
	}

	content := string(data)
	bullet := "- " + message

	// Try to find an existing section for this version.
	re := regexp.MustCompile(`(?m)^##\s+v?` + regexp.QuoteMeta(ver) + `\s*$`)
	loc := re.FindStringIndex(content)

	if loc != nil {
		// Find the end of this section (next ## heading or EOF).
		rest := content[loc[1]:]
		nextSection := regexp.MustCompile(`(?m)^##\s+`).FindStringIndex(rest)
		var insertAt int
		if nextSection != nil {
			// Insert before the next section, trimming trailing blank lines.
			chunk := strings.TrimRight(rest[:nextSection[0]], "\n")
			insertAt = loc[1] + len(chunk)
			content = content[:insertAt] + "\n" + bullet + "\n" + content[insertAt:]
		} else {
			chunk := strings.TrimRight(rest, "\n")
			insertAt = loc[1] + len(chunk)
			content = content[:insertAt] + "\n" + bullet + "\n"
		}
	} else {
		// Section doesn't exist — create it at the top (after the # heading if present).
		section := "## " + ver + "\n\n" + bullet + "\n\n"
		if strings.HasPrefix(content, "# ") {
			nl := strings.Index(content, "\n")
			content = content[:nl+1] + "\n" + section + strings.TrimLeft(content[nl+1:], "\n")
		} else if content == "" {
			content = "# Changelog\n\n" + section
		} else {
			content = section + content
		}
	}

	return os.WriteFile(path, []byte(content), 0o644)
}
