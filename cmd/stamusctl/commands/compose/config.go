package compose

import (
	"fmt"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
	"github.com/spf13/cobra"
)

var (
	format = ""
)

func run() {
	p := compose.NewParametersFromEnv(v)
	if format != "" {
		out := p.Format(format)
		fmt.Print(out)
	} else {
		p.Logs(func(s string) string {
			fmt.Print(s)
			return s
		})
	}
}

func NewGetConfig() *cobra.Command {
	var command = &cobra.Command{
		Use:   "config",
		Short: "clean docker compose file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	get := &cobra.Command{
		Use:   "config",
		Short: "clean docker compose file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	command.Flags().StringVarP(&format, "format", "f", "", "format")

	command.AddCommand(get)
	return command
}

func init() {

}
