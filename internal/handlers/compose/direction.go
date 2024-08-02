package handlers

import (
	// Common

	"strconv"

	// External
	"github.com/spf13/cobra"

	// Internal
	"stamus-ctl/internal/app"
	compose "stamus-ctl/internal/docker-compose"
	"stamus-ctl/pkg/mocker"
)

func HandleUp(configPath string) error {
	if app.Mode.IsTest() {
		mocker.Mocked.Up(configPath)
		return nil
	}
	return handleUp(configPath)
}

// HandleUp handles the up command, similar to the up command in docker-compose
func handleUp(configPath string) error {
	// Get command
	command := compose.GetComposeCmd("up")
	// Set flags
	command.Flags().Lookup("folder").DefValue = configPath
	command.Flags().Lookup("folder").Value.Set(configPath)
	command.Flags().Lookup("detach").DefValue = "true"
	command.Flags().Lookup("detach").Value.Set("true")
	// Create root command
	var cmd *cobra.Command = &cobra.Command{Use: "compose"}
	cmd.SetArgs([]string{"up"})
	cmd.AddCommand(command)
	// Run command
	return cmd.Execute()
}

func HandleDown(configPath string, removeOrphans bool, volumes bool) error {
	if app.Mode.IsTest() {
		mocker.Mocked.Down(configPath)
		return nil
	}
	return handleDown(configPath, removeOrphans, volumes)
}

// HandleDown handles the down command, similar to the down command in docker-compose
func handleDown(configPath string, removeOrphans bool, volumes bool) error {
	// Get command
	command := compose.GetComposeCmd("down")
	// Set flags
	command.Flags().Lookup("folder").DefValue = configPath
	command.Flags().Lookup("folder").Value.Set(configPath)
	command.Flags().Lookup("remove-orphans").DefValue = strconv.FormatBool(removeOrphans)
	command.Flags().Lookup("remove-orphans").Value.Set(strconv.FormatBool(removeOrphans))
	command.Flags().Lookup("volumes").DefValue = strconv.FormatBool(volumes)
	command.Flags().Lookup("volumes").Value.Set(strconv.FormatBool(volumes))
	// Create root command
	var cmd *cobra.Command = &cobra.Command{Use: "compose"}
	cmd.SetArgs([]string{"up"})
	cmd.AddCommand(command)
	// Run command
	return cmd.Execute()
}
