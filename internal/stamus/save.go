package stamus

func SaveLogin(registry string, user string, token string) error {
	// Get config content
	Config, err := GetStamusConfig()
	if err != nil {
		return err
	}

	// Save in struct
	Config.SetRegistry(Registry(registry), User(user), Token(token))

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
