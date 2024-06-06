package compose

import (
	// Common

	// External

	"fmt"
	"strings"

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
)

// Flags
var input = models.Parameter{
	Name:      "folder",
	Shorthand: "F",
	Usage:     "Declare the folder where the configuration files are saved",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
}
var format = models.Parameter{
	Name:    "format",
	Usage:   "Format of the output (go template)",
	Type:    "string",
	Default: models.CreateVariableString("{{.}}"),
}
var reload bool = false

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
	input.AddAsFlag(cmd, false)

	// Add Commands
	cmd.AddCommand(getCmd())
	cmd.AddCommand(setCmd())
	return cmd
}

// Subcommands
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
	input.AddAsFlag(cmd, false)
	return cmd
}
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
	input.AddAsFlag(cmd, false)
	cmd.Flags().BoolVar(&reload, "reload", false, "Reload the configuration file")
	cmd.Flags().MarkHidden("reload")
	return cmd
}

// Inits
func init() {
	// Setup
	initSelksFolder(DefaultSelksPath)
}

// Handlers
func getHandler(cmd *cobra.Command, args []string) error {
	// Load the config
	value, err := input.GetValue()
	if err != nil {
		return err
	}
	inputAsString := value.(string)
	inputAsFile, err := models.CreateFileInstance(inputAsString, "config.yaml")
	if err != nil {
		return err
	}
	config, err := models.LoadConfigFrom(inputAsFile, reload)
	if err != nil {
		return err
	}

	// Print the config
	printConfig(config, args)

	return nil
}

func printConfig(config *models.Config, args []string) {
	// Print the config
	values := config.GetParams().GetValues(args...)
	groupedValues := make(map[string]interface{})
	for _, param := range config.GetParams().GetOrdered() {
		if value, ok := values[param]; ok {
			parts := strings.Split(param, ".")
			addToGroup(parts, value, groupedValues)
		}
	}
	printGroupedValues(groupedValues, "")
}

func addToGroup(parts []string, value string, group map[string]interface{}) {
	if len(parts) == 1 {
		group[parts[0]] = value
	} else {
		if _, ok := group[parts[0]]; !ok {
			group[parts[0]] = make(map[string]interface{})
		}
		addToGroup(parts[1:], value, group[parts[0]].(map[string]interface{}))
	}
}

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

func setHandler(cmd *cobra.Command, args []string) error {
	// Load the config
	value, err := input.GetValue()
	if err != nil {
		return err
	}
	inputAsString := value.(string)
	inputAsFile, err := models.CreateFileInstance(inputAsString, "config.yaml")
	if err != nil {
		return err
	}
	config, err := models.LoadConfigFrom(inputAsFile, reload)
	if err != nil {
		return err
	}
	// Set from default
	err = config.GetParams().SetToDefaults()
	if err != nil {
		return err
	}
	// Extract and set parameters from args
	paramsArgs := utils.ExtractArgs(args)
	config.GetParams().SetLooseValues(paramsArgs)
	config.SetArbitrary(paramsArgs)
	// Validate
	err = config.GetParams().ValidateAll()
	if err != nil {
		return err
	}
	// Save the configuration
	outputAsString := *output.Variable.String
	outputAsFile, err := models.CreateFileInstance(outputAsString, "config.yaml")
	if err != nil {
		return err
	}
	config.SaveConfigTo(outputAsFile)
	return nil
}
