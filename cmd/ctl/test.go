package ctl

import (
	// External
	"github.com/spf13/cobra"
	// Custom

	flags "stamus-ctl/internal/handlers"
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/internal/utils"
)

func testCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:    "test",
		Short:  "Create a test configuration",
		Long:   "Internal utility command to create a test configuration",
		Hidden: true,
		RunE:   testHandler,
	}
	// Flags
	flags.IsDefaultParam.AddAsFlag(cmd, false)
	flags.OutputPath.AddAsFlag(cmd, false)
	flags.Values.AddAsFlag(cmd, false)
	flags.FromFile.AddAsFlag(cmd, false)
	return cmd
}

func testHandler(cmd *cobra.Command, args []string) error { // Get flags
	isDefault, err := flags.IsDefaultParam.GetValue()
	if err != nil {
		return err
	}
	outputPath, err := flags.OutputPath.GetValue()
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

	// Call handler
	selksInitParams := handlers.InitHandlerInputs{
		IsDefault:        isDefault.(bool),
		BackupFolderPath: ".test/config",
		OutputPath:       outputPath.(string),
		Arbitrary:        utils.ExtractArgs(args),
		Project:          "",
		Version:          "",
		Values:           values.(string),
		FromFile:         fromFile.(string),
	}
	return handlers.InitHandler(true, selksInitParams)
}
