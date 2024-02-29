package compose

import (
	"fmt"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
	"github.com/spf13/cobra"
)

func NewTemplate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "template",
		Short: "create docker compose file and output it to stdout",
		PreRun: func(cmd *cobra.Command, args []string) {
			compose.ValidateInputFlag(params)
		},
		Run: func(cmd *cobra.Command, args []string) {
			manifest := compose.GenerateComposeFileFromCli(cmd, params, nonInteractive)

			fmt.Print(manifest)
		},
	}

	return command
}

func init() {

}
