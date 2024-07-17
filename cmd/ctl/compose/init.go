package compose

import (
	// Common

	// External

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/embeds"
	flags "stamus-ctl/internal/handlers"
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
	flags.OutputPath.AddAsFlag(cmd, false)
	flags.IsDefaultParam.AddAsFlag(cmd, false)
	flags.Values.AddAsFlag(cmd, false)
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
		IsDefault:        *flags.IsDefaultParam.Variable.Bool,
		BackupFolderPath: app.DefaultSelksPath,
		OutputPath:       *flags.OutputPath.Variable.String,
		Arbitrary:        utils.ExtractArgs(args),
		Project:          "selks",
		Version:          "latest",
		Values:           *flags.Values.Variable.String,
	}
	return handlers.InitHandler(true, selksInitParams)
}
