package compose

import (
	"os"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/spf13/cobra"
)

var (
	nonInteractive = false
	outputFile     string

	params compose.Parameters
)

func NewInit() *cobra.Command {
	var command = &cobra.Command{
		Use:   "init",
		Short: "create docker compose file",
		PreRun: func(cmd *cobra.Command, args []string) {
			if params.RestartMode != "no" &&
				params.RestartMode != "always" &&
				params.RestartMode != "on-failure" &&
				params.RestartMode != "unless-stopped" {
				logging.Sugar.Fatalf("Please provid a valid value for --restart. %s is not valid", params.RestartMode)
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
					&params,
				)
			}

			if params.InterfacesList == "" {
				logging.Sugar.Fatal("please provide a valid network interface.")
			}

			manifest := compose.GenerateComposeFile(params)

			f.WriteString(manifest)

		},
	}
	command.Flags().StringVarP(&outputFile, "output", "o", "docker-compose.yaml", "Defines the path where SELKS will store it's data.")
	command.PersistentFlags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "set interactive mode.")

	command.PersistentFlags().StringVarP(&params.InterfacesList, "interface", "i", "", "Defines an interface on which SELKS should listen.")
	command.PersistentFlags().StringVarP(&params.SciriusToken, "token", "t", "", "Scirius secret key.")

	command.PersistentFlags().StringVar(&params.VolumeDataPath, "container-datapath", "", "Defines the path where SELKS will store it's data.")
	command.PersistentFlags().StringVar(&params.Registry, "registry", "", "Defines the path where SELKS will store it's data.")

	command.PersistentFlags().StringVar(&params.ElasticPath, "es-datapath", "ghcr.io/stamusnetworks", "Defines the path where Elasticsearch will store it's data.")
	command.PersistentFlags().StringVar(&params.ElasticMemory, "es-memory", "3G", "Defines the path where Elasticsearch will store it's data.")
	command.PersistentFlags().StringVar(&params.ElkVersion, "elk-version", "7.16.1", "Defines the version of the ELK stack to use.")

	command.PersistentFlags().StringVarP(&params.RestartMode, "restart", "r", "unless-stopped",
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
