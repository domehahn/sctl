package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newSBOMCmd() *cobra.Command {
	var format string
	var outFile string
	var lockPath string

	cmd := &cobra.Command{
		Use:   "sbom",
		Short: "Generate a Software Bill of Materials for installed skills",
		Long: `Reads the lockfile and emits an SBOM in SPDX 2.3 or CycloneDX 1.5 format.

Supported formats:
  spdx       SPDX 2.3 tag-value text (default)
  cyclonedx  CycloneDX 1.5 JSON`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			var out string
			switch strings.ToLower(format) {
			case "cyclonedx", "cdx":
				out, err = SbomCycloneDX(lf)
			default:
				out, err = SbomSPDX(lf)
			}
			if err != nil {
				return &InternalError{Message: "generate SBOM", Cause: err}
			}

			if outFile != "" && !globalDryRun {
				if writeErr := os.WriteFile(outFile, []byte(out), 0o644); writeErr != nil {
					return &InternalError{Message: "write SBOM file", Cause: writeErr}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "SBOM written to %s (%d component(s))\n", outFile, len(lf.Skills))
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), out)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "spdx", "SBOM format: spdx or cyclonedx")
	cmd.Flags().StringVarP(&outFile, "out", "f", "", "Write SBOM to file instead of stdout")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	return cmd
}

func SbomSPDX(lf *lockfile.LockFile) (string, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var sb strings.Builder

	sb.WriteString("SPDXVersion: SPDX-2.3\n")
	sb.WriteString("DataLicense: CC0-1.0\n")
	sb.WriteString("SPDXID: SPDXRef-DOCUMENT\n")
	sb.WriteString("DocumentName: agent-skill-sbom\n")
	sb.WriteString(fmt.Sprintf("DocumentNamespace: https://skpm.io/sbom/%s\n", now))
	sb.WriteString(fmt.Sprintf("Created: %s\n", now))
	sb.WriteString("Creator: Tool: skpm\n")
	sb.WriteString("\n")

	for _, sl := range lf.Skills {
		spdxID := "SPDXRef-Package-" + strings.ReplaceAll(sl.Name, "/", "-")
		sb.WriteString(fmt.Sprintf("PackageName: %s\n", sl.Name))
		sb.WriteString(fmt.Sprintf("SPDXID: %s\n", spdxID))
		sb.WriteString(fmt.Sprintf("PackageVersion: %s\n", sl.Version))
		if sl.DownloadURL != "" {
			sb.WriteString(fmt.Sprintf("PackageDownloadLocation: %s\n", sl.DownloadURL))
		} else {
			sb.WriteString("PackageDownloadLocation: NOASSERTION\n")
		}
		sb.WriteString("FilesAnalyzed: false\n")
		if sl.SHA256 != "" {
			sb.WriteString(fmt.Sprintf("PackageChecksum: SHA256: %s\n", sl.SHA256))
		}
		purl := fmt.Sprintf("pkg:skpm/%s@%s", sl.Name, sl.Version)
		if sl.Namespace != "" {
			purl = fmt.Sprintf("pkg:skpm/%s/%s@%s", sl.Namespace, sl.Name, sl.Version)
		}
		sb.WriteString(fmt.Sprintf("ExternalRef: PACKAGE-MANAGER purl %s\n", purl))
		if sl.Signature != "" {
			sb.WriteString(fmt.Sprintf("PackageComment: signature:%s\n", sl.Signature))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func SbomCycloneDX(lf *lockfile.LockFile) (string, error) {
	type hash struct {
		Alg     string `json:"alg"`
		Content string `json:"content"`
	}
	type component struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Version string `json:"version"`
		PURL    string `json:"purl"`
		Hashes  []hash `json:"hashes,omitempty"`
	}
	type metadata struct {
		Timestamp string `json:"timestamp"`
		Tools     []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"tools"`
	}
	type bom struct {
		BOMFormat   string      `json:"bomFormat"`
		SpecVersion string      `json:"specVersion"`
		Version     int         `json:"version"`
		Metadata    metadata    `json:"metadata"`
		Components  []component `json:"components"`
	}

	doc := bom{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.5",
		Version:     1,
		Metadata: metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Tools: []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			}{{Name: "skpm", Version: Version}},
		},
	}

	for _, sl := range lf.Skills {
		purl := fmt.Sprintf("pkg:skpm/%s@%s", sl.Name, sl.Version)
		if sl.Namespace != "" {
			purl = fmt.Sprintf("pkg:skpm/%s/%s@%s", sl.Namespace, sl.Name, sl.Version)
		}
		c := component{
			Type:    "library",
			Name:    sl.Name,
			Version: sl.Version,
			PURL:    purl,
		}
		if sl.SHA256 != "" {
			c.Hashes = []hash{{Alg: "SHA-256", Content: sl.SHA256}}
		}
		doc.Components = append(doc.Components, c)
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}
