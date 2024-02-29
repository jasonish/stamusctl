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
			compose.ValidateInputFlag(params)
		},
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logging.Sugar.Fatal("cannot create output file")
			}

			defer f.Close()

			manifest := compose.GenerateComposeFileFromCli(cmd, params, nonInteractive)

			f.WriteString(manifest)

		},
	}
	command.Flags().StringVarP(&outputFile, "output", "o", "docker-compose.yaml", "Defines the path where SELKS will store it's data.")
	command.PersistentFlags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "set interactive mode.")

	command.PersistentFlags().StringVarP(&params.InterfacesList, "interface", "i", "", "Defines an interface on which SELKS should listen.")
	command.PersistentFlags().StringVarP(&params.SciriusToken, "token", "t", "", "Scirius secret key.")

	command.PersistentFlags().StringVar(&params.VolumeDataPath, "container-datapath", "", "Defines the path where SELKS will store it's data.")
	command.PersistentFlags().StringVar(&params.Registry, "registry", "", "Defines the path where SELKS will store it's data.")

	command.PersistentFlags().StringVar(&params.SciriusToken, "scirius-version", "master", "Defines the version of the scirius to use.")
	command.PersistentFlags().StringVar(&params.ArkimeviewerVersion, "arkimeviewer-version", "master", "Defines the version of arkimeviewer to use.")
	command.PersistentFlags().StringVar(&params.ElkVersion, "elk-version", "7.16.1", "Defines the version of the ELK stack to use.")

	command.PersistentFlags().StringVar(&params.ElasticPath, "es-datapath", "ghcr.io/stamusnetworks", "Defines the path where Elasticsearch will store it's data.")

	command.PersistentFlags().StringVar(&params.ElasticMemory, "es-memory", "3G", "Amount of memory to give to the elasticsearch container.")
	command.PersistentFlags().StringVar(&params.LogstashMemory, "ls-memory", "2G", "Amount of memory to give to the logstash container.")

	command.PersistentFlags().StringVarP(&params.RestartMode, "restart", "r", "unless-stopped",
		`restart mode.
'no': never restart automatically the containers
'always': automatically restart the containers even if they have been manually stopped
'on-failure': only restart the containers if they failed
'unless-stopped': always restart the container except if it has been manually stopped`,
	)
	command.PersistentFlags().BoolVarP(&params.DebugMode, "debug", "d", false, "Activate debug mode for scirius and nginx.")

	return command
}

func init() {

}
