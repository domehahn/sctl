package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/spf13/cobra"
)

func newPatchCmd() *cobra.Command {
	var lockPath string
	var patchFile string
	var stripN int
	var reverse bool

	cmd := &cobra.Command{
		Use:   "patch <skill-name>",
		Short: "Apply a local patch file to an installed skill without changing the lockfile",
		Long: `Applies a unified diff (.patch file) to the files of an installed skill.
Useful when upstream has a bug and a fix isn't released yet.

The patch is applied using the system 'patch' utility (must be in PATH).
The lockfile version is NOT changed — this is an emergency workaround, not
a version management operation. Document the patch in your repo's README.

Examples:
  skpm patch security-reviewer --file fix-regex.patch
  skpm patch security-reviewer --file fix-regex.patch --reverse   # undo`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			skillName := args[0]

			if patchFile == "" {
				return &UserError{Message: "--file is required"}
			}
			absPath, err := filepath.Abs(patchFile)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("resolve patch file: %v", err)}
			}
			if _, err := os.Stat(absPath); err != nil {
				return &UserError{Message: fmt.Sprintf("patch file not found: %s", absPath)}
			}

			if lockPath == "" {
				lockPath = lockfile.DefaultFilename
			}
			lf, err := lockfile.Read(lockPath)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read %s: %v", lockPath, err)}
			}

			var target *lockfile.SkillLock
			for i := range lf.Skills {
				if lf.Skills[i].Name == skillName {
					target = &lf.Skills[i]
					break
				}
			}
			if target == nil {
				return &UserError{Message: fmt.Sprintf("skill %q not found in lockfile", skillName)}
			}

			var installDir string
			for _, p := range target.InstalledTo {
				if info, err := os.Stat(p); err == nil && info.IsDir() {
					installDir = p
					break
				}
			}
			if installDir == "" {
				return &UserError{Message: fmt.Sprintf("no valid install directory found for %q — run: skpm install", skillName)}
			}

			patchBin, err := exec.LookPath("patch")
			if err != nil {
				return &UserError{Message: "'patch' utility not found in PATH — install it (e.g. brew install patch)"}
			}

			patchArgs := []string{
				fmt.Sprintf("-p%d", stripN),
				"--input", absPath,
				"--directory", installDir,
			}
			if reverse {
				patchArgs = append(patchArgs, "--reverse")
			}
			if globalDryRun {
				patchArgs = append(patchArgs, "--dry-run")
			}

			c := exec.Command(patchBin, patchArgs...)
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()

			if err := c.Run(); err != nil {
				if globalDryRun {
					return &UserError{Message: fmt.Sprintf("patch dry-run reported problems: %v", err)}
				}
				return &UserError{Message: fmt.Sprintf("patch failed: %v", err)}
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDry run: patch applies cleanly to %s@%s\n", skillName, target.Version)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "\nPatched %s@%s in %s\n", skillName, target.Version, installDir)
				if reverse {
					fmt.Fprintf(cmd.OutOrStdout(), "Patch reversed.\n")
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lock", "", "Lockfile path (default: agent-skills.lock)")
	cmd.Flags().StringVarP(&patchFile, "file", "f", "", "Path to the .patch file (required)")
	cmd.Flags().IntVarP(&stripN, "strip", "p", 1, "Strip N leading path components (passed to patch -p)")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "Reverse the patch (undo a previously applied patch)")
	return cmd
}
