package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	compose "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
	"strings"

	cp "github.com/otiai10/copy"
)

type SetHandlerInputs struct {
	Config   string   // Path to the config folder
	Values   string   // Path to the values.yaml file
	Reload   bool     // Reload the configuration, don't keep arbitrary parameters
	Apply    bool     // Apply the new configuration, reload the services
	Args     []string // Cmd arguments
	FromFile string   // Path to the file containing the values
}

// func SetHandler(configPath string, args []string, reload bool, apply bool) error {
func SetHandler(params SetHandlerInputs) error {
	// Load the config
	file, err := models.CreateFileInstance(params.Config, "values.yaml")
	if err != nil {
		return err
	}
	config, err := models.LoadConfigFrom(file, params.Reload)
	if err != nil {
		return err
	}

	// Extract and set parameters from args
	paramsArgs := utils.ExtractArgs(params.Args)
	config.GetParams().SetLooseValues(paramsArgs)
	config.SetArbitrary(paramsArgs)
	config.GetParams().ProcessOptionnalParams(false)
	// Set values from file
	err = config.GetParams().SetValuesFromFiles(params.FromFile)
	if err != nil {
		return err
	}
	err = config.GetParams().SetValuesFromFile(params.Values)
	if err != nil {
		return err
	}
	// Validate
	err = config.GetParams().ValidateAll()
	if err != nil {
		return err
	}

	// Save the configuration
	outputAsFile, err := models.CreateFileInstance(params.Config, "values.yaml")
	if err != nil {
		return err
	}
	config.SaveConfigTo(outputAsFile)
	// Apply the configuration
	if params.Apply {
		err = compose.HandleUp(params.Config)
		if err != nil {
			return err
		}
	}
	return nil
}

// For each argument, copy the input path to the output path
func SetContentHandler(path string, args []string) error {
	// For each argument
	for _, arg := range args {
		// Split argument
		split := strings.Split(arg, ":")
		if len(split) != 2 {
			return fmt.Errorf("Invalid argument: %s", arg)
		}
		// Get paths
		inputPath := split[0]
		outputPath := split[1]
		// Call handler
		err := copy(inputPath, filepath.Join(path, outputPath))
		if err != nil {
			return err
		}
	}
	return nil
}

func copy(inputPath string, outputPath string) error {
	fmt.Println("Setting content from ", inputPath, " to ", outputPath)
	// Check input path exists
	info, err := os.Stat(inputPath)
	if err != nil {
		log.Println(info, err)
		return fmt.Errorf("Input path does not exist: %s", inputPath)
	}

	err = cp.Copy(inputPath, outputPath)
	if err != nil {
		return err
	}

	return nil
}
