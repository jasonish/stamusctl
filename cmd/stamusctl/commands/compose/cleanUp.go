package compose

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	print bool
)

func NewCleanup() *cobra.Command {
	var command = &cobra.Command{
		Use:   "clean",
		Short: "clean docker compose file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			os.Remove("docker-compose.yaml")
			os.Remove("containers-data")
		},
	}
	command.PersistentFlags().BoolVarP(&print, "print", "p", false, "Print before deleting.")

	return command
}

func init() {

}
