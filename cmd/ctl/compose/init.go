package compose

import (
	// Common
	"log"
	// External

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/embeds"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
)

// Flags
var output = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
	Usage:     "Declare the folder where to save configuration files",
}
var interactive = models.Parameter{
	Name:      "interactive",
	Shorthand: "i",
	Type:      "bool",
	Default:   models.CreateVariableBool(true),
	Usage:     "Interactive mode",
}

// Commands
func initCmd() *cobra.Command {
	// Setup
	initSelksFolder(DefaultSelksPath)
	// Create
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init compose config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return SELKSHandler(cmd)
		},
	}
	// Flags
	output.AddAsFlag(cmd, false)
	interactive.AddAsFlag(cmd, false)
	// Commands
	cmd.AddCommand(SELKSCmd())
	return cmd
}

func SELKSCmd() *cobra.Command {
	// Create
	cmd := &cobra.Command{
		Use:   "selks",
		Short: "Init SELKS container compose file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return SELKSHandler(cmd)
		},
	}
	return cmd

}

// Create SELKS folder if it does not exist
func initSelksFolder(path string) {
	selksConfigExist, err := utils.FolderExists(path)
	if err != nil {
		panic(err)
	}
	if !selksConfigExist {
		err = embeds.Extract()
		if err != nil {
			panic(err)
		}
	}
}

func SELKSHandler(cmd *cobra.Command) error {
	// Instanciate config
	var config models.Config
	confFile, err := models.CreateFileInstance(DefaultSelksPath, "config.yaml")
	if err != nil {
		return err
	}
	configPointer, err := models.NewConfigFrom(confFile)
	if err != nil {
		return err
	}
	config = *configPointer
	// Read the folder configuration
	_, _, err = config.ExtractParams()
	if err != nil {
		return err
	}
	// Ask for the parameters
	if *interactive.Variable.Bool {
		log.Println("Interactive mode")
		err := config.GetParams().AskAll()
		if err != nil {
			return err
		}
	}
	// Save the configuration
	outputFile, err := models.CreateFileInstance(*output.Variable.String, "config.yaml")
	if err != nil {
		return err
	}
	config.SaveConfigTo(outputFile)
	return nil
}
