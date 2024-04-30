package compose

import (
	"fmt"

	compose "stamus-ctl/internal/docker-compose"

	"github.com/spf13/cobra"
)

func NewTemplate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "template",
		Short: "create docker compose file and output it to stdout",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return compose.ValidateInputFlag(params)
		},
		Run: func(cmd *cobra.Command, args []string) {
			manifest := compose.GenerateComposeFileFromCli(cmd, &params, nonInteractive)

			fmt.Print(manifest)
		},
	}

	return command
}

func init() {

}
