package compose

import (
	"fmt"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/spf13/cobra"
)

func NewTemplate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "template",
		Short: "create docker compose file and output it to stdout",
		PreRun: func(cmd *cobra.Command, args []string) {
			if restart != "no" &&
				restart != "always" &&
				restart != "on-failure" &&
				restart != "unless-stopped" {
				logging.Sugar.Fatalf("Please provid a valid value for --restart. %s is not valid", restart)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := compose.CheckVersions(); err != nil {
				logging.Sugar.Fatal(err.Error())
			}

			if nonInteractive == false {
				compose.Ask(
					cmd,
					&netInterface,
					&restart,
					&elasticPath,
					&dataPath,
					&registry,
					&token,
				)
			}

			if netInterface == "" {
				logging.Sugar.Fatal("please provide a valid network interface.")
			}

			manifest := compose.GenerateComposeFile(
				netInterface,
				restart,
				elasticPath,
				dataPath,
				registry,
				token,
				elkVersion,
			)

			fmt.Print(manifest)
		},
	}

	return command
}

func init() {

}
