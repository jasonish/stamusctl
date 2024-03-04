package compose

import (
	"fmt"
	"os"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
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
			p := compose.NewParametersFromEnv(v)

			if print {
				fmt.Printf("Deleting %s\n", p.OutputFile)
			}
			os.Remove(p.OutputFile)
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
