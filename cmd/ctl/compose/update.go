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
	parameters.Registry.AddAsFlag(cmd, false)
	parameters.Username.AddAsFlag(cmd, false)
	parameters.Password.AddAsFlag(cmd, false)
	parameters.Version.AddAsFlag(cmd, false)
	parameters.Config.AddAsFlag(cmd, false)
	return cmd
}

func updateHandler(cmd *cobra.Command, args []string) error {
	// Validate parameters from flags
	if *parameters.Registry.Variable.String == "" {
		err := parameters.Registry.AskUser()
		if err != nil {
			return err
		}
	}
	if *parameters.Username.Variable.String == "" {
		err := parameters.Username.AskUser()
		if err != nil {
			return err
		}
	}
	if *parameters.Password.Variable.String == "" {
		err := parameters.Password.AskUser()
		if err != nil {
			return err
		}
	}
	// Get values from flags
	var registryVal, usernameVal, passwordVal, versionVal string
	registryVal = *parameters.Registry.Variable.String
	usernameVal = *parameters.Username.Variable.String
	passwordVal = *parameters.Password.Variable.String
	versionVal = *parameters.Version.Variable.String

	// Call handler
	return handlers.UpdateHandler(
		*parameters.Config.Variable.String,
		args,
		registryVal,
		usernameVal,
		passwordVal,
		versionVal,
	)

}
