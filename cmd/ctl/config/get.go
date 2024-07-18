package config

import (
	"fmt"
	flags "stamus-ctl/internal/handlers"
	config "stamus-ctl/internal/handlers/config"

	"github.com/spf13/cobra"
)

// Command
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
	// Subcommands
	cmd.AddCommand(getContentCmd())
	// Flags
	flags.ConfigPath.AddAsFlag(cmd, false)
	return cmd
}

// Subcommands
func getContentCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Get folder content architecture",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getContent(cmd, args)
		},
	}
	// Flags
	flags.ConfigPath.AddAsFlag(cmd, false)
	return cmd
}

// Handlers
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
	groupedValues, err := config.GetGroupedConfig(configPath.(string), args, reload.(bool))
	if err != nil {
		return err
	}
	// Print the values
	printGroupedValues(groupedValues, "")
	return nil
}

func getContent(cmd *cobra.Command, args []string) error {
	// Get properties
	configPath, err := flags.ConfigPath.GetValue()
	if err != nil {
		return err
	}
	// Call handler
	groupedContent, err := config.GetGroupedContent(configPath.(string), args)
	if err != nil {
		return err
	}

	// Print the content
	printColoredGroupedValues(groupedContent, "")
	return nil
}

// Utilities

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

func printColoredGroupedValues(group map[string]interface{}, prefix string) {
	for key, value := range group {
		switch v := value.(type) {
		case string:
			fmt.Printf("\033[2m%s%s\033[0m\n", prefix, key)
		case map[string]interface{}:
			fmt.Printf("\033[2m%s\033[0m\033[1m%s/\033[0m\n", prefix, key)
			printColoredGroupedValues(v, prefix+"|  ")
		}
	}
}
