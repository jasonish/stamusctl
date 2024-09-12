package config

import (
	// Custom
	"fmt"
	flags "stamus-ctl/internal/handlers"
	config "stamus-ctl/internal/handlers/config"

	// External
	"github.com/spf13/cobra"
)

// Command
func setCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "set [keys=values...]",
		Short: "Set config file parameters",
		Long: `Set config file parameters
Input keys and values of parameters to set.
Example: set scirius.token=AwesomeToken`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return setHandler(cmd, args)
		},
	}
	// Subcommands
	cmd.AddCommand(setContentCmd())
	// Flags
	flags.Config.AddAsFlag(cmd, false)
	flags.Values.AddAsFlag(cmd, false)
	flags.Reload.AddAsFlag(cmd, false)
	flags.Apply.AddAsFlag(cmd, false)
	flags.FromFile.AddAsFlag(cmd, false)
	return cmd
}

// Subcommands
func setContentCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Place a file or folder content in the configuration",
		Long: `Place a file or folder content in the configuration
Example: config content /nginx:/etc/nginx /nginx.conf:/etc/nginx/nginx.conf,
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return setContentHandler(cmd, args)
		},
	}
	// Flags
	flags.Config.AddAsFlag(cmd, false)
	return cmd
}

// Handlers
func setHandler(cmd *cobra.Command, args []string) error {
	// Get properties
	configPath, err := flags.Config.GetValue()
	if err != nil {
		return err
	}
	isValidConfig := flags.Config.IsValid()
	if !isValidConfig {
		return fmt.Errorf("Invalid output path")
	}
	reload, err := flags.Reload.GetValue()
	if err != nil {
		return err
	}
	apply, err := flags.Apply.GetValue()
	if err != nil {
		return err
	}
	values, err := flags.Values.GetValue()
	if err != nil {
		return err
	}
	fromFile, err := flags.FromFile.GetValue()
	if err != nil {
		return err
	}

	// Set the values
	params := config.SetHandlerInputs{
		Config:   configPath.(string),
		Args:     args,
		Reload:   reload.(bool),
		Apply:    apply.(bool),
		Values:   values.(string),
		FromFile: fromFile.(string),
	}
	err = config.SetHandler(params)
	if err != nil {
		return err
	}
	return nil
}

func setContentHandler(cmd *cobra.Command, args []string) error {
	// Get flags
	configPath, err := flags.Config.GetValue()
	if err != nil {
		return err
	}

	// Call handler
	return config.SetContentHandler(configPath.(string), args)
}
