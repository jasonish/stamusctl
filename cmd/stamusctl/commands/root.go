package cmd

import (
	"fmt"
	"os"
	"strings"

	"git.stamus-networks.com/lanath/stamus-ctl/cmd/stamusctl/commands/compose"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/app"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "stamusctl",
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntP("verbose", "v", 0, "verbose level")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(compose.NewCompose())
}

func initConfig() {
	viper.SetEnvPrefix(app.Name)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	logging.SetLogger()
}
