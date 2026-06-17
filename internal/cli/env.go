package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/domehahn/skpm/v2/internal/manifest"
	"github.com/spf13/cobra"
)

func newEnvCmd() *cobra.Command {
	var exportMode bool

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Show the resolved skpm runtime environment",
		Long: `Prints the active configuration: config file path, cache directory,
default registry, all configured registries with auth status, and
the locations of the manifest and lockfile in the current directory.

Useful for diagnosing auth problems, wrong cache paths, or misconfigured
registries without having to inspect files manually.

Use --export to print tokens as shell export statements for eval in CI:
  eval $(skpm env --export)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// --export: print registry tokens as shell export statements.
			if exportMode {
				return envExport(cmd)
			}

			format := outputFormat()

			cfgPath, _ := config.Path()
			cfg, cfgErr := config.Load()

			type registryInfo struct {
				Name      string `json:"name"`
				Type      string `json:"type"`
				URL       string `json:"url,omitempty"`
				Default   bool   `json:"default"`
				AuthType  string `json:"auth_type"`
				HasToken  bool   `json:"has_token"`
				TokenEnv  string `json:"token_env,omitempty"`
				TokenFrom string `json:"token_from,omitempty"`
			}

			var registries []registryInfo
			var defaultRegistry string
			cacheDir := ""

			if cfg != nil {
				defaultRegistry = cfg.DefaultRegistry
				cacheDir = cfg.CacheDir

				names := make([]string, 0, len(cfg.Registries))
				for n := range cfg.Registries {
					names = append(names, n)
				}
				sort.Strings(names)

				for _, n := range names {
					rc := cfg.Registries[n]
					info := registryInfo{
						Name:    n,
						Type:    rc.Type,
						URL:     rc.URL,
						Default: n == cfg.DefaultRegistry,
					}
					info.AuthType = rc.Auth.Type
					if info.AuthType == "" {
						info.AuthType = "none"
					}

					if rc.Auth.TokenEnv != "" {
						info.TokenEnv = rc.Auth.TokenEnv
						if v := os.Getenv(rc.Auth.TokenEnv); v != "" {
							info.HasToken = true
							info.TokenFrom = "env:" + rc.Auth.TokenEnv
						} else {
							info.TokenFrom = "env:" + rc.Auth.TokenEnv + " (unset)"
						}
					} else if rc.Auth.Token != "" || rc.Token != "" {
						info.HasToken = true
						info.TokenFrom = "config"
					}

					registries = append(registries, info)
				}
			}

			// Manifest and lockfile status.
			mfStatus := fileStatus{Path: manifest.DefaultFilename}
			if mf, err := manifest.Read(manifest.DefaultFilename); err == nil {
				mfStatus.Exists = true
				mfStatus.Count = len(mf.Skills)
			}

			lfStatus := fileStatus{Path: lockfile.DefaultFilename}
			if lf, err := lockfile.Read(lockfile.DefaultFilename); err == nil {
				lfStatus.Exists = true
				lfStatus.Count = len(lf.Skills)
			}

			if format == OutputJSON {
				type jsonOut struct {
					ConfigPath      string         `json:"config_path"`
					ConfigExists    bool           `json:"config_exists"`
					ConfigError     string         `json:"config_error,omitempty"`
					CacheDir        string         `json:"cache_dir"`
					DefaultRegistry string         `json:"default_registry,omitempty"`
					Registries      []registryInfo `json:"registries"`
					Manifest        fileStatus     `json:"manifest"`
					Lockfile        fileStatus     `json:"lockfile"`
				}
				out := jsonOut{
					ConfigPath:      cfgPath,
					ConfigExists:    cfgErr == nil,
					CacheDir:        cacheDir,
					DefaultRegistry: defaultRegistry,
					Registries:      registries,
					Manifest:        mfStatus,
					Lockfile:        lfStatus,
				}
				if cfgErr != nil {
					out.ConfigError = cfgErr.Error()
				}
				PrintResult(format, CommandResult{Success: true, Command: "env", Data: out})
				return nil
			}

			// ── Text output ───────────────────────────────────────────
			fmt.Fprintf(cmd.OutOrStdout(), "skpm environment\n\n")

			cfgExistsLabel := "exists"
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				cfgExistsLabel = "not found (defaults apply)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  config:         %s  [%s]\n", cfgPath, cfgExistsLabel)
			if cfgErr != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  config error:   %v\n", cfgErr)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  cache dir:      %s\n", cacheDir)

			if defaultRegistry != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  default reg:    %s\n", defaultRegistry)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "  default reg:    (none configured)\n")
			}

			fmt.Fprintln(cmd.OutOrStdout())

			if len(registries) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  registries:     (none configured)\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "  registries:\n")
				for _, r := range registries {
					marker := "  "
					if r.Default {
						marker = "* "
					}
					authLabel := r.AuthType
					if r.HasToken {
						authLabel += "  token: " + r.TokenFrom
					} else if r.AuthType != "none" {
						authLabel += "  token: missing"
					}
					urlPart := ""
					if r.URL != "" {
						urlPart = "  " + r.URL
					}
					fmt.Fprintf(cmd.OutOrStdout(), "    %s%-28s %-12s  %s%s\n", marker, r.Name, r.Type, authLabel, urlPart)
				}
			}

			fmt.Fprintln(cmd.OutOrStdout())
			printFileStatus(cmd, "manifest", mfStatus)
			printFileStatus(cmd, "lockfile", lfStatus)

			// Links file.
			if _, err := os.Stat(linksFile); err == nil {
				lm, _ := readLinksManifest()
				count := 0
				if lm != nil {
					count = len(lm.Links)
					names := make([]string, 0, count)
					for n := range lm.Links {
						names = append(names, n)
					}
					sort.Strings(names)
					fmt.Fprintf(cmd.OutOrStdout(), "  links:          %s  (%d linked: %s)\n",
						linksFile, count, strings.Join(names, ", "))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&exportMode, "export", false, "Print registry tokens as shell export statements (for eval)")
	return cmd
}

// envExport prints each registry token as a shell export statement suitable
// for eval $(skpm env --export) in CI pipelines.
func envExport(cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return &InternalError{Message: "load config", Cause: err}
	}

	names := make([]string, 0, len(cfg.Registries))
	for n := range cfg.Registries {
		names = append(names, n)
	}
	sort.Strings(names)

	count := 0
	for _, n := range names {
		rc := cfg.Registries[n]
		tok := rc.Token
		if tok == "" {
			tok = rc.Auth.Token
		}
		if tok == "" && rc.Auth.TokenEnv != "" {
			tok = os.Getenv(rc.Auth.TokenEnv)
		}
		if tok == "" {
			continue
		}
		// Variable name: SKPM_TOKEN_<REGISTRY_NAME_UPPER> with non-alnum replaced by _.
		varName := "SKPM_TOKEN_" + envVarSafe(n)
		fmt.Fprintf(cmd.OutOrStdout(), "export %s=%q\n", varName, tok)
		count++
	}
	if count == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "# no registry tokens configured")
	}
	return nil
}

// envVarSafe converts a registry name to an uppercase shell-safe identifier.
func envVarSafe(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(s) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

type fileStatus struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
	Count  int    `json:"count,omitempty"`
}

func printFileStatus(cmd *cobra.Command, label string, fs fileStatus) {
	if fs.Exists {
		fmt.Fprintf(cmd.OutOrStdout(), "  %-16s%s  (%d skill(s))\n", label+":", fs.Path, fs.Count)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "  %-16s%s  (not found)\n", label+":", fs.Path)
	}
}
