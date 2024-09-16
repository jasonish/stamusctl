package handlers

import (
	// Common

	"strconv"

	// External
	"github.com/spf13/cobra"

	// Internal
	"stamus-ctl/internal/app"
	compose "stamus-ctl/internal/docker-compose"
	"stamus-ctl/internal/stamus"
	"stamus-ctl/pkg/mocker"
)

func HandleUp() error {
	conf, err := stamus.GetCurrent()
	if err != nil {
		return err
	}
	if app.Mode.IsTest() {
		return mocker.Mocked.Up(conf)
	}
	return handleUp(conf)
}

// HandleUp handles the up command, similar to the up command in docker-compose
func handleUp(configName string) error {
	// Get command
	command := compose.GetComposeCmd("up")
	// Set flags
	command.Flags().Lookup("config").DefValue = configName
	command.Flags().Lookup("config").Value.Set(configName)
	command.Flags().Lookup("detach").DefValue = "true"
	command.Flags().Lookup("detach").Value.Set("true")
	// Create root command
	var cmd *cobra.Command = &cobra.Command{Use: "compose"}
	cmd.SetArgs([]string{"up"})
	cmd.AddCommand(command)
	// Run command
	return cmd.Execute()
}

func HandleDown(removeOrphans bool, volumes bool) error {
	conf, err := stamus.GetCurrent()
	if err != nil {
		return err
	}
	if app.Mode.IsTest() {
		return mocker.Mocked.Down(conf)
	}
	return handleDown(conf, removeOrphans, volumes)
}

// HandleDown handles the down command, similar to the down command in docker-compose
func handleDown(configName string, removeOrphans bool, volumes bool) error {
	// Get command
	command := compose.GetComposeCmd("down")
	// Set flags
	command.Flags().Lookup("config").DefValue = configName
	command.Flags().Lookup("config").Value.Set(configName)
	command.Flags().Lookup("remove-orphans").DefValue = strconv.FormatBool(removeOrphans)
	command.Flags().Lookup("remove-orphans").Value.Set(strconv.FormatBool(removeOrphans))
	command.Flags().Lookup("volumes").DefValue = strconv.FormatBool(volumes)
	command.Flags().Lookup("volumes").Value.Set(strconv.FormatBool(volumes))
	// Create root command
	var cmd *cobra.Command = &cobra.Command{Use: "compose"}
	cmd.SetArgs([]string{"down"})
	cmd.AddCommand(command)
	// Run command
	return cmd.Execute()
}
