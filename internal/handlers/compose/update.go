package handlers

import (
	// Common

	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// External

	// Custom
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/stamus"
	"stamus-ctl/internal/utils"
)

type UpdateHandlerParams struct {
	Config  string
	Args    []string
	Version string
}

func UpdateHandler(params UpdateHandlerParams) error {
	// Unpack params
	configPath := params.Config
	args := params.Args
	versionVal := params.Version

	// Get registry info
	image := "/selks:" + versionVal
	destPath := filepath.Join(app.TemplatesFolder + "selks/")
	latestPath := filepath.Join(destPath, "latest/")

	// Get registries infos
	stamusConf, err := stamus.GetStamusConfig()
	if err != nil {
		return err
	}

	// Pull config
	fmt.Println("Getting configuration")
	for _, registryInfo := range stamusConf.Registries.AsList() {
		err = registryInfo.PullConfig(destPath, image)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	// Execute update script
	prerunPath := filepath.Join(destPath, "sbin/pre-run")
	postrunPath := filepath.Join(destPath, "sbin/post-run")
	runOutput, err := runArbitraryScript(prerunPath)
	if err != nil {
		return err
	}

	// Save output
	outputFile, err := os.Create(filepath.Join(configPath, "values.yaml"))
	if err != nil {
		return err
	}
	defer outputFile.Close()
	if _, err := outputFile.WriteString(runOutput.String()); err != nil {
		return err
	}

	// Load existing config
	confFile, err := models.CreateFileInstance(configPath, "values.yaml")
	if err != nil {
		return err
	}
	existingConfig, err := models.LoadConfigFrom(confFile, false)
	if err != nil {
		return err
	}

	// Create new config
	newConfFile, err := models.CreateFileInstance(latestPath, "config.yaml")
	if err != nil {
		return err
	}
	newConfig, err := models.NewConfigFrom(newConfFile)
	if err != nil {
		return err
	}
	_, _, err = newConfig.ExtractParams(true)
	if err != nil {
		return err
	}

	// Extract and set values from args and existing config
	paramsArgs := utils.ExtractArgs(args)
	newConfig.GetParams().SetValues(existingConfig.GetParams().GetVariablesValues())
	newConfig.GetParams().SetLooseValues(paramsArgs)
	newConfig.GetArbitrary().SetArbitrary(paramsArgs)
	newConfig.GetParams().ProcessOptionnalParams(false)

	// Ask for missing parameters
	err = newConfig.GetParams().AskMissing()
	if err != nil {
		return err
	}

	// Save the configuration
	err = newConfig.SaveConfigTo(confFile)
	if err != nil {
		return err
	}

	// Run post-run script
	runOutput, err = runArbitraryScript(postrunPath)
	if err != nil {
		return err
	}
	fmt.Println("")

	return nil
}

func runArbitraryScript(path string) (*strings.Builder, error) {
	// Execute arbitrary script
	arbitrary := exec.Command(path)
	// Display output to terminal
	runOutput := new(strings.Builder)
	arbitrary.Stdout = runOutput
	arbitrary.Stderr = os.Stderr
	// Change execution rights
	os.Chmod(path, 0755)
	// Run arbitrary script
	if err := arbitrary.Run(); err != nil {
		return nil, err
	}
	return runOutput, nil
}
