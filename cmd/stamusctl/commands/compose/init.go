package compose

import (
	"os"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/spf13/cobra"
)

var (
	outputFile string
)

func NewInit() *cobra.Command {
	var command = &cobra.Command{
		Use:   "init",
		Short: "create docker compose file",
		PreRun: func(cmd *cobra.Command, args []string) {
			if restart != "no" &&
				restart != "always" &&
				restart != "on-failure" &&
				restart != "unless-stopped" {
				logging.Sugar.Fatalf("Please provid a valid value for --restart. %s is not valid", restart)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logging.Sugar.Fatal("cannot create output file")
			}

			defer f.Close()

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
			)

			f.WriteString(manifest)

		},
	}
	command.Flags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "set interactive mode.")
	command.Flags().StringVarP(&netInterface, "interface", "i", "", "Defines an interface on which SELKS should listen.")
	command.Flags().StringVar(&elasticPath, "es-datapath", "/var/lib/docker", "Defines the path where Elasticsearch will store it's data.")
	command.Flags().StringVar(&dataPath, "container-datapath", "", "Defines the path where SELKS will store it's data.")
	command.Flags().StringVarP(&outputFile, "output", "o", "docker-compose.yaml", "Defines the path where SELKS will store it's data.")
	command.Flags().StringVarP(&restart, "restart", "r", "unless-stopped",
		`restart mode.
'no': never restart automatically the containers
'always': automatically restart the containers even if they have been manually stopped
'on-failure': only restart the containers if they failed
'unless-stopped': always restart the container except if it has been manually stopped`,
	)
	return command
}

func init() {

}
