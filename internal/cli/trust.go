package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultTrustFile = ".skpm-trust.yaml"

// TrustEntry is one trusted signer key in the trust store.
type TrustEntry struct {
	Name  string `yaml:"name"`
	Key   string `yaml:"key"` // base64-encoded
	Added string `yaml:"added"`
}

// TrustStore is the full contents of .skpm-trust.yaml.
type TrustStore struct {
	Version int          `yaml:"version"`
	Keys    []TrustEntry `yaml:"keys"`
}

// LoadTrustStore reads and parses the trust store file.
func LoadTrustStore(path string) (*TrustStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &TrustStore{Version: 1, Keys: nil}, nil
		}
		return nil, &UserError{Message: fmt.Sprintf("read trust store %s: %v", path, err)}
	}
	var ts TrustStore
	if err := yaml.Unmarshal(data, &ts); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse trust store %s: %v", path, err)}
	}
	if ts.Version == 0 {
		ts.Version = 1
	}
	return &ts, nil
}

func saveTrustStore(path string, ts *TrustStore) error {
	data, err := yaml.Marshal(ts)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func newTrustCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trust",
		Short: "Manage the trusted signing-key registry",
		Long: `Maintains a project-level keyring (.skpm-trust.yaml) of trusted signing keys.
Each entry has a name, a base64-encoded key, and the date it was added.

Use 'skpm sign verify' with --key-file or --key-env to verify individual entries.
The trust store is intended for auditing which keys are authorized in your project.`,
	}
	cmd.AddCommand(newTrustListCmd())
	cmd.AddCommand(newTrustAddCmd())
	cmd.AddCommand(newTrustRevokeCmd())
	return cmd
}

func newTrustListCmd() *cobra.Command {
	var trustFile string

	return &cobra.Command{
		Use:   "list",
		Short: "List all trusted signing keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			if trustFile == "" {
				trustFile = defaultTrustFile
			}
			ts, err := LoadTrustStore(trustFile)
			if err != nil {
				return err
			}
			if len(ts.Keys) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No trusted keys configured.")
				return nil
			}

			format := outputFormat()
			if format == OutputJSON {
				PrintResult(format, CommandResult{Success: true, Command: "trust list", Data: ts.Keys})
				return nil
			}

			for _, e := range ts.Keys {
				// Show only first 12 chars of key for auditing without full exposure.
				preview := e.Key
				if len(preview) > 12 {
					preview = preview[:12] + "…"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %-24s  %-16s  added: %s\n", e.Name, preview, e.Added)
			}
			return nil
		},
	}
}

func newTrustAddCmd() *cobra.Command {
	var trustFile string
	var name string
	var keyFile string
	var keyEnv string
	var keyValue string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Register a trusted signing key",
		Long: `Adds a key to the trust store (.skpm-trust.yaml).

Key sources (in priority order):
  --key-file  path to a file containing a base64-encoded key
  --key-env   name of an environment variable (base64-encoded)
  --key-value base64-encoded key value directly (avoid in CI — use --key-env)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if trustFile == "" {
				trustFile = defaultTrustFile
			}
			if name == "" {
				return &UserError{Message: "--name is required"}
			}

			var raw string
			switch {
			case keyFile != "":
				data, err := os.ReadFile(keyFile)
				if err != nil {
					return &UserError{Message: fmt.Sprintf("read key file: %v", err)}
				}
				raw = strings.TrimSpace(string(data))
			case keyEnv != "":
				raw = os.Getenv(keyEnv)
				if raw == "" {
					return &UserError{Message: fmt.Sprintf("environment variable %s is empty or not set", keyEnv)}
				}
			case keyValue != "":
				raw = keyValue
			default:
				return &UserError{Message: "provide --key-file, --key-env, or --key-value"}
			}

			// Validate the key is valid base64 and meets minimum length.
			decoded, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("key is not valid base64: %v", err)}
			}
			if len(decoded) < 16 {
				return &UserError{Message: "key is too short (minimum 16 bytes after decoding)"}
			}

			ts, err := LoadTrustStore(trustFile)
			if err != nil {
				return err
			}
			for _, e := range ts.Keys {
				if e.Name == name {
					return &UserError{Message: fmt.Sprintf("key %q already exists — revoke it first", name)}
				}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would add trusted key %q to %s\n", name, trustFile)
				return nil
			}

			ts.Keys = append(ts.Keys, TrustEntry{
				Name:  name,
				Key:   raw,
				Added: time.Now().UTC().Format(time.RFC3339),
			})
			if err := saveTrustStore(trustFile, ts); err != nil {
				return &InternalError{Message: "save trust store", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added trusted key %q → %s\n", name, trustFile)
			return nil
		},
	}

	cmd.Flags().StringVar(&trustFile, "trust-file", "", "Trust store path (default: .skpm-trust.yaml)")
	cmd.Flags().StringVar(&name, "name", "", "Name for this trusted key (required)")
	cmd.Flags().StringVar(&keyFile, "key-file", "", "Path to file with base64-encoded key")
	cmd.Flags().StringVar(&keyEnv, "key-env", "", "Environment variable with base64-encoded key")
	cmd.Flags().StringVar(&keyValue, "key-value", "", "Base64-encoded key value")
	return cmd
}

func newTrustRevokeCmd() *cobra.Command {
	var trustFile string
	var name string

	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Remove a trusted signing key from the trust store",
		RunE: func(cmd *cobra.Command, args []string) error {
			if trustFile == "" {
				trustFile = defaultTrustFile
			}
			if name == "" && len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				return &UserError{Message: "--name is required"}
			}

			ts, err := LoadTrustStore(trustFile)
			if err != nil {
				return err
			}

			before := len(ts.Keys)
			filtered := ts.Keys[:0]
			for _, e := range ts.Keys {
				if e.Name != name {
					filtered = append(filtered, e)
				}
			}
			if len(filtered) == before {
				return &UserError{Message: fmt.Sprintf("trusted key %q not found", name)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would remove trusted key %q from %s\n", name, trustFile)
				return nil
			}

			ts.Keys = filtered
			sort.Slice(ts.Keys, func(i, j int) bool { return ts.Keys[i].Name < ts.Keys[j].Name })
			if err := saveTrustStore(trustFile, ts); err != nil {
				return &InternalError{Message: "save trust store", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Revoked trusted key %q from %s\n", name, trustFile)
			return nil
		},
	}

	cmd.Flags().StringVar(&trustFile, "trust-file", "", "Trust store path (default: .skpm-trust.yaml)")
	cmd.Flags().StringVar(&name, "name", "", "Name of the key to revoke")
	return cmd
}
