package models

import (
	// Common

	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	// External

	cp "github.com/otiai10/copy"
	"github.com/spf13/viper"
	// Custom
)

const defaultConfPath = ".configs/selks/embedded/"

type Config struct {
	file          file
	parameters    *Parameters
	viperInstance *viper.Viper
}

func NewConfigFrom(file file) (*Config, error) {
	// Instanciate viper
	viperInstance, err := instanciateViper(file)
	if err != nil {
		return nil, err
	}
	// Create the config
	conf := Config{
		file:          file,
		viperInstance: viperInstance,
	}
	return &conf, nil
}

func LoadConfigFrom(path file) (*Config, error) {
	// Load the config
	configured, err := NewConfigFrom(path)
	if err != nil {
		return nil, err
	}
	// Extract config data
	values := configured.ExtractValues()
	stamusConfPathPointer := values["stamusconfig"]
	stamusConfPath := *stamusConfPathPointer.String
	file, err := CreateFileInstance(stamusConfPath, "config.yaml")
	if err != nil {
		return nil, err
	}
	// Load origin config
	originConf, err := NewConfigFrom(file)
	if err != nil {
		return nil, err
	}
	_, _, err = originConf.ExtractParams(true)
	if err != nil {
		return nil, err
	}
	// Merge
	originConf.parameters.SetValues(values)
	originConf.parameters.ProcessOptionnalParams(false)
	return originConf, nil
}

// Return list of config files to include and list of parameters for current config
func (f *Config) extracKeys() ([]string, []string) {
	// Extract includes list
	includes := []string{}
	// Fron viperInstance, get the key "includes" as string list
	includes = f.viperInstance.GetStringSlice("includes")
	// Extract parameters list
	parametersMap := map[string]bool{}
	for _, key := range f.viperInstance.AllKeys() {
		// Extract the parameter name
		parameterAsArray := strings.Split(key, ".")
		parameter := strings.Join(parameterAsArray[:len(parameterAsArray)-1], ".")
		if len(parameter) != 0 {
			parametersMap[parameter] = true
		}
	}
	// Convert map to list
	parametersList := []string{}
	for key, _ := range parametersMap {
		parametersList = append(parametersList, key)
	}
	return includes, parametersList
}

// Returns the parameter extracted from the config file
func (f *Config) extractParam(parameter string, isDeep bool) (*Parameter, error) {
	// Extract the parameter
	currentParam := Parameter{
		Name:         f.getStringParamValue(parameter, "name"),
		Shorthand:    f.getStringParamValue(parameter, "shorthand"),
		Type:         f.getStringParamValue(parameter, "type"),
		Usage:        f.getStringParamValue(parameter, "usage"),
		ValidateFunc: GetValidateFunc(f.getStringParamValue(parameter, "validate")),
	}
	if !isDeep {
		return &currentParam, nil
	}
	// Extract variables
	switch currentParam.Type {
	case "string":
		currentParam.Default = CreateVariableString(f.getStringParamValue(parameter, "default"))
	case "bool", "optional":
		currentParam.Default = CreateVariableBool(f.getBoolParamValue(parameter, "default"))
	case "int":
		currentParam.Default = CreateVariableInt(f.getIntParamValue(parameter, "default"))
	}
	choices, err := GetChoices(f.getStringParamValue(parameter, "choices"))
	if err != nil {
		return nil, err
	}
	currentParam.Choices = choices
	if f.getStringParamValue(parameter, "default") == "" {
		currentParam.Default = Variable{}
	}
	return &currentParam, nil
}

func (f *Config) ExtractParams(isDeep bool) (*Parameters, []string, error) {
	// To return
	var parameters Parameters = make(Parameters)
	var includes []string = []string{}
	// Extract lists
	includesList, parametersList := f.extracKeys()
	includes = append(includes, includesList...)
	// Extract parameters
	for _, parameter := range parametersList {
		param, err := f.extractParam(parameter, isDeep)
		if err != nil {
			return nil, nil, err
		}
		parameters.AddAsParameter(parameter, param)
	}
	// Extract includes parameters
	for _, include := range includesList {
		// Create config instance for the include
		file, err := createFileInstanceFromPath(defaultConfPath + include)
		if err != nil {
			return nil, nil, err
		}
		conf, err := NewConfigFrom(file)
		if err != nil {
			return nil, nil, err
		}
		// Extract parameters
		fileParams, fileIncludes, err := conf.ExtractParams(isDeep)
		if err != nil {
			return nil, nil, err
		}
		// Merge data
		parameters.AddAsParameters(fileParams)
		includes = append(includes, fileIncludes...)
	}
	f.parameters = &parameters
	return &parameters, includes, nil
}

func (f *Config) ExtractValues() map[string]*Variable {
	// Extract parameters list
	parametersList := f.viperInstance.AllKeys()
	// Extract values
	var paramMap = make(map[string]*Variable)
	for _, parameter := range parametersList {
		str := f.viperInstance.GetString(parameter)
		boolean := f.viperInstance.GetBool(parameter)
		integer := f.viperInstance.GetInt(parameter)
		paramMap[parameter] = &Variable{
			String: &str,
			Bool:   &boolean,
			Int:    &integer,
		}
	}
	return paramMap
}

func (f *Config) getStringParamValue(name string, param string) string {
	return f.viperInstance.GetString(fmt.Sprintf("%s.%s", name, param))
}
func (f *Config) getBoolParamValue(name string, param string) bool {
	return f.viperInstance.GetBool(fmt.Sprintf("%s.%s", name, param))
}
func (f *Config) getIntParamValue(name string, param string) int {
	return f.viperInstance.GetInt(fmt.Sprintf("%s.%s", name, param))
}

func (f *Config) GetParams() *Parameters {
	return f.parameters
}

// Copy everything from the f.path to the destination path
func (f *Config) CopyToPath(dest string) error {
	return cp.Copy(f.file.Path, dest)
}

func (f *Config) SaveConfigTo(dest file) error {
	// Get flat map of parameters and convert to nested map
	var data = map[string]any{}
	for key, param := range *f.parameters {
		value, err := param.GetValue()
		if err != nil {
			fmt.Println("Error getting parameter value", key, err)
			fmt.Printf("Use %s=<your value> to set it\n", key)
			return err
		}
		data[key] = value
	}
	data = nestMap(data)
	// Process templates
	err := processTemplates(f.file.Path, dest.Path, data)
	if err != nil {
		log.Println("Error processing templates", err)
		return err
	}
	// Save parameters values to config file
	f.saveParamsTo(dest)
	// Get list of all included subconfigs
	_, includes, err := f.ExtractParams(false)
	if err != nil {
		return err
	}
	// Delete all included subconfig filess
	for _, include := range includes {
		err := os.Remove(filepath.Join(dest.Path, include))
		if err != nil {
			return err
		}
	}
	// Clean destination folder
	err = deleteEmptyFiles(dest.Path)
	if err != nil {
		log.Println("Error deleting empty files", err)
		return err
	}
	err = deleteEmptyFolders(dest.Path)
	if err != nil {
		log.Println("Error deleting empty folders", err)
		return err
	}
	return nil
}

// Save parameters values to config file
func (f *Config) saveParamsTo(dest file) error {
	//Clear the file
	err := os.Remove(dest.completePath())
	if err != nil {
		fmt.Println("Error removing config file", err)
		return err
	}
	//ReCreate the file
	file, err := os.Create(dest.completePath())
	if err != nil {
		fmt.Println("Error creating config file", err)
		return err
	}
	defer file.Close()

	// Instanciate config to dest file
	conf, err := NewConfigFrom(dest)
	if err != nil {
		fmt.Println("Error creating config instance", err)
		return err
	}
	conf.parameters = f.parameters
	//Get current config parameters values
	paramsValues := make(map[string]any)
	for key, param := range *conf.parameters {
		value, err := param.GetValue()
		if err != nil {
			fmt.Println("Error getting parameter value", key, err)
			log.Printf("Use %s=<your value> to set it", key)
			return err
		}
		paramsValues[key] = value
	}
	// Set the new values
	conf.viperInstance.Set("stamusconfig", f.file.Path)
	for key, value := range paramsValues {
		conf.viperInstance.Set(key, value)
	}
	// Write the new config file
	err = conf.viperInstance.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config file", err)
		return err
	}

	return nil
}
