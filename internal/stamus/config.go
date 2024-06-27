package stamus

import (
	"encoding/json"
)

type Registry string
type User string
type Token string

type Registries map[Registry]Logins
type Logins map[User]Token

type Config struct {
	Registries Registries `json:"registries"`
	Configs    []string   `json:"configs"`
	Current    string     `json:"current"`
}

func (conf Config) setStamusConfig() error {
	// Open or create ~/.stamus/config.json
	file, err := getOrCreateStamusConfigFile()
	if err != nil {
		return err
	}

	// Write the new content
	bytes, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	_, err = file.WriteAt(bytes, 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SetRegistry(registry Registry, user User, token Token) {
	// Create Registries if it does not exist in Config
	if c.Registries == nil {
		c.Registries = make(Registries)
	}
	// Create Registry if it does not exist in Registries
	if c.Registries[Registry(registry)] == nil {
		c.Registries[Registry(registry)] = make(Logins)
	}

	c.Registries[Registry(registry)][User(user)] = Token(token)
}

func (c *Config) SetConfig(config string) {
	c.Configs = append(c.Configs, config)
	uniqueConfigs := make([]string, 0, len(c.Configs))
	seen := make(map[string]bool)
	for _, c := range c.Configs {
		if !seen[c] {
			uniqueConfigs = append(uniqueConfigs, c)
			seen[c] = true
		}
	}
	c.Configs = uniqueConfigs
}
