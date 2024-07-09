package models

import (
	"fmt"
	"stamus-ctl/internal/app"
	"strings"

	"github.com/spf13/viper"
)

func InstanciateViper(file File) (*viper.Viper, error) {
	// Create a new viper instance
	viperInstance := viper.New()
	// General configuration
	viperInstance.SetEnvPrefix(app.Name)
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viperInstance.AutomaticEnv()
	// Specific configuration
	viperInstance.SetConfigName(file.Name)
	viperInstance.SetConfigType(file.Type)
	viperInstance.AddConfigPath(file.Path)
	// Read the config file
	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	return viperInstance, nil
}
