package compose

import (
	// Common
	// External
	"github.com/spf13/cobra"

	// Custom
	"stamus-ctl/internal/app"
)

// Constants
var DefaultSelksPath string
var LatestSelksPath string

func init() {
	DefaultSelksPath = app.TemplatesFolder + "selks/embedded/"
	LatestSelksPath = app.TemplatesFolder + "selks/latest/"
}

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
	cmd.AddCommand(updateCmd())
	cmd.AddCommand(wrappedCmd()...)

	return cmd
}
