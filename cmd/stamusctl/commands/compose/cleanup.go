package compose

import (
	// Common
	"fmt"
	"os"

	// External
	"github.com/spf13/cobra"

	// Custom
	compose "stamus-ctl/internal/docker-compose"
	utils "stamus-ctl/internal/utils"
)

var (
	print bool
)

func CleanupHandler() *cobra.Command {
	var command = &cobra.Command{
		Use:   "clean",
		Short: "Clean container compose file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Confirm
			if !utils.AskForConfirmation("Are you sure you want to delete the compose file and associated volumes ? (y/n)") {
				return
			}
			// Get parameters
			p := compose.NewParametersFromEnv(v)
			// Delete files
			if print {
				fmt.Printf("Deleting %s\n", p.OutputFile)
			}
			os.Remove(p.OutputFile)
			// Delete volumes
			if print {
				fmt.Printf("Deleting %s\n", p.VolumeDataPath)
			}
			os.RemoveAll(p.VolumeDataPath)
		},
	}

	command.PersistentFlags().BoolVarP(&print, "print", "p", false, "Print before deleting.")

	return command
}

func init() {

}
