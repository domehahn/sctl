package cli

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newSignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign or verify lockfile entries for supply-chain security",
		Long: `Signs each skill entry in the lockfile with HMAC-SHA256 using a shared key.
The signature covers: name, version, and sha256 digest.

Key sources (in priority order):
  1. --key-file <path>   path to a file containing a base64-encoded key
  2. SKPM_SIGNING_KEY    environment variable (base64-encoded)

Use 'skpm sign verify' to check signatures at install time.`,
	}
	cmd.AddCommand(newSignApplyCmd())
	cmd.AddCommand(newSignVerifyCmd())
	return cmd
}

func newSignApplyCmd() *cobra.Command {
	var keyFile string
	var keyEnv string
	var lockPath string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Compute and store HMAC-SHA256 signatures in the lockfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			key, err := LoadSigningKey(keyFile, keyEnv)
			if err != nil {
				return err
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}
			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills in lockfile.")
				return nil
			}

			for i, sl := range lf.Skills {
				lf.Skills[i].Signature = SkillSignature(key, sl.Name, sl.Version, sl.SHA256)
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would sign %d skill(s)\n", len(lf.Skills))
				return nil
			}
			if err := lf.Write(lockPath); err != nil {
				return &InternalError{Message: "write lockfile", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Signed %d skill(s) → %s\n", len(lf.Skills), lockPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&keyFile, "key-file", "", "Path to signing key file (base64-encoded)")
	cmd.Flags().StringVar(&keyEnv, "key-env", "SKPM_SIGNING_KEY", "Environment variable containing the signing key")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	return cmd
}

func newSignVerifyCmd() *cobra.Command {
	var keyFile string
	var keyEnv string
	var lockPath string

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify HMAC-SHA256 signatures stored in the lockfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			key, err := LoadSigningKey(keyFile, keyEnv)
			if err != nil {
				return err
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}
			if len(lf.Skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No skills in lockfile.")
				return nil
			}

			format := outputFormat()
			type result struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Status  string `json:"status"`
				Note    string `json:"note,omitempty"`
			}

			var results []result
			failures := 0
			for _, sl := range lf.Skills {
				r := result{Name: sl.Name, Version: sl.Version}
				if sl.Signature == "" {
					r.Status = "unsigned"
					r.Note = "no signature in lockfile — run: skpm sign apply"
					failures++
				} else {
					expected := SkillSignature(key, sl.Name, sl.Version, sl.SHA256)
					if hmac.Equal([]byte(sl.Signature), []byte(expected)) {
						r.Status = "ok"
					} else {
						r.Status = "invalid"
						r.Note = "signature mismatch — possible tampering"
						failures++
					}
				}
				results = append(results, r)
			}

			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: failures == 0,
					Command: "sign verify",
					Data:    results,
				})
				return nil
			}

			icons := map[string]string{"ok": "✓", "unsigned": "?", "invalid": "✗"}
			for _, r := range results {
				note := ""
				if r.Note != "" {
					note = "  (" + r.Note + ")"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s  %-30s  %s%s\n", icons[r.Status], r.Name+"@"+r.Version, r.Status, note)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d skill(s): %d ok, %d failed\n", len(results), len(results)-failures, failures)
			if failures > 0 {
				return &UserError{Message: "signature verification failed"}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&keyFile, "key-file", "", "Path to signing key file (base64-encoded)")
	cmd.Flags().StringVar(&keyEnv, "key-env", "SKPM_SIGNING_KEY", "Environment variable containing the signing key")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Path to lockfile (default: agent-skills.lock)")
	return cmd
}

// LoadSigningKey reads the HMAC key from a file or an environment variable.
func LoadSigningKey(keyFile, envName string) ([]byte, error) {
	var raw string
	if keyFile != "" {
		data, err := os.ReadFile(keyFile)
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("read key file %s: %v", keyFile, err)}
		}
		raw = strings.TrimSpace(string(data))
	} else {
		raw = os.Getenv(envName)
		if raw == "" {
			return nil, &UserError{Message: fmt.Sprintf("no signing key: provide --key-file or set %s (base64-encoded)", envName)}
		}
	}
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("decode signing key (expected base64): %v", err)}
	}
	if len(key) < 16 {
		return nil, &UserError{Message: "signing key is too short (minimum 16 bytes after decoding)"}
	}
	return key, nil
}

// SkillSignature computes HMAC-SHA256 over "name@version\nsha256:<digest>"
// and returns the result as a base64-encoded string.
func SkillSignature(key []byte, name, version, digest string) string {
	h := hmac.New(sha256.New, key)
	fmt.Fprintf(h, "%s@%s\nsha256:%s", name, version, digest)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
