package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const protectFile = ".skpm-protect"

type protectStore struct {
	Protected []string `yaml:"protected"`
}

func loadProtectStore(path string) (*protectStore, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &protectStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var ps protectStore
	if err := yaml.Unmarshal(data, &ps); err != nil {
		return nil, err
	}
	return &ps, nil
}

func (ps *protectStore) save(path string) error {
	data, err := yaml.Marshal(ps)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (ps *protectStore) isProtected(name string) bool {
	for _, n := range ps.Protected {
		if n == name {
			return true
		}
	}
	return false
}

func newProtectCmd() *cobra.Command {
	var lockDir string

	cmd := &cobra.Command{
		Use:   "protect <name>",
		Short: "Mark a skill as protected from upgrade and update",
		Long: `Adds the skill to ` + protectFile + ` in the lockfile directory.
Protected skills are skipped by 'skpm upgrade' and 'skpm update'.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			storePath := filepath.Join(lockDir, protectFile)

			ps, err := loadProtectStore(storePath)
			if err != nil {
				return &InternalError{Message: "load protect store", Cause: err}
			}

			if ps.isProtected(name) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already protected.\n", name)
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would protect %s\n", name)
				return nil
			}

			ps.Protected = append(ps.Protected, name)
			if err := ps.save(storePath); err != nil {
				return &InternalError{Message: "write protect store", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Protected %s\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&lockDir, "lock-dir", ".", "Directory containing "+protectFile)
	return cmd
}

func newUnprotectCmd() *cobra.Command {
	var lockDir string

	cmd := &cobra.Command{
		Use:   "unprotect <name>",
		Short: "Remove a skill's protected status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			storePath := filepath.Join(lockDir, protectFile)

			ps, err := loadProtectStore(storePath)
			if err != nil {
				return &InternalError{Message: "load protect store", Cause: err}
			}

			if !ps.isProtected(name) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is not protected.\n", name)
				return nil
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would unprotect %s\n", name)
				return nil
			}

			remaining := ps.Protected[:0]
			for _, n := range ps.Protected {
				if n != name {
					remaining = append(remaining, n)
				}
			}
			ps.Protected = remaining
			if err := ps.save(storePath); err != nil {
				return &InternalError{Message: "write protect store", Cause: err}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Unprotected %s\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&lockDir, "lock-dir", ".", "Directory containing "+protectFile)
	return cmd
}

// IsProtected returns true if the named skill is in the protect store at path.
func IsProtected(storePath, name string) bool {
	ps, err := loadProtectStore(storePath)
	if err != nil {
		return false
	}
	return ps.isProtected(name)
}

// ProtectedList returns the list of protected skill names.
func ProtectedList(storePath string) []string {
	ps, err := loadProtectStore(storePath)
	if err != nil {
		return nil
	}
	out := make([]string, len(ps.Protected))
	copy(out, ps.Protected)
	return out
}

// matchesGlob checks whether a skill name matches a simple shell glob pattern.
// Only '*' (zero or more chars) and '?' (single char) are supported.
func matchesGlob(pattern, name string) bool {
	return globMatch(pattern, name)
}

func globMatch(pat, s string) bool {
	for len(pat) > 0 {
		switch pat[0] {
		case '*':
			pat = pat[1:]
			if len(pat) == 0 {
				return true
			}
			for i := 0; i <= len(s); i++ {
				if globMatch(pat, s[i:]) {
					return true
				}
			}
			return false
		case '?':
			if len(s) == 0 {
				return false
			}
			pat = pat[1:]
			s = s[1:]
		default:
			if len(s) == 0 || pat[0] != s[0] {
				return false
			}
			pat = pat[1:]
			s = s[1:]
		}
	}
	return len(s) == 0
}
