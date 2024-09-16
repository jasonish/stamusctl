package compose

import (
	// External
	"github.com/spf13/cobra"

	// Internal

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
	// Command
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init compose config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return SELKSHandler(cmd, args)
		},
	}
	// Flags
	flags.IsDefaultParam.AddAsFlag(cmd, false)
	flags.Values.AddAsFlag(cmd, false)
	flags.FromFile.AddAsFlag(cmd, false)
	flags.Config.AddAsFlag(cmd, false)
	// Commands
	cmd.AddCommand(SELKSCmd())
	return cmd
}

func SELKSCmd() *cobra.Command {
	// Setup
	initSelksFolder(app.DefaultSelksPath)
	// Command
	cmd := &cobra.Command{
		Use:   "selks",
		Short: "Init SELKS container compose file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return SELKSHandler(cmd, args)
		},
	}
	// Flags
	flags.IsDefaultParam.AddAsFlag(cmd, false)
	flags.Values.AddAsFlag(cmd, false)
	flags.FromFile.AddAsFlag(cmd, false)
	flags.Config.AddAsFlag(cmd, false)
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
	// Get flags
	isDefault, err := flags.IsDefaultParam.GetValue()
	if err != nil {
		return err
	}
	values, err := flags.Values.GetValue()
	if err != nil {
		return err
	}
	fromFile, err := flags.FromFile.GetValue()
	if err != nil {
		return err
	}
	config, err := flags.Config.GetValue()
	if err != nil {
		return err
	}

	// Call handler
	selksInitParams := handlers.InitHandlerInputs{
		IsDefault:        isDefault.(bool),
		BackupFolderPath: app.DefaultSelksPath,
		Arbitrary:        utils.ExtractArgs(args),
		Project:          "selks",
		Version:          "latest",
		Values:           values.(string),
		Config:           config.(string),
		FromFile:         fromFile.(string),
	}
	return handlers.InitHandler(true, selksInitParams)
}
