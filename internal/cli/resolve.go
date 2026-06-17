package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newResolveCmd() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "resolve <name>",
		Short: "Show how a skill was resolved and where it is installed",
		Long: `Looks up a skill in the lockfile and prints full resolution details:
version, source registry, SHA-256 digest, install paths, signature, and provenance.

Useful for debugging lock conflicts or auditing exactly where a skill came from.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			sl, ok := lf.Find(name)
			if !ok {
				return &UserError{Message: fmt.Sprintf("skill %q not found in %s", name, lockPath)}
			}

			format := outputFormat()
			if format == OutputJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(CommandResult{
					Success: true,
					Command: "resolve",
					Data:    sl,
				})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "  name:          %s\n", sl.Name)
			if sl.Namespace != "" {
				fmt.Fprintf(w, "  namespace:     %s\n", sl.Namespace)
			}
			fmt.Fprintf(w, "  version:       %s\n", sl.Version)
			if sl.Constraint != "" {
				fmt.Fprintf(w, "  constraint:    %s\n", sl.Constraint)
			}
			fmt.Fprintf(w, "  source:        %s\n", sl.Source)
			if sl.RegistryURL != "" {
				fmt.Fprintf(w, "  registry_url:  %s\n", sl.RegistryURL)
			}
			if sl.RegistryType != "" {
				fmt.Fprintf(w, "  registry_type: %s\n", sl.RegistryType)
			}
			if sl.DownloadURL != "" {
				fmt.Fprintf(w, "  download_url:  %s\n", sl.DownloadURL)
			}
			sha := sl.SHA256
			if len(sha) > 16 {
				sha = sha[:16] + "..."
			}
			fmt.Fprintf(w, "  sha256:        %s\n", sha)
			if sl.PackageType != "" {
				fmt.Fprintf(w, "  package_type:  %s\n", sl.PackageType)
			}
			if sl.SourceCommit != "" {
				fmt.Fprintf(w, "  source_commit: %s\n", sl.SourceCommit)
			}
			if len(sl.InstalledTo) > 0 {
				fmt.Fprintf(w, "  installed_to:  %s\n", strings.Join(sl.InstalledTo, ", "))
			}
			if sl.Signature != "" {
				fmt.Fprintf(w, "  signature:     %s\n", sl.Signature)
			} else {
				fmt.Fprintf(w, "  signature:     (none)\n")
			}
			if sl.Provenance != "" {
				fmt.Fprintf(w, "  provenance:    %s\n", sl.Provenance)
			} else {
				fmt.Fprintf(w, "  provenance:    (none)\n")
			}
			if len(sl.CompatibleWith) > 0 {
				platforms := make([]string, len(sl.CompatibleWith))
				for i, p := range sl.CompatibleWith {
					platforms[i] = string(p)
				}
				fmt.Fprintf(w, "  compatible:    %s\n", strings.Join(platforms, ", "))
			}
			for k, v := range sl.Metadata {
				fmt.Fprintf(w, "  meta.%s: %s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	return cmd
}
