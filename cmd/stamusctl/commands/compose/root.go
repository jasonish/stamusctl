package compose

import (
	"github.com/spf13/cobra"
)

var (
	nonInteractive = false
	netInterface   = ""
	restart        = ""
	elasticPath    = ""
	dataPath       = ""
)

func NewCompose() *cobra.Command {
	var command = &cobra.Command{
		Use:   "compose",
		Short: "work with docker-compose",
	}
	command.AddCommand(NewInit())
	command.AddCommand(NewTemplate())
	return command
}

func init() {

}
