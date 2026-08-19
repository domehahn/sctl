package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultPolicyFile = "skpm-policy.yaml"

// Policy is the top-level structure for skpm-policy.yaml.
type Policy struct {
	Version int         `yaml:"version"`
	Rules   PolicyRules `yaml:"rules"`
}

// PolicyRules holds all rule sets that CheckPolicy evaluates.
type PolicyRules struct {
	AllowedRegistries []string          `yaml:"allowed_registries,omitempty"`
	DenyExperimental  bool              `yaml:"deny_experimental,omitempty"`
	BannedSkills      []string          `yaml:"banned_skills,omitempty"`
	RequireSignatures bool              `yaml:"require_signatures,omitempty"`
	MinVersions       map[string]string `yaml:"min_versions,omitempty"`
}

// PolicyViolation describes a single rule failure.
type PolicyViolation struct {
	Rule    string `yaml:"rule" json:"rule"`
	Skill   string `yaml:"skill" json:"skill"`
	Message string `yaml:"message" json:"message"`
}

// CheckPolicy evaluates all policy rules against the lockfile and returns
// every violation. An empty slice means the lockfile is compliant.
func CheckPolicy(p *Policy, lf *lockfile.LockFile) []PolicyViolation {
	var violations []PolicyViolation

	allowedSet := make(map[string]bool, len(p.Rules.AllowedRegistries))
	for _, r := range p.Rules.AllowedRegistries {
		allowedSet[strings.TrimRight(r, "/")] = true
	}

	bannedSet := make(map[string]bool, len(p.Rules.BannedSkills))
	for _, b := range p.Rules.BannedSkills {
		bannedSet[b] = true
	}

	for _, sl := range lf.Skills {
		qualName := sl.Name
		if sl.Namespace != "" && sl.Namespace != "default" {
			qualName = sl.Namespace + "/" + sl.Name
		}

		if len(allowedSet) > 0 {
			src := strings.TrimRight(sl.Source, "/")
			if !allowedSet[src] {
				violations = append(violations, PolicyViolation{
					Rule:    "allowed_registries",
					Skill:   qualName,
					Message: fmt.Sprintf("source %q is not in allowed_registries", sl.Source),
				})
			}
		}

		if bannedSet[sl.Name] || bannedSet[qualName] {
			violations = append(violations, PolicyViolation{
				Rule:    "banned_skills",
				Skill:   qualName,
				Message: fmt.Sprintf("skill %q is banned by policy", qualName),
			})
		}

		if p.Rules.RequireSignatures && sl.Signature == "" {
			violations = append(violations, PolicyViolation{
				Rule:    "require_signatures",
				Skill:   qualName,
				Message: "no signature — run: skpm sign apply",
			})
		}

		if minVer, ok := p.Rules.MinVersions[sl.Name]; ok {
			if !semverAtLeast(sl.Version, minVer) {
				violations = append(violations, PolicyViolation{
					Rule:    "min_versions",
					Skill:   qualName,
					Message: fmt.Sprintf("installed %s < required minimum %s", sl.Version, minVer),
				})
			}
		}
	}
	return violations
}

func newPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Rule-based installation gates for your skill lockfile",
		Long: `Define and enforce constraints on installed skills via skpm-policy.yaml.

Available subcommands:
  init   Scaffold a starter skpm-policy.yaml
  check  Evaluate the lockfile against the policy file
  show   Pretty-print the active policy`,
	}
	cmd.AddCommand(newPolicyInitCmd())
	cmd.AddCommand(newPolicyCheckCmd())
	cmd.AddCommand(newPolicyShowCmd())
	return cmd
}

func newPolicyInitCmd() *cobra.Command {
	var policyFile string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Scaffold a starter skpm-policy.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if policyFile == "" {
				policyFile = defaultPolicyFile
			}
			if _, err := os.Stat(policyFile); err == nil {
				return &UserError{Message: fmt.Sprintf("%s already exists — delete it first to re-scaffold", policyFile)}
			}

			starter := Policy{
				Version: 1,
				Rules: PolicyRules{
					AllowedRegistries: []string{"https://registry.agentskills.io"},
					DenyExperimental:  false,
					BannedSkills:      []string{},
					RequireSignatures: false,
					MinVersions:       map[string]string{},
				},
			}

			data, err := yaml.Marshal(starter)
			if err != nil {
				return &InternalError{Message: "marshal policy", Cause: err}
			}

			header := "# skpm-policy.yaml — skill installation policy\n" +
				"# Run 'skpm policy check' to validate the lockfile against these rules.\n\n"

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would write %s\n", policyFile)
				return nil
			}
			if err := os.WriteFile(policyFile, append([]byte(header), data...), 0o644); err != nil {
				return &InternalError{Message: "write policy file", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s — edit the rules and run 'skpm policy check'\n", policyFile)
			return nil
		},
	}

	cmd.Flags().StringVar(&policyFile, "policy", "", "Output path (default: skpm-policy.yaml)")
	return cmd
}

func newPolicyCheckCmd() *cobra.Command {
	var policyFile string
	var lockPath string
	var failOnViolation bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Validate the lockfile against skpm-policy.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if policyFile == "" {
				policyFile = defaultPolicyFile
			}
			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}

			p, err := loadPolicyFile(policyFile)
			if err != nil {
				return err
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			violations := CheckPolicy(p, lf)

			format := outputFormat()
			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: len(violations) == 0,
					Command: "policy check",
					Data:    violations,
				})
				if failOnViolation && len(violations) > 0 {
					return &UserError{Message: "policy violations found"}
				}
				return nil
			}

			if len(violations) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  policy check passed — no violations")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "  %d policy violation(s):\n\n", len(violations))
			for _, v := range violations {
				fmt.Fprintf(cmd.OutOrStdout(), "  [%s]  %-28s  %s\n", v.Rule, v.Skill, v.Message)
			}
			if failOnViolation {
				return &UserError{Message: "policy violations found"}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&policyFile, "policy", "", "Policy file (default: skpm-policy.yaml)")
	cmd.Flags().StringVar(&lockPath, "lock", "", "Lockfile path (default: agent-skills.lock)")
	cmd.Flags().BoolVar(&failOnViolation, "fail", false, "Exit with non-zero code when violations exist")
	return cmd
}

func newPolicyShowCmd() *cobra.Command {
	var policyFile string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Pretty-print the active policy file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if policyFile == "" {
				policyFile = defaultPolicyFile
			}
			p, err := loadPolicyFile(policyFile)
			if err != nil {
				return err
			}

			format := outputFormat()
			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "policy show",
					Data:    p,
				})
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Policy file: %s  (version %d)\n\n", policyFile, p.Version)
			if len(p.Rules.AllowedRegistries) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  Allowed registries:")
				for _, r := range p.Rules.AllowedRegistries {
					fmt.Fprintf(cmd.OutOrStdout(), "    - %s\n", r)
				}
			}
			if len(p.Rules.BannedSkills) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  Banned skills:")
				for _, s := range p.Rules.BannedSkills {
					fmt.Fprintf(cmd.OutOrStdout(), "    - %s\n", s)
				}
			}
			if len(p.Rules.MinVersions) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  Minimum versions:")
				for k, v := range p.Rules.MinVersions {
					fmt.Fprintf(cmd.OutOrStdout(), "    %s >= %s\n", k, v)
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  require_signatures: %v\n", p.Rules.RequireSignatures)
			fmt.Fprintf(cmd.OutOrStdout(), "  deny_experimental:  %v\n", p.Rules.DenyExperimental)
			return nil
		},
	}

	cmd.Flags().StringVar(&policyFile, "policy", "", "Policy file (default: skpm-policy.yaml)")
	return cmd
}

func loadPolicyFile(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &UserError{Message: fmt.Sprintf("read policy file %s: %v — run 'skpm policy init' to create one", path, err)}
	}
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse policy file %s: %v", path, err)}
	}
	if p.Version == 0 {
		p.Version = 1
	}
	return &p, nil
}

// semverAtLeast returns true if installed >= minimum using simple dot-split comparison.
// It is intentionally simplified — only works with N.N.N numerics.
func semverAtLeast(installed, minimum string) bool {
	iv := strings.SplitN(strings.TrimPrefix(installed, "v"), ".", 3)
	mv := strings.SplitN(strings.TrimPrefix(minimum, "v"), ".", 3)
	for len(iv) < 3 {
		iv = append(iv, "0")
	}
	for len(mv) < 3 {
		mv = append(mv, "0")
	}
	for i := 0; i < 3; i++ {
		var ip, mp int
		fmt.Sscanf(iv[i], "%d", &ip)
		fmt.Sscanf(mv[i], "%d", &mp)
		if ip < mp {
			return false
		}
		if ip > mp {
			return true
		}
	}
	return true
}
