package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

type integrityCheck struct {
	Skill  string `json:"skill"`
	Path   string `json:"path"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func newCheckIntegrityCmd() *cobra.Command {
	var lockPath string
	var failOnMismatch bool

	cmd := &cobra.Command{
		Use:   "check-integrity",
		Short: "Verify installed skill files against lockfile SHA256 digests",
		Long: `For each locked skill, reads every file listed in installed_to and computes
its SHA256. Reports skills whose on-disk content has drifted from the locked digest.

A skill is checked by hashing its SKILL.md file in each installed_to path.
Use --fail-on-mismatch to exit non-zero if any mismatch is detected.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			var checks []integrityCheck
			mismatches := 0

			for _, sl := range lf.Skills {
				if sl.SHA256 == "" {
					checks = append(checks, integrityCheck{
						Skill:  sl.Name,
						Status: "skip",
						Detail: "no SHA256 in lockfile",
					})
					continue
				}

				if len(sl.InstalledTo) == 0 {
					checks = append(checks, integrityCheck{
						Skill:  sl.Name,
						Status: "skip",
						Detail: "not installed (no installed_to paths)",
					})
					continue
				}

				for _, installPath := range sl.InstalledTo {
					skillMD := filepath.Join(installPath, "SKILL.md")
					digest, err := fileSHA256(skillMD)
					if err != nil {
						checks = append(checks, integrityCheck{
							Skill:  sl.Name,
							Path:   skillMD,
							Status: "missing",
							Detail: err.Error(),
						})
						mismatches++
						continue
					}

					if digest != sl.SHA256 {
						checks = append(checks, integrityCheck{
							Skill:  sl.Name,
							Path:   skillMD,
							Status: "mismatch",
							Detail: fmt.Sprintf("expected %s, got %s", sl.SHA256[:12]+"...", digest[:12]+"..."),
						})
						mismatches++
					} else {
						checks = append(checks, integrityCheck{
							Skill:  sl.Name,
							Path:   skillMD,
							Status: "ok",
						})
					}
				}
			}

			format := outputFormat()
			if format == OutputJSON {
				if err := json.NewEncoder(cmd.OutOrStdout()).Encode(checks); err != nil {
					return err
				}
			} else {
				for _, c := range checks {
					switch c.Status {
					case "ok":
						fmt.Fprintf(cmd.OutOrStdout(), "  ok       %s\n", c.Skill)
					case "skip":
						fmt.Fprintf(cmd.OutOrStdout(), "  skip     %s  (%s)\n", c.Skill, c.Detail)
					case "missing":
						fmt.Fprintf(cmd.OutOrStdout(), "  missing  %s  %s\n", c.Skill, c.Path)
					case "mismatch":
						fmt.Fprintf(cmd.OutOrStdout(), "  MISMATCH %s  %s  %s\n", c.Skill, c.Path, c.Detail)
					}
				}
				if mismatches > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "\n%d integrity issue(s) detected.\n", mismatches)
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), "\nAll checked skills OK.")
				}
			}

			if failOnMismatch && mismatches > 0 {
				return fmt.Errorf("%d integrity issue(s) detected", mismatches)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().BoolVar(&failOnMismatch, "fail-on-mismatch", false, "Exit non-zero if any mismatch is found")
	return cmd
}

// fileSHA256 computes the SHA256 hex digest of a file.
func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
