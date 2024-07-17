package compose

import (
	// Common

	// External

	"fmt"

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/app"
	flags "stamus-ctl/internal/handlers"
	handlers "stamus-ctl/internal/handlers/compose"
)

// Init
func init() {
	// Setup
	initSelksFolder(app.DefaultSelksPath)
}

// Commands
func configCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Interact with compose config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHandler(cmd, args)
		},
	}
	// Flags
	flags.ConfigPath.AddAsFlag(cmd, false)

	// Add Commands
	cmd.AddCommand(getCmd())
	cmd.AddCommand(setCmd())
	return cmd
}

// Subcommands
func setCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [keys=values...]",
		Short: "Set compose config file parameters",
		Long: `Set compose config file parameters
Input keys and values of parameters to set.
Example: set scirius.token=AwesomeToken`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return setHandler(cmd, args)
		},
	}
	flags.ConfigPath.AddAsFlag(cmd, false)
	flags.Values.AddAsFlag(cmd, false)
	flags.Reload.AddAsFlag(cmd, false)
	flags.Apply.AddAsFlag(cmd, false)
	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [keys...]",
		Short: "Get compose config file parameters values",
		Long: `Get compose config file parameters values
Input the keys of parameters to get
Example: get scirius`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHandler(cmd, args)
		},
	}
	flags.ConfigPath.AddAsFlag(cmd, false)
	return cmd
}

// Handlers
func setHandler(cmd *cobra.Command, args []string) error {
	// Get properties
	configPath, err := flags.ConfigPath.GetValue()
	if err != nil {
		return err
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
	// Set the values
	params := handlers.SetHandlerInputs{
		Config: configPath.(string),
		Args:   args,
		Reload: reload.(bool),
		Apply:  apply.(bool),
		Values: values.(string),
	}
	err = handlers.SetHandler(params)
	if err != nil {
		return err
	}
	return nil
}

func getHandler(cmd *cobra.Command, args []string) error {
	// Get properties
	configPath, err := flags.ConfigPath.GetValue()
	if err != nil {
		return err
	}
	reload, err := flags.Reload.GetValue()
	if err != nil {
		return err
	}
	// Load the config values
	groupedValues, err := handlers.GetGroupedConfig(configPath.(string), args, reload.(bool))
	if err != nil {
		return err
	}
	// Print the values
	printGroupedValues(groupedValues, "")
	return nil
}

// Utility function
// From the grouped values, print the values in a readable format
func printGroupedValues(group map[string]interface{}, prefix string) {
	for key, value := range group {
		switch v := value.(type) {
		case string:
			fmt.Printf("%s%s: %s\n", prefix, key, v)
		case map[string]interface{}:
			fmt.Printf("%s%s:\n", prefix, key)
			printGroupedValues(v, prefix+"  ")
		}
	}
}
