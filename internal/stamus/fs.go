package stamus

import (
	// Common
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	// Custom
	"stamus-ctl/internal/app"
)

func getOrCreateStamusConfigFile() (*os.File, error) {
	// Create ~/stamus directory
	err := os.MkdirAll(app.ConfigFolder, 0755)
	if err != nil {
		return nil, err
	}

	// Open or create ~/stamus/config.json
	f, err := os.OpenFile(filepath.Join(app.ConfigFolder, "config.json"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func tryGetStamusConfigFile() (*os.File, error) {
	// Open or create ~/stamus/config.json
	f, err := os.OpenFile(filepath.Join(app.ConfigFolder, "config.json"), os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func GetStamusConfig() (*Config, error) {
	// Open or create ~/stamus/config.json
	file, err := tryGetStamusConfigFile()
	if err != nil {
		return &Config{}, nil
	}
	// Read the file contents
	bytes, err := io.ReadAll(file)
	if err != nil {
		return &Config{}, nil
	}
	// Unmarshal the file contents
	config := &Config{}
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return &Config{}, nil
		}
	}

	return config, nil
}
