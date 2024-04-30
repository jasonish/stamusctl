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
	configPointer, err := models.NewConfigFrom(DefaultSelksPath)
	if err != nil {
		return err
	}
	config = *configPointer
	// Read the folder configuration
	config.ExtractParams()
	// Ask for the parameters
	if *interactive.Variable.Bool {
		log.Println("Interactive mode")
		config.GetProjectParams().AskAll()
	}
	// Validate the configuration
	err = config.Validate()
	if err != nil {
		return err
	}
	// Save the configuration
	config.SaveConfigTo(*output.Variable.String)
	return nil
}
