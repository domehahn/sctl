package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newRunCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "run <script> [args...]",
		Short: "Execute a script defined in skill.yaml",
		Long: `Runs a named script from the scripts section of skill.yaml.

  skpm run test
  skpm run build --dir ./my-skill
  skpm run lint -- --strict

Scripts are defined in skill.yaml under the 'scripts' key:

  scripts:
    test: skpm validate . && skpm lint .
    build: skpm package .
    lint: skpm lint .`,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptName := args[0]
			extraArgs := args[1:]

			skillDir := dir
			if skillDir == "" {
				skillDir = "."
			}

			scripts, err := readSkillScripts(skillDir)
			if err != nil {
				return err
			}

			if len(scripts) == 0 {
				return &UserError{Message: "no scripts defined in skill.yaml — add a 'scripts' section"}
			}

			script, ok := scripts[scriptName]
			if !ok {
				names := make([]string, 0, len(scripts))
				for k := range scripts {
					names = append(names, k)
				}
				sort.Strings(names)
				return &UserError{Message: fmt.Sprintf("unknown script %q; available: %s", scriptName, strings.Join(names, ", "))}
			}

			if len(extraArgs) > 0 {
				script = script + " " + strings.Join(extraArgs, " ")
			}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would run %q in %s\n  > %s\n", scriptName, skillDir, script)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "> %s\n", script)

			shell, shellFlag := shellInterpreter()
			c := exec.CommandContext(cmd.Context(), shell, shellFlag, script) //nolint:gosec
			c.Dir = skillDir
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			c.Env = append(os.Environ(),
				"SKPM_SKILL_DIR="+skillDir,
				"SKPM_SCRIPT="+scriptName,
			)

			if err := c.Run(); err != nil {
				if exit, ok := err.(*exec.ExitError); ok {
					return &UserError{Message: fmt.Sprintf("script %q exited with code %d", scriptName, exit.ExitCode())}
				}
				return &InternalError{Message: "run script", Cause: err}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Skill directory containing skill.yaml (default: .)")

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		skillDir := dir
		if skillDir == "" {
			skillDir = "."
		}
		scripts, err := readSkillScripts(skillDir)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		names := make([]string, 0, len(scripts))
		for k := range scripts {
			names = append(names, k)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func readSkillScripts(dir string) (map[string]string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "skill.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &UserError{Message: fmt.Sprintf("skill.yaml not found in %s", dir)}
		}
		return nil, &InternalError{Message: "read skill.yaml", Cause: err}
	}
	var sy struct {
		Scripts map[string]string `yaml:"scripts"`
	}
	if err := yaml.Unmarshal(data, &sy); err != nil {
		return nil, &UserError{Message: fmt.Sprintf("parse skill.yaml: %v", err)}
	}
	if sy.Scripts == nil {
		sy.Scripts = map[string]string{}
	}
	return sy.Scripts, nil
}

func shellInterpreter() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/C"
	}
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	return sh, "-c"
}
