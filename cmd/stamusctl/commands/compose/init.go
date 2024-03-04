package compose

import (
	"errors"
	"os"

	compose "git.stamus-networks.com/lanath/stamus-ctl/internal/docker-compose"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	nonInteractive = false
	v              = viper.New()

	params compose.Parameters
)

func NewInit() *cobra.Command {
	var command = &cobra.Command{
		Use:   "init",
		Short: "create docker compose file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return compose.ValidateInputFlag(params)
		},
		Run: func(cmd *cobra.Command, args []string) {

			f, err := os.OpenFile(params.OutputFile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logging.Sugar.Fatal("cannot create docker file", "error", err)
			}

			defer f.Close()

			_, err = os.OpenFile(InputFileConfigName, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logging.Sugar.Fatalw("cannot create config file", "error", err)
			}

			manifest := compose.GenerateComposeFileFromCli(cmd, &params, nonInteractive)

			v.Set("suricata.interfaces", params.InterfacesList)
			v.Set("scirius.token", params.SciriusToken)
			v.Set("scirius.debugMode", params.DebugMode)
			v.Set("scirius.version", params.SciriusVersion)
			v.Set("arkimeviewer.version", params.ArkimeviewerVersion)
			v.Set("elk.version", params.ElkVersion)
			v.Set("elk.elastic.path", params.ElasticPath)
			v.Set("elk.elastic.memory", params.ElasticMemory)
			v.Set("elk.elastic.ml", params.MLEnabled)
			v.Set("elk.logstash.memory", params.LogstashMemory)
			v.Set("global.volumes.path", params.VolumeDataPath)
			v.Set("global.restartMode", params.RestartMode)
			v.Set("global.registry", params.Registry)
			v.Set("config.outputFile", params.OutputFile)
			v.Set("nginx.exec", params.NginxExec)

			err = v.WriteConfig()
			if err != nil {
				logging.Sugar.Fatalw("cannot write config file", "error", err)
			}
			if _, err := os.Stat(params.VolumeDataPath + "/nginx/ssl"); errors.Is(err, os.ErrNotExist) {
				compose.GenerateSSLWithDocker(params.VolumeDataPath + "/nginx/ssl")
			} else {
				logging.Sugar.Debugw("cert already exist. skiped.", "path", params.VolumeDataPath+"/nginx/ssl")
			}

			f.WriteString(manifest)
			compose.WriteConfigFiles(params.VolumeDataPath)

		},
	}
	command.Flags().StringVarP(&params.OutputFile, "output", "o", "docker-compose.yaml", "Defines the path where to write the docker-compose file.")
	command.PersistentFlags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "set interactive mode.")

	command.PersistentFlags().StringVarP(&params.InterfacesList, "interface", "i", "", "Defines an interface on which SELKS should listen.")
	command.PersistentFlags().StringVarP(&params.SciriusToken, "token", "t", "", "Scirius secret key.")

	command.PersistentFlags().StringVar(&params.VolumeDataPath, "container-datapath", utils.IgnoreError(os.Getwd())+"/containers-data", "Defines the path where SELKS will store it's data.")
	command.PersistentFlags().StringVar(&params.Registry, "registry", "", "Defines the path where SELKS will store it's data.")

	command.PersistentFlags().StringVar(&params.SciriusVersion, "scirius-version", "master", "Defines the version of the scirius to use.")
	command.PersistentFlags().StringVar(&params.ArkimeviewerVersion, "arkimeviewer-version", "master", "Defines the version of arkimeviewer to use.")
	command.PersistentFlags().StringVar(&params.ElkVersion, "elk-version", "7.16.1", "Defines the version of the ELK stack to use.")

	command.PersistentFlags().StringVar(&params.ElasticPath, "es-datapath", "/var/lib/docker", "Defines the path where Elasticsearch will store it's data.")

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

	command.AddCommand(NewTemplate())

	return command
}

func init() {

}
