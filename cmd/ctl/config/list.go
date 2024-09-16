package config

import (
	// Core
	"fmt"
	"log"

	// External
	"github.com/spf13/cobra"

	// Internal

	"stamus-ctl/internal/stamus"
)

// Command
func listCmd() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get list of configurations on the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listHandler()
		},
	}
	return cmd
}

// Handlers
func listHandler() error {
	// Get list
	list, err := stamus.GetConfigsList()
	if err != nil {
		log.Println("Error getting list of configurations")
		return err
	}
	current, err := stamus.GetCurrent()
	if err != nil {
		log.Println("Error getting current configuration")
		return err
	}
	fmt.Println(" List of configurations")
	// Print list
	for _, config := range list {
		if current == config {
			fmt.Println(" - " + config + " (current)")
		} else {
			fmt.Println(" - " + config)
		}
	}
	return nil
}
