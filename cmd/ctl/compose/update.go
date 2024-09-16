package compose

import (
	// Common

	// External

	"github.com/spf13/cobra"

	// Custom
	flags "stamus-ctl/internal/handlers"
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
	flags.Version.AddAsFlag(cmd, false)
	flags.Config.AddAsFlag(cmd, false)
	return cmd
}

func updateHandler(cmd *cobra.Command, args []string) error {
	// Validate flags
	version, err := flags.Version.GetValue()
	if err != nil {
		return err
	}
	config, err := flags.Config.GetValue()
	if err != nil {
		return err
	}
	// Call handler
	params := handlers.UpdateHandlerParams{
		Version: version.(string),
		Config:  config.(string),
		Args:    args,
	}
	return handlers.UpdateHandler(params)

}
