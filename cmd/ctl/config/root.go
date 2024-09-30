package config

import (
	// Common

	// External

	"github.com/spf13/cobra"
	// Custom
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/embeds"
	"stamus-ctl/internal/utils"
)

// Init
func init() {
	// Setup
	initSelksFolder(app.DefaultSelksPath)
}

// Create SELKS folder if it does not exist
func initSelksFolder(path string) {
	selksConfigExist, err := utils.FolderExists(path)
	if err != nil {
		panic(err)
	}
	if !selksConfigExist && app.Embed.IsTrue() {
		err = embeds.ExtractEmbedTo("selks", app.TemplatesFolder+"selks/embedded/")
		if err != nil {
			panic(err)
		}
	}
}

// Commands
func ConfigCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Interact with compose config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHandler(cmd, args)
		},
	}
	// Add Commands
	cmd.AddCommand(getCmd())
	cmd.AddCommand(setCmd())
	return cmd
}
