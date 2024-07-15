package handlers

import (
	"path/filepath"
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/stamus"
)

type InitHandlerInputs struct {
	IsDefault        bool
	FolderPath       string
	BackupFolderPath string
	OutputPath       string
	Project          string
	Version          string
	Arbitrary        map[string]string
}

func InitHandler(isCli bool, params InitHandlerInputs) error {
	// Get registry info
	image := "/" + params.Project + ":" + params.Version
	destPath := filepath.Join(app.TemplatesFolder, params.Project)
	// Pull latest template
	err := pullLatestTemplate(destPath, image)
	if err != nil {
		return err
	}
	// Instanciate config
	config, err := instanciateConfig(destPath, params.BackupFolderPath)
	if err != nil {
		return err
	}
	// Read the folder configuration
	_, _, err = config.ExtractParams(true)
	if err != nil {
		return err
	}
	// Set parameters
	err = setParameters(isCli, config, params)
	if err != nil {
		return err
	}
	// Validate parameters
	err = config.GetParams().ValidateAll()
	if err != nil {
		return err
	}
	// Save the configuration
	outputFile, err := models.CreateFileInstance(params.OutputPath, "values.yaml")
	if err != nil {
		return err
	}
	config.SaveConfigTo(outputFile)
	return nil
}

// Pull latest template from saved registries
func pullLatestTemplate(destPath string, image string) error {
	// Get registries infos
	stamusConf, err := stamus.GetStamusConfig()
	if err != nil {
		return err
	}
	// Pull latest config
	for _, registryInfo := range stamusConf.Registries.AsList() {
		err = registryInfo.PullConfig(destPath, image)
		if err == nil {
			break
		}
	}
	return err
}

// Instanciate config from folder or backup folders
func instanciateConfig(folderPath string, backupFolderPath string) (*models.Config, error) {
	var config *models.Config
	confFile, err := models.CreateFileInstance(folderPath, "config.yaml")
	if err != nil {
		confFile, err = models.CreateFileInstance(backupFolderPath, "config.yaml")
		if err != nil {
			return nil, err
		}
	}
	config, err = models.NewConfigFrom(confFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Set parameters, from args, and defaults / asks rest
func setParameters(isCli bool, config *models.Config, params InitHandlerInputs) error {
	// Extract and set values from args
	err := config.GetParams().SetLooseValues(params.Arbitrary)
	config.SetArbitrary(params.Arbitrary)
	if err != nil {
		return err
	}
	// Set from default
	if params.IsDefault {
		err = config.GetParams().SetToDefaults()
		if err != nil {
			return err
		}
	}
	// Ask for missing parameters
	if isCli {
		err = config.GetParams().AskMissing()
		if err != nil {
			return err
		}
	}
	return nil
}
