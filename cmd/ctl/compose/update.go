package compose

import (
	// Common

	// External

	"github.com/spf13/cobra"

	// Custom
	parameters "stamus-ctl/internal/handlers"
	handlers "stamus-ctl/internal/handlers/compose"
)

// Commands
func updateCmd() *cobra.Command {
	// Create cmd
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update compose configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateHandler(cmd, args)
		},
	}
	// Add flags
	parameters.Version.AddAsFlag(cmd, false)
	parameters.Config.AddAsFlag(cmd, false)
	return cmd
}

func updateHandler(cmd *cobra.Command, args []string) error {
	// Call handler
	params := handlers.UpdateHandlerParams{
		Config:  *parameters.Config.Variable.String,
		Version: *parameters.Version.Variable.String,
		Args:    args,
	}
	return handlers.UpdateHandler(params)

}
