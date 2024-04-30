package compose

import (
	// Common
	"fmt"

	// External
	"github.com/spf13/cobra"

	// Custom
	compose "stamus-ctl/internal/docker-compose"
	"stamus-ctl/internal/logging"
)

var (
	format = ""
)

func getConfig() {
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

func ConfigHandler() *cobra.Command {
	var command = &cobra.Command{
		Use:   "config",
		Short: "Interact with container compose config file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			getConfig()
		},
	}

	var subCommands = []*cobra.Command{
		{
			Use:   "get",
			Short: "Get container compose config file",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {
				getConfig()
			},
		},
		{
			Use:   "set",
			Short: "Set container compose config file",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {
				logging.Sugar.Errorw("compose config set not yet implemented")
			},
		},
	}

	command.Flags().StringVarP(&format, "format", "f", "", "format")

	for _, subCommand := range subCommands {
		command.AddCommand(subCommand)
	}
	return command
}

func init() {

}
