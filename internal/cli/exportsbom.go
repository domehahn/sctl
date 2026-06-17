package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type sbomComponent struct {
	Name      string   `json:"name"`
	Version   string   `json:"version,omitempty"`
	Source    string   `json:"source,omitempty"`
	SHA256    string   `json:"sha256,omitempty"`
	Signature string   `json:"signature,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
}

type sbomDocument struct {
	Format      string          `json:"format"`
	Version     string          `json:"version"`
	GeneratedAt string          `json:"generated_at"`
	Components  []sbomComponent `json:"components"`
}

func newExportSBOMCmd() *cobra.Command {
	var lockPath string
	var out string

	cmd := &cobra.Command{
		Use:   "export-sbom",
		Short: "Export a software bill-of-materials JSON for all locked skills",
		Long: `Serialises every entry in the lockfile into a machine-readable SBOM document
(name, version, source, SHA256, signature, tags, namespace) following a minimal
custom format.

Write to a file with --out or print to stdout with --out -.

Useful for compliance audits, supply-chain reviews, and CI artefact storage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			components := make([]sbomComponent, 0, len(lf.Skills))
			for _, sl := range lf.Skills {
				c := sbomComponent{
					Name:      sl.Name,
					Version:   sl.Version,
					Source:    sl.Source,
					SHA256:    sl.SHA256,
					Signature: sl.Signature,
					Namespace: sl.Namespace,
				}
				if tags := skillTags(sl); len(tags) > 0 {
					c.Tags = tags
				}
				components = append(components, c)
			}

			doc := sbomDocument{
				Format:      "agentskills-sbom",
				Version:     "1",
				GeneratedAt: time.Now().UTC().Format(time.RFC3339),
				Components:  components,
			}

			encode := func(w interface{ Write([]byte) (int, error) }) error {
				enc := json.NewEncoder(w)
				enc.SetIndent("", "  ")
				return enc.Encode(doc)
			}

			if out == "" || out == "-" {
				if err := encode(cmd.OutOrStdout()); err != nil {
					return err
				}
				if out == "-" {
					return nil
				}
				return nil
			}

			f, err := os.Create(out)
			if err != nil {
				return fmt.Errorf("create %s: %w", out, err)
			}
			if err := encode(f); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("close %s: %w", out, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Wrote SBOM: %s (%d component(s))\n", out, len(components))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&out, "out", "-", "Output path; use '-' to print to stdout")
	return cmd
}
