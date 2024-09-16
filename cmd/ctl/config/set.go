package config

import (
	// Custom

	flags "stamus-ctl/internal/handlers"
	config "stamus-ctl/internal/handlers/config"
	"stamus-ctl/internal/stamus"

	// External
	"github.com/spf13/cobra"
)

// Command
func setCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "set [keys=values...]",
		Short: "Set config related stuff",
		Long: `To set current config values, input keys and values of parameters to set.
Example: set scirius.token=AwesomeToken
Or, use subcommands to set content or current configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return setHandler(cmd, args)
		},
	}
	// Subcommands
	cmd.AddCommand(setContentCmd())
	cmd.AddCommand(setCurrentCmd())
	// Flags
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
	return cmd
}

func setCurrentCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Set the current configuration to use",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setCurrentHandler(args[0])
		},
	}
	return cmd
}

// Handlers
func setHandler(cmd *cobra.Command, args []string) error {
	// Get properties
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
	// Call handler
	return config.SetContentHandler(args)
}

func setCurrentHandler(name string) error {
	return stamus.SetCurrent(name)
}
