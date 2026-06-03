package cli

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var (
	globalOutput      string
	globalDryRun      bool
	globalVerbose     bool
	globalConcurrency int
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "sctl",
		Short: "Skill Control — AI agent skill package manager",
		Long: `sctl manages AI agent skills as versioned software artifacts.

Skills are SKILL.md-based capability bundles for Claude Code, GitLab Duo,
GitHub Copilot, and Codex. sctl installs, validates, and packages them.`,
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

	root.AddCommand(newInstallCmd())
	root.AddCommand(newAddCmd())
	root.AddCommand(newValidateCmd())
	root.AddCommand(newPackageCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print sctl version information",
		Run: func(cmd *cobra.Command, args []string) {
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
			cmd.Printf("sctl %s (commit %s, built %s)\n", Version, Commit, Date)
		},
	}
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
