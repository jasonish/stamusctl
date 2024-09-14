package stamus

import "stamus-ctl/internal/models"

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
