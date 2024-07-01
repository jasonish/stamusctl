package compose

import (
	// Common

	// External

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/embeds"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
)

// Constants
const embed string = "selks"

// Flags
var output = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
	Usage:     "Declare the folder where to save configuration files",
}
var defaultSettings = models.Parameter{
	Name:      "default",
	Shorthand: "d",
	Type:      "bool",
	Default:   models.CreateVariableBool(false),
	Usage:     "Set to default settings",
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
			return SELKSHandler(cmd, args)
		},
	}
	// Flags
	output.AddAsFlag(cmd, false)
	defaultSettings.AddAsFlag(cmd, false)
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
			return SELKSHandler(cmd, args)
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
		err = embeds.ExtractEmbedTo(embed, app.TemplatesFolder+"selks/embedded/")
		if err != nil {
			panic(err)
		}
	}
}

func SELKSHandler(cmd *cobra.Command, args []string) error {
	// Instanciate config
	var config *models.Config
	confFile, err := models.CreateFileInstance(LatestSelksPath, "config.yaml")
	if err != nil {
		confFile, err = models.CreateFileInstance(DefaultSelksPath, "config.yaml")
		if err != nil {
			return err
		}
	}
	config, err = models.NewConfigFrom(confFile)
	if err != nil {
		return err
	}
	// Read the folder configuration
	_, _, err = config.ExtractParams(true)
	if err != nil {
		return err
	}
	// Set parameters
	if *defaultSettings.Variable.Bool {
		// Extract and set values from args
		extractedArgs := utils.ExtractArgs(args)
		err = config.GetParams().SetLooseValues(extractedArgs)
		config.SetArbitrary(extractedArgs)
		if err != nil {
			return err
		}
		// Set from default
		err := config.GetParams().SetToDefaults()
		if err != nil {
			return err
		}
		// Ask for missing parameters
		err = config.GetParams().AskMissing()
		if err != nil {
			return err
		}
	} else {
		//Set from user input
		err := config.GetParams().AskAll()
		if err != nil {
			return err
		}
	}
	// Validate parameters
	err = config.GetParams().ValidateAll()
	if err != nil {
		return err
	}
	// Save the configuration
	outputFile, err := models.CreateFileInstance(*output.Variable.String, "values.yaml")
	if err != nil {
		return err
	}
	config.SaveConfigTo(outputFile)
	return nil
}
