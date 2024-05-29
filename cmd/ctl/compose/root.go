package compose

import (
	// Common
	// External

	"github.com/spf13/cobra"
)

// Constants
const DefaultSelksPath = ".configs/selks/embedded"

// Commands
func ComposeCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "compose",
		Short: "Create container compose file",
	}

	// Custom commands
	cmd.AddCommand(initCmd())
	cmd.AddCommand(configCmd())
	cmd.AddCommand(wrappedCmd()...)

	return cmd
}
