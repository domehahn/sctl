package cli

import (
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

const metaTagsKey = "tags"

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Add or remove custom tags on lockfile entries",
	}
	cmd.AddCommand(newTagAddCmd())
	cmd.AddCommand(newTagRemoveCmd())
	cmd.AddCommand(newTagListCmd())
	return cmd
}

func newTagAddCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "add <name> <tag>",
		Short: "Add a tag to a skill's lockfile metadata",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, tag := args[0], strings.TrimSpace(args[1])

			if tag == "" {
				return &UserError{Message: "tag must not be empty"}
			}
			if strings.Contains(tag, ",") {
				return &UserError{Message: fmt.Sprintf("tag %q must not contain commas", tag)}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			slPtr, ok := lf.Find(name)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", name)}
			}
			sl := *slPtr

			existing := skillTags(sl)
			for _, t := range existing {
				if t == tag {
					fmt.Fprintf(cmd.OutOrStdout(), "%s already has tag %q.\n", name, tag)
					return nil
				}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would add tag %q to %s\n", tag, name)
				return nil
			}

			if sl.Metadata == nil {
				sl.Metadata = map[string]string{}
			}
			if existing == nil {
				sl.Metadata[metaTagsKey] = tag
			} else {
				sl.Metadata[metaTagsKey] = joinTags(append(existing, tag))
			}

			lf.Upsert(sl)
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Tagged %s with %q\n", name, tag)
			return nil
		},
	}
	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}

func newTagRemoveCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "remove <name> <tag>",
		Short: "Remove a tag from a skill's lockfile metadata",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, tag := args[0], strings.TrimSpace(args[1])

			if tag == "" {
				return &UserError{Message: "tag must not be empty"}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			slPtr, ok := lf.Find(name)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", name)}
			}
			sl := *slPtr

			tags := skillTags(sl)
			var remaining []string
			found := false
			for _, t := range tags {
				if t == tag {
					found = true
				} else {
					remaining = append(remaining, t)
				}
			}
			if !found {
				fmt.Fprintf(cmd.OutOrStdout(), "%s does not have tag %q.\n", name, tag)
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would remove tag %q from %s\n", tag, name)
				return nil
			}

			if sl.Metadata == nil {
				sl.Metadata = map[string]string{}
			}
			if len(remaining) == 0 {
				delete(sl.Metadata, metaTagsKey)
			} else {
				sl.Metadata[metaTagsKey] = joinTags(remaining)
			}

			lf.Upsert(sl)
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed tag %q from %s\n", tag, name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}

func newTagListCmd() *cobra.Command {
	var lockPath string
	var filterTag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tags in the lockfile, or skills with a specific tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			if filterTag != "" {
				var matches []string
				for _, sl := range lf.Skills {
					for _, t := range skillTags(sl) {
						if t == filterTag {
							matches = append(matches, sl.Name)
							break
						}
					}
				}
				if len(matches) == 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "No skills tagged %q.\n", filterTag)
					return nil
				}
				for _, m := range matches {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", m)
				}
				return nil
			}

			// List all unique tags with counts.
			tagCount := map[string]int{}
			for _, sl := range lf.Skills {
				for _, t := range skillTags(sl) {
					tagCount[t]++
				}
			}
			if len(tagCount) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No tags found.")
				return nil
			}
			for t, count := range tagCount {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-24s  %d skill(s)\n", t, count)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&filterTag, "tag", "", "List only skills with this tag")
	return cmd
}

// skillTags parses the comma-separated tags stored in sl.Metadata["tags"].
func skillTags(sl lockfile.SkillLock) []string {
	if sl.Metadata == nil {
		return nil
	}
	raw := sl.Metadata[metaTagsKey]
	if raw == "" {
		return nil
	}
	var out []string
	for _, t := range splitTags(raw) {
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func splitTags(s string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := strings.TrimSpace(s[start:i])
			if part != "" {
				parts = append(parts, part)
			}
			start = i + 1
		}
	}
	return parts
}

func joinTags(tags []string) string {
	result := ""
	for i, t := range tags {
		if i > 0 {
			result += ","
		}
		result += t
	}
	return result
}
