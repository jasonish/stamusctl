package compose

import (
	// Common
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// External
	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/models"
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
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/config.yaml", input.GetValue().(string)))
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
	config, err := models.LoadConfigFrom(input.GetValue().(string))
	if err != nil {
		return err
	}
	// Extract the parameters to set
	paramsArgs := make(map[string]string)
	for _, arg := range args {
		splited := strings.Split(arg, "=")
		if len(splited) != 2 {
			fmt.Println("Error: invalid argument", arg)
		} else {
			paramsArgs[splited[0]] = splited[1]
		}
	}
	// Set the parameters
	config.GetProjectParams().SetLooseValues(paramsArgs)
	// Save the configuration
	config.SaveConfigTo(*output.Variable.String)
	return nil
}
