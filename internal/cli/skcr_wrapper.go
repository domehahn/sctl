package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// runAuthoringCompatibility executes an authoring command (such as create, template,
// format, changelog, watch, version bump/set). If the 'skcr' binary is present in PATH,
// it delegates execution directly to 'skcr'. Otherwise, it displays a clear compatibility
// notice on stderr and runs fallback logic to prevent breaking existing workflows.
func runAuthoringCompatibility(cmd *cobra.Command, skcrArgs []string, fallback func() error) error {
	skcrPath, err := exec.LookPath("skcr")
	if err == nil && skcrPath != "" {
		fmt.Fprintf(os.Stderr, "[notice] Delegating authoring command to skcr: skcr %v\n", skcrArgs)
		c := exec.CommandContext(cmd.Context(), skcrPath, skcrArgs...)
		c.Stdin = os.Stdin
		c.Stdout = cmd.OutOrStdout()
		c.Stderr = cmd.ErrOrStderr()
		if err := c.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					os.Exit(ws.ExitStatus())
				}
			}
			return &UserError{Message: fmt.Sprintf("skcr delegation failed: %v", err)}
		}
		return nil
	}

	fmt.Fprintf(os.Stderr, "[notice] '%s' is an authoring command preserved for compatibility. For dedicated skill authoring, use 'skcr'.\n", cmd.CommandPath())
	return fallback()
}
