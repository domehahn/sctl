package cli

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Version, Commit, Date are set at build time via -ldflags.
// When installed via `go install`, they fall back to the embedded module version.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func init() {
	if Version != "dev" {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		Version = info.Main.Version
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= 7 {
				Commit = s.Value[:7]
			}
		case "vcs.time":
			Date = s.Value
		}
	}
}

var (
	globalOutput      string
	globalDryRun      bool
	globalVerbose     bool
	globalConcurrency int
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "skpm",
		Short: "Skill Package Manager for AI agent skills",
		Long: `skpm manages AI agent skills as versioned software artifacts.

Skills are SKILL.md-based capability bundles for Claude Code, GitLab Duo,
GitHub Copilot, and Codex. skpm installs, validates, and packages them.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initLogger(globalVerbose, globalOutput)
			return nil
		},
	}

	root.PersistentFlags().StringVarP(&globalOutput, "output", "o", "text", "Output format: text or json")
	root.PersistentFlags().BoolVar(&globalDryRun, "dry-run", false, "Print what would be done without making changes")
	root.PersistentFlags().BoolVarP(&globalVerbose, "verbose", "v", false, "Enable verbose/debug logging")
	root.PersistentFlags().IntVar(&globalConcurrency, "concurrency", 4, "Maximum parallel downloads")

	root.AddCommand(newInitCmd())
	root.AddCommand(newConfigCmd())
	root.AddCommand(newAddCmd())
	root.AddCommand(newRemoveCmd())
	root.AddCommand(newLockCmd())
	root.AddCommand(newInstallCmd())
	root.AddCommand(newUpdateCmd())
	root.AddCommand(newOutdatedCmd())
	root.AddCommand(newListCmd())
	root.AddCommand(newSearchCmd())
	root.AddCommand(newInfoCmd())
	root.AddCommand(newRegistryCmd())
	root.AddCommand(newValidateCmd())
	root.AddCommand(newLintCmd())
	root.AddCommand(newFormatCmd())
	root.AddCommand(newPackageCmd())
	root.AddCommand(newPublishCmd())
	root.AddCommand(newVerifyCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newCacheCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func newVersionCmd() *cobra.Command {
	var debugBuildInfo bool
	c := &cobra.Command{
		Use:   "version",
		Short: "Print skpm version information or manage skill versions",
		Run: func(cmd *cobra.Command, args []string) {
			if debugBuildInfo {
				info, ok := debug.ReadBuildInfo()
				if !ok {
					fmt.Fprintln(cmd.OutOrStdout(), "debug.ReadBuildInfo() returned false")
					return
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Main.Path:    %q\n", info.Main.Path)
				fmt.Fprintf(cmd.OutOrStdout(), "Main.Version: %q\n", info.Main.Version)
				for _, s := range info.Settings {
					fmt.Fprintf(cmd.OutOrStdout(), "Setting:      %s=%q\n", s.Key, s.Value)
				}
				return
			}
			format := outputFormat()
			if format == OutputJSON {
				PrintResult(format, CommandResult{
					Success: true,
					Command: "version",
					Data: map[string]string{
						"version": Version,
						"commit":  Commit,
						"date":    Date,
					},
				})
				return
			}
			cmd.Printf("skpm %s (commit %s, built %s)\n", Version, Commit, Date)
		},
	}
	c.Flags().BoolVar(&debugBuildInfo, "debug-build-info", false, "")
	_ = c.Flags().MarkHidden("debug-build-info")
	c.AddCommand(newVersionShowCmd())
	c.AddCommand(newVersionBumpCmd())
	c.AddCommand(newVersionSetCmd())
	return c
}

func outputFormat() OutputFormat {
	if strings.ToLower(globalOutput) == "json" {
		return OutputJSON
	}
	return OutputText
}

func initLogger(verbose bool, output string) {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	if strings.ToLower(output) == "json" {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
	}
}
