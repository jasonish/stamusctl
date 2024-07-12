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
	prerun := exec.Command(prerunPath)
	postrun := exec.Command(postrunPath)
	// Display output to terminal
	runOutput := new(strings.Builder)
	prerun.Stdout = runOutput
	prerun.Stderr = os.Stderr
	// Change execution rights
	os.Chmod(prerunPath, 0755)
	os.Chmod(postrunPath, 0755)
	// Run pre-run script
	if err := prerun.Run(); err != nil {
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
	newConfig.SetArbitrary(paramsArgs)
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
	postrunOutput := new(strings.Builder)
	postrun.Stdout = postrunOutput
	postrun.Stderr = os.Stderr
	// Run pre-run script
	if err := postrun.Run(); err != nil {
		return err
	}
	fmt.Println("")

	return nil
}
