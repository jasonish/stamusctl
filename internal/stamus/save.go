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

func SaveConfig(name string) error {
	// Get config content
	Config, err := GetStamusConfig()
	if err != nil {
		return err
	}

	// Save in struct
	Config.SetConfig(name)

	// Save config
	Config.setStamusConfig()

	return nil
}
