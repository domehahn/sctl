package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

// knownHooks lists all lifecycle hook names in execution order.
var knownHooks = []string{
	"pre_add", "post_add",
	"pre_install", "post_install",
	"pre_update", "post_update",
	"pre_publish", "post_publish",
	"pre_release", "post_release",
}

func newHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage and run lifecycle hooks defined in agent-skills.yaml",
		Long: `Lifecycle hooks are shell commands defined in agent-skills.yaml under 'hooks'.
skpm runs them automatically before and after key operations.

  hooks:
    pre_install:  echo "installing skills..."
    post_install: ./scripts/verify-install.sh
    post_add:     git add agent-skills.yaml agent-skills.lock
    pre_publish:  skpm validate .

Available hook names:
  pre_add / post_add
  pre_install / post_install
  pre_update / post_update
  pre_publish / post_publish
  pre_release / post_release`,
	}
	cmd.AddCommand(newHooksListCmd())
	cmd.AddCommand(newHooksRunCmd())
	return cmd
}

func newHooksListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show hooks configured in agent-skills.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "No agent-skills.yaml found.")
					return nil
				}
				return &InternalError{Message: "read manifest", Cause: err}
			}
			if len(mf.Hooks) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No hooks configured. Add a 'hooks' section to agent-skills.yaml.")
				return nil
			}
			// Print in canonical order, then any extra keys alphabetically.
			printed := map[string]bool{}
			for _, name := range knownHooks {
				if script, ok := mf.Hooks[name]; ok {
					fmt.Fprintf(cmd.OutOrStdout(), "  %-24s %s\n", name, script)
					printed[name] = true
				}
			}
			extras := []string{}
			for k := range mf.Hooks {
				if !printed[k] {
					extras = append(extras, k)
				}
			}
			sort.Strings(extras)
			for _, name := range extras {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-24s %s\n", name, mf.Hooks[name])
			}
			return nil
		},
	}
}

func newHooksRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "run <hook-name>",
		Short:     "Manually trigger a lifecycle hook",
		ValidArgs: knownHooks,
		Args:      cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]
			mf, err := manifest.Read(manifest.DefaultFilename)
			if err != nil {
				if os.IsNotExist(err) {
					return &UserError{Message: "no agent-skills.yaml found"}
				}
				return &InternalError{Message: "read manifest", Cause: err}
			}
			script, ok := mf.Hooks[hookName]
			if !ok {
				return &UserError{Message: fmt.Sprintf("hook %q is not defined in agent-skills.yaml", hookName)}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Running hook %q\n> %s\n", hookName, script)
			if globalDryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "(dry run — not executed)")
				return nil
			}
			return execHook(hookName, script, cmd)
		},
	}
}

// runProjectHook reads agent-skills.yaml and executes the named hook if defined.
// Errors from hooks are printed as warnings but do not abort the operation unless
// the hook name starts with "pre_" (pre-hooks failing are fatal).
func runProjectHook(cmd *cobra.Command, hookName string) error {
	mf, err := manifest.Read(manifest.DefaultFilename)
	if err != nil {
		return nil // no manifest — silently skip
	}
	script, ok := mf.Hooks[hookName]
	if !ok || script == "" {
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "  [hook] %s\n", hookName)
	if globalDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "    (dry run) > %s\n", script)
		return nil
	}

	if err := execHook(hookName, script, cmd); err != nil {
		if strings.HasPrefix(hookName, "pre_") {
			return &UserError{Message: fmt.Sprintf("pre-hook %q failed: %v", hookName, err)}
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "  warning: post-hook %q failed: %v\n", hookName, err)
	}
	return nil
}

func execHook(name, script string, cmd *cobra.Command) error {
	sh, flag := shellInterpreter()
	c := exec.Command(sh, flag, script) //nolint:gosec
	c.Stdout = cmd.OutOrStdout()
	c.Stderr = cmd.ErrOrStderr()
	c.Env = append(os.Environ(), "SKPM_HOOK="+name)
	if err := c.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("exited with code %d", exit.ExitCode())
		}
		return err
	}
	return nil
}
