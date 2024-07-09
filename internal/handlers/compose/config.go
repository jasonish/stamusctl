package handlers

import (
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
	"strings"
)

func SetHandler(configPath string, args []string, reload bool) error {
	// Load the config
	file, err := models.CreateFileInstance(configPath, "values.yaml")
	if err != nil {
		return err
	}
	config, err := models.LoadConfigFrom(file, reload)
	if err != nil {
		return err
	}
	// Extract and set parameters from args
	paramsArgs := utils.ExtractArgs(args)
	config.GetParams().SetLooseValues(paramsArgs)
	config.SetArbitrary(paramsArgs)
	config.GetParams().ProcessOptionnalParams(false)
	// Validate
	err = config.GetParams().ValidateAll()
	if err != nil {
		return err
	}
	// Save the configuration
	outputAsFile, err := models.CreateFileInstance(configPath, "values.yaml")
	if err != nil {
		return err
	}
	config.SaveConfigTo(outputAsFile)
	return nil
}

// Get the grouped config values
// Essentially, this function reads the config values file and groups the values
func GetGroupedConfig(configPath string, args []string, reload bool) (map[string]interface{}, error) {
	// File instance
	inputAsFile, err := models.CreateFileInstance(configPath, "values.yaml")
	if err != nil {
		return nil, err
	}
	// Load the config
	config, err := models.LoadConfigFrom(inputAsFile, reload)
	if err != nil {
		return nil, err
	}
	// Process optionnal parameters
	err = config.GetParams().ProcessOptionnalParams(false)
	if err != nil {
		return nil, err
	}
	// Group values
	groupedValues := groupValues(config, args)
	// Return
	return groupedValues, nil
}

// Group values
// Utility function to group values from the config to nested maps
func groupValues(config *models.Config, args []string) map[string]interface{} {
	values := config.GetParams().GetValues(args...)
	groupedValues := make(map[string]interface{})
	for _, param := range config.GetParams().GetOrdered() {
		if value, ok := values[param]; ok {
			parts := strings.Split(param, ".")
			addToGroup(parts, value, groupedValues)
		}
	}
	return groupedValues
}

func addToGroup(parts []string, value string, group map[string]interface{}) {
	if len(parts) == 1 {
		group[parts[0]] = value
	} else {
		if _, ok := group[parts[0]]; !ok {
			group[parts[0]] = make(map[string]interface{})
		}
		addToGroup(parts[1:], value, group[parts[0]].(map[string]interface{}))
	}
}
