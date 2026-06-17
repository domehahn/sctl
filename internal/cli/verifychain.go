package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const keyringFile = ".skpm-keys"

type chainResult struct {
	Skill  string `json:"skill"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func newVerifyChainCmd() *cobra.Command {
	var lockPath string
	var keysFile string
	var keyEnv string
	var all bool

	cmd := &cobra.Command{
		Use:   "verify-chain [name]",
		Short: "Verify a skill's HMAC signature against trusted keys in a keyring",
		Long: `Loads trusted keys from ` + keyringFile + ` (or --keys-file) and verifies each
skill's Signature field. Accepts if any key in the keyring produces a matching
signature. Use --all to check every skill in the lockfile at once.

The keyring file is a YAML list of base64-encoded HMAC keys, or plain text
with one base64-encoded key per line.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return &UserError{Message: "provide a skill name or use --all"}
			}

			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read lockfile: %v", err)}
			}

			keys, err := loadKeyring(keysFile, keyEnv)
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return &UserError{Message: fmt.Sprintf("no trusted keys found in %s (use --keys-file or --key-env)", keysFile)}
			}

			var toCheck []lockfile.SkillLock
			if all {
				toCheck = lf.Skills
				if len(toCheck) == 0 {
					return &UserError{Message: "lockfile is empty — nothing to verify"}
				}
			} else {
				sl, ok := lf.Find(args[0])
				if !ok {
					return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", args[0])}
				}
				toCheck = []lockfile.SkillLock{*sl}
			}

			results := make([]chainResult, 0, len(toCheck))
			failures := 0
			for _, sl := range toCheck {
				r := verifySkillChain(sl, keys)
				results = append(results, r)
				if r.Status != "ok" {
					failures++
				}
			}

			if strings.ToLower(globalOutput) == "json" {
				if err := json.NewEncoder(cmd.OutOrStdout()).Encode(results); err != nil {
					return err
				}
				if failures > 0 {
					return fmt.Errorf("%d skill(s) failed verification", failures)
				}
				return nil
			}

			for _, r := range results {
				if r.Status == "ok" {
					fmt.Fprintf(cmd.OutOrStdout(), "  ok      %s\n", r.Skill)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "  %-7s %s  — %s\n", r.Status, r.Skill, r.Detail)
				}
			}

			if failures > 0 {
				return fmt.Errorf("%d skill(s) failed verification", failures)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nAll %d skill(s) verified.\n", len(results))
			return nil
		},
	}

	cmd.Flags().StringVarP(&lockPath, "lock", "l", "agent-skills.lock.yaml", "Path to lockfile")
	cmd.Flags().StringVar(&keysFile, "keys-file", keyringFile, "Path to keyring file")
	cmd.Flags().StringVar(&keyEnv, "key-env", "", "Env var containing a single trusted key (base64); appended to keyring")
	cmd.Flags().BoolVar(&all, "all", false, "Verify every skill in the lockfile")
	return cmd
}

func verifySkillChain(sl lockfile.SkillLock, keys [][]byte) chainResult {
	if sl.Signature == "" {
		return chainResult{Skill: sl.Name, Status: "unsigned", Detail: "no signature recorded"}
	}
	for _, key := range keys {
		if sl.Signature == SkillSignature(key, sl.Name, sl.Version, sl.SHA256) {
			return chainResult{Skill: sl.Name, Status: "ok"}
		}
	}
	return chainResult{Skill: sl.Name, Status: "invalid", Detail: "signature does not match any trusted key"}
}

// loadKeyring reads base64-encoded HMAC keys from a file and optionally from an env var.
// The file may be a YAML list of strings or plain newline-separated lines.
func loadKeyring(path, envName string) ([][]byte, error) {
	var rawKeys []string

	if envName != "" {
		if v := os.Getenv(envName); v != "" {
			rawKeys = append(rawKeys, strings.TrimSpace(v))
		}
	}

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, &InternalError{Message: "read keyring", Cause: err}
	}
	if err == nil {
		var yamlList []string
		if yaml.Unmarshal(data, &yamlList) == nil && len(yamlList) > 0 {
			rawKeys = append(rawKeys, yamlList...)
		} else {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					rawKeys = append(rawKeys, line)
				}
			}
		}
	}

	keys := make([][]byte, 0, len(rawKeys))
	for _, raw := range rawKeys {
		key, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return nil, &UserError{Message: fmt.Sprintf("invalid base64 key in keyring: %v", err)}
		}
		if len(key) < 16 {
			return nil, &UserError{Message: "keyring key is too short (minimum 16 bytes decoded)"}
		}
		keys = append(keys, key)
	}
	return keys, nil
}
