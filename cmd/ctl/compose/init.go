package compose

import (
	// Common

	// External

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/embeds"
	parameters "stamus-ctl/internal/handlers"
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/internal/utils"
)

// Constants
const embed string = "selks"

// Commands
func initCmd() *cobra.Command {
	// Setup
	initSelksFolder(app.DefaultSelksPath)
	// Create
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init compose config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return SELKSHandler(cmd, args)
		},
	}
	// Flags
	parameters.OutputPath.AddAsFlag(cmd, false)
	parameters.IsDefaultParam.AddAsFlag(cmd, false)
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
	selksInitParams := handlers.InitHandlerInputs{
		IsDefault:        *parameters.IsDefaultParam.Variable.Bool,
		FolderPath:       app.LatestSelksPath,
		BackupFolderPath: app.DefaultSelksPath,
		OutputPath:       *parameters.OutputPath.Variable.String,
		Arbitrary:        utils.ExtractArgs(args),
		Project:          "selks",
		Version:          "latest",
	}
	return handlers.InitHandler(true, selksInitParams)
}
