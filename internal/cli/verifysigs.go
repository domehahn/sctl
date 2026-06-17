package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newVerifySignaturesCmd() *cobra.Command {
	var lockPath string
	var trustFile string
	var keyFile string
	var keyEnv string
	var failOnMissing bool

	cmd := &cobra.Command{
		Use:   "verify-signatures",
		Short: "Verify HMAC-SHA256 signatures for all signed skills in the lockfile",
		Long: `Re-derives the HMAC-SHA256 signature for every lockfile entry that carries a
signature field, using the key from --key-file or --key-env.

If --trust-file is provided, keys are looked up per-skill from the trust store.
Skills without a signature are reported but do not fail unless --fail-on-missing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &InternalError{Message: "read lockfile", Cause: err}
			}

			// Load a default key (used when trust store doesn't have a per-skill key).
			var defaultKey []byte
			if keyFile != "" {
				defaultKey, err = os.ReadFile(keyFile)
				if err != nil {
					return &UserError{Message: fmt.Sprintf("read key file: %v", err)}
				}
			} else if keyEnv != "" {
				v := os.Getenv(keyEnv)
				if v == "" {
					return &UserError{Message: fmt.Sprintf("env var %q is empty", keyEnv)}
				}
				defaultKey = []byte(v)
			}

			// Load trust store for per-skill key lookup.
			var ts *TrustStore
			if trustFile != "" {
				ts, err = LoadTrustStore(trustFile)
				if err != nil {
					return &UserError{Message: fmt.Sprintf("load trust store: %v", err)}
				}
			}

			type result struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Status  string `json:"status"`
				Reason  string `json:"reason,omitempty"`
			}

			var results []result
			passed, failed, skipped, missing := 0, 0, 0, 0

			for _, sl := range lf.Skills {
				if sl.Signature == "" {
					missing++
					results = append(results, result{
						Name:    sl.Name,
						Version: sl.Version,
						Status:  "unsigned",
					})
					continue
				}

				// Resolve key: trust store per-skill > default key.
				key := defaultKey
				if ts != nil {
					for _, entry := range ts.Keys {
						if entry.Name == sl.Name {
							key = []byte(entry.Key)
							break
						}
					}
				}

				if len(key) == 0 {
					skipped++
					results = append(results, result{
						Name:    sl.Name,
						Version: sl.Version,
						Status:  "skip",
						Reason:  "no key available",
					})
					continue
				}

				expected := SkillSignature(key, sl.Name, sl.Version, sl.SHA256)
				if expected == sl.Signature {
					passed++
					results = append(results, result{
						Name:    sl.Name,
						Version: sl.Version,
						Status:  "ok",
					})
				} else {
					failed++
					results = append(results, result{
						Name:    sl.Name,
						Version: sl.Version,
						Status:  "fail",
						Reason:  "signature mismatch",
					})
				}
			}

			format := outputFormat()
			if format == OutputJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(CommandResult{
					Success: failed == 0,
					Command: "verify-signatures",
					Data:    results,
				})
			}

			for _, r := range results {
				switch r.Status {
				case "ok":
					fmt.Fprintf(cmd.OutOrStdout(), "  ✓  %s@%s\n", r.Name, r.Version)
				case "fail":
					fmt.Fprintf(cmd.OutOrStdout(), "  ✗  %s@%s  %s\n", r.Name, r.Version, r.Reason)
				case "skip":
					fmt.Fprintf(cmd.OutOrStdout(), "  -  %s@%s  (skipped: %s)\n", r.Name, r.Version, r.Reason)
				case "unsigned":
					fmt.Fprintf(cmd.OutOrStdout(), "  ?  %s@%s  (unsigned)\n", r.Name, r.Version)
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n✓ %d  ✗ %d  - %d  ? %d\n", passed, failed, skipped, missing)

			if failed > 0 {
				return fmt.Errorf("%d signature(s) failed verification", failed)
			}
			if failOnMissing && missing > 0 {
				return fmt.Errorf("%d skill(s) have no signature", missing)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&trustFile, "trust-file", ".skpm-trust.yaml", "Path to trust store")
	cmd.Flags().StringVar(&keyFile, "key-file", "", "Path to HMAC key file (overrides trust store default)")
	cmd.Flags().StringVar(&keyEnv, "key-env", "", "Env var containing HMAC key")
	cmd.Flags().BoolVar(&failOnMissing, "fail-on-missing", false, "Exit non-zero when skills have no signature")
	return cmd
}
