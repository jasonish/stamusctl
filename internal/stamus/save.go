package stamus

import (
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/models"
)

func SaveLogin(registryInfo models.RegistryInfo) error {
	// Get config content
	Config, err := GetStamusConfig()
	if err != nil {
		return err
	}

	// Save in struct
	Config.SetRegistry(
		Registry(registryInfo.Registry),
		User(registryInfo.Username),
		Token(registryInfo.Password),
	)

	// Save config
	Config.setStamusConfig()

	return nil
}

func SetCurrent(name string) error {
	// Get config content
	Config, err := GetStamusConfig()
	if err != nil {
		return err
	}

	// Save in struct
	err = Config.SetCurrent(name)
	if err != nil {
		return err
	}

	// Save config
	return Config.setStamusConfig()
}

func GetCurrent() (string, error) {
	// Get config content
	config, err := GetStamusConfig()
	if err != nil {
		return "", err
	}
	// Check current
	current := config.Current
	if current == "" {
		current = app.DefaultConfigName
	}
	// Get current
	return current, nil
}
