package handlers

import (
	"stamus-ctl/internal/models"
)

type InitHandlerInputs struct {
	IsDefault        bool
	FolderPath       string
	BackupFolderPath string
	OutputPath       string
	Arbitrary        map[string]string
}

func InitHandler(cli bool, params InitHandlerInputs) error {
	// Instanciate config
	var config *models.Config
	confFile, err := models.CreateFileInstance(params.FolderPath, "config.yaml")
	if err != nil {
		confFile, err = models.CreateFileInstance(params.BackupFolderPath, "config.yaml")
		if err != nil {
			return err
		}
	}
	config, err = models.NewConfigFrom(confFile)
	if err != nil {
		return err
	}
	// Read the folder configuration
	_, _, err = config.ExtractParams(true)
	if err != nil {
		return err
	}
	// Set parameters
	if params.IsDefault {
		// Extract and set values from args
		err = config.GetParams().SetLooseValues(params.Arbitrary)
		config.SetArbitrary(params.Arbitrary)
		if err != nil {
			return err
		}
		// Set from default
		err := config.GetParams().SetToDefaults()
		if err != nil {
			return err
		}
		// Ask for missing parameters
		if cli {
			err = config.GetParams().AskMissing()
			if err != nil {
				return err
			}
		}
	} else {
		// Extract and set values from args
		err = config.GetParams().SetLooseValues(params.Arbitrary)
		config.SetArbitrary(params.Arbitrary)
		if err != nil {
			return err
		}
		//Set from user input
		if cli {
			err := config.GetParams().AskAll()
			if err != nil {
				return err
			}
		}
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
