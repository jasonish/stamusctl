package compose

import (
	"github.com/spf13/cobra"
)

func NewCompose() *cobra.Command {
	var command = &cobra.Command{
		Use:   "compose",
		Short: "work with docker-compose",
	}
	initCommand := NewInit()

	command.AddCommand(initCommand)
	command.AddCommand(NewCleanup())

	initCommand.AddCommand(NewTemplate())

	return command
}

func init() {

}
