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
	f, err := os.OpenFile(filepath.Join(app.ConfigFolder, "config.json"), os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func GetStamusConfig() (*Config, error) {
	// Open or create ~/stamus/config.json
	file, err := getOrCreateStamusConfigFile()
	if err != nil {
		return nil, err
	}
	// Read the file contents
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Unmarshal the file contents
	Config := &Config{}
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &Config)
		if err != nil {
			return nil, err
		}
	}

	return Config, nil
}
