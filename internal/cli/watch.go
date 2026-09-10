package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	var (
		dir      string
		debounce time.Duration
		exts     []string
	)

	cmd := &cobra.Command{
		Use:   "watch <script>",
		Short: "Watch skill files for changes and re-run a script automatically",
		Long: `Watches the skill directory for file changes and re-runs the named script
on every save. Useful for a fast validate/lint/test loop during authoring.

  skpm watch test
  skpm watch lint --dir ./my-skill
  skpm watch validate --ext .md,.yaml

The script must be defined in skill.yaml under 'scripts'. Use --debounce to
control how long to wait after the last change before re-running (default: 300ms).

Press Ctrl+C to stop.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthoringCompatibility(cmd, append([]string{"watch"}, args...), func() error {
				scriptName := args[0]

				skillDir := dir
				if skillDir == "" {
					skillDir = "."
				}
				absDir, err := filepath.Abs(skillDir)
				if err != nil {
					return &InternalError{Message: "resolve path", Cause: err}
				}

				// Verify the script exists before starting the watcher.
				scripts, err := readSkillScripts(absDir)
				if err != nil {
					return err
				}
				if _, ok := scripts[scriptName]; !ok {
					names := scriptNames(scripts)
					return &UserError{Message: fmt.Sprintf("unknown script %q; available: %s", scriptName, strings.Join(names, ", "))}
				}

				watchExts := map[string]bool{}
				for _, e := range exts {
					e = strings.TrimPrefix(e, ".")
					watchExts["."+e] = true
				}
				if len(watchExts) == 0 {
					for _, e := range []string{".md", ".yaml", ".yml", ".txt", ".json"} {
						watchExts[e] = true
					}
				}

				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					return &InternalError{Message: "create watcher", Cause: err}
				}
				defer watcher.Close()

				if err := watcher.Add(absDir); err != nil {
					return &InternalError{Message: "watch directory", Cause: err}
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Watching %s for changes (script: %q, debounce: %s)\n", absDir, scriptName, debounce)
				fmt.Fprintf(cmd.OutOrStdout(), "Press Ctrl+C to stop.\n\n")

				// Run once immediately on start.
				runWatchScript(cmd, absDir, scriptName)

				var timer *time.Timer
				ctx := cmd.Context()

				for {
					select {
					case <-ctx.Done():
						return nil

					case event, ok := <-watcher.Events:
						if !ok {
							return nil
						}
						if !watchExts[filepath.Ext(event.Name)] {
							continue
						}
						if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
							if timer != nil {
								timer.Stop()
							}
							timer = time.AfterFunc(debounce, func() {
								fmt.Fprintf(cmd.OutOrStdout(), "\n[%s] change detected: %s\n", time.Now().Format("15:04:05"), filepath.Base(event.Name))
								runWatchScript(cmd, absDir, scriptName)
							})
						}

					case watchErr, ok := <-watcher.Errors:
						if !ok {
							return nil
						}
						fmt.Fprintf(cmd.ErrOrStderr(), "watcher error: %v\n", watchErr)
					}
				}
			})
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Skill directory to watch (default: .)")
	cmd.Flags().DurationVar(&debounce, "debounce", 300*time.Millisecond, "Delay after last change before re-running")
	cmd.Flags().StringSliceVar(&exts, "ext", nil, "File extensions to watch (default: .md,.yaml,.yml,.txt,.json)")

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
		return scriptNames(scripts), cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func runWatchScript(cmd *cobra.Command, dir, scriptName string) {
	scripts, err := readSkillScripts(dir)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error reading skill.yaml: %v\n", err)
		return
	}
	script, ok := scripts[scriptName]
	if !ok {
		fmt.Fprintf(cmd.ErrOrStderr(), "script %q no longer defined in skill.yaml\n", scriptName)
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "> %s\n", script)
	sh, flag := shellInterpreter()
	c := exec.Command(sh, flag, script) //nolint:gosec
	c.Dir = dir
	c.Stdout = cmd.OutOrStdout()
	c.Stderr = cmd.ErrOrStderr()
	c.Env = append(os.Environ(),
		"SKPM_SKILL_DIR="+dir,
		"SKPM_SCRIPT="+scriptName,
	)
	if err := c.Run(); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "script exited with error: %v\n", err)
	}
}

func scriptNames(scripts map[string]string) []string {
	names := make([]string, 0, len(scripts))
	for k := range scripts {
		names = append(names, k)
	}
	return names
}
