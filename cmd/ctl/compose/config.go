package compose

import (
	// Common
	"fmt"
	"io/ioutil"
	"os"

	// External
	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
)

// Flags
var input = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
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

// Commands
func configCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Interact with compose config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHandler(cmd)
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
		Use:   "get",
		Short: "Get compose config file parameters values",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHandler(cmd)
		},
	}
	return cmd
}
func setCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set compose config file parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setHandler(cmd, args)
		},
	}
	return cmd
}

// Inits
func init() {
	// Setup
	initSelksFolder(DefaultSelksPath)
}

// Handlers
func getHandler(cmd *cobra.Command) error {
	// Read the file content
	value, err := input.GetValue()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/config.yaml", value.(string)))
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Print the content to the terminal
	fmt.Println(string(content))

	return nil
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
	config, err := models.LoadConfigFrom(inputAsFile)
	if err != nil {
		return err
	}
	// Extract and set parameters from args
	paramsArgs := utils.ExtractArgs(args)
	config.GetParams().SetLooseValues(paramsArgs)
	err = config.GetParams().ValidateAll()
	if err != nil {
		return err
	}
	// Set from default
	err = config.GetParams().SetToDefaults()
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
