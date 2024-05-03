package compose

import (
	// Common
	// External
	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/models"
)

// Constants
const DefaultSelksPath = ".configs/selks/embedded"

// Flags
var interactive = models.Parameter{
	Name:      "interactive",
	Shorthand: "i",
	Type:      "bool",
	Default:   models.CreateVariableBool(true),
	Usage:     "Interactive mode",
}

// Commands
func ComposeCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "compose",
		Short: "Create container compose file",
	}
	// Flags
	interactive.AddAsFlag(cmd, false)
	// Commands
	cmd.AddCommand(initCmd())
	cmd.AddCommand(configCmd())

	return cmd
}
