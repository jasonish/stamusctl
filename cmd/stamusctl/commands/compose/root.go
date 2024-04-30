package compose

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	InputFileConfigName = ".config.yaml"
	InputFileConfigType = "yaml"
	InputFileConfigPath = "."
)

func NewCompose() *cobra.Command {
	var command = &cobra.Command{
		Use:   "compose",
		Short: "Interact with container compose",
	}

	command.AddCommand(NewInit())
	command.AddCommand(CleanupHandler())
	command.AddCommand(ConfigHandler())

	cobra.OnInitialize(initConfig)

	return command
}

func initConfig() {

	v.SetConfigName(InputFileConfigName)
	v.SetConfigType(InputFileConfigType)
	v.AddConfigPath(InputFileConfigPath)
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

}

func init() {

}
