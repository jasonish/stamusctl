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
	"stamus-ctl/internal/app"
)

var defaultConfPath string

func init() {
	defaultConfPath = app.TemplatesFolder + "selks/embedded/"
}

// Config is a struct that represents a configuration file
// It contains the path to the file, the arbitrary values, the parameters values and the viper instnace to interact with the file
// It can be used to get or set values, validates them etc
type Config struct {
	file          File
	arbitrary     map[string]any
	parameters    *Parameters
	viperInstance *viper.Viper
}

// Create a new config instance from a file
func NewConfigFrom(file File) (*Config, error) {
	// Instanciate viper
	viperInstance, err := InstanciateViper(file)
	if err != nil {
		log.Println("Error instanciating viper", err)
		return nil, err
	}
	// Create the config
	conf := Config{
		file:          file,
		viperInstance: viperInstance,
		arbitrary:     make(map[string]any),
	}
	return &conf, nil
}

// Create a new config instance from a path, extract the values and return the instance
// Reload is used to not keep the arbitrary values
func LoadConfigFrom(path File, reload bool) (*Config, error) {
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
	// Set arbitrary
	if !reload {
		for key, value := range values {
			originConf.SetArbitrary(map[string]string{key: value.asString()})
		}
	}
	// Merge
	originConf.parameters.SetValues(values)
	return originConf, nil
}

func (f *Config) SetArbitrary(arbitrary map[string]string) {
	for key, value := range arbitrary {
		f.arbitrary[key] = asLooseTyped(value)
	}
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

// Extract parameters and includes from the config file
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
		file, err := CreateFileInstanceFromPath(defaultConfPath + include)
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

// Extract values from the config file
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

// Save the config to a folder
func (f *Config) SaveConfigTo(dest File) error {
	// Create config value map
	var data = map[string]any{}
	var configData = map[string]any{}
	for key, param := range *f.parameters {
		value, err := param.GetValue()
		if err != nil {
			fmt.Println("Error getting parameter value", key, err)
			fmt.Printf("Use %s=<your value> to set it\n", key)
			return err
		}
		configData[key] = value
	}
	// Merge with arbitrary config values and cerate a nested map
	for key, value := range f.arbitrary {
		data[key] = value
	}
	for key, value := range configData {
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

	// Clean destination folder
	err = f.Clean(dest)
	if err != nil {
		return err
	}

	return nil
}

// Cleans the config folder
func (f *Config) Clean(folder File) error {
	// Get list of all included subconfigs
	_, includes, err := f.ExtractParams(false)
	if err != nil {
		return err
	}
	// Delete all included subconfig files
	for _, include := range includes {
		err := os.Remove(filepath.Join(folder.Path, include))
		if err != nil {
			return err
		}
	}
	// Clean destination folder
	err = deleteEmptyFiles(folder.Path)
	if err != nil {
		log.Println("Error deleting empty files", err)
		return err
	}
	err = deleteEmptyFolders(folder.Path)
	if err != nil {
		log.Println("Error deleting empty folders", err)
		return err
	}
	return nil
}

// Save parameters values to config file
func (f *Config) saveParamsTo(dest File) error {
	//Clear the file
	os.Remove(dest.completePath())
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
	for key, value := range f.arbitrary {
		conf.viperInstance.Set(key, value)
	}
	// conf.viperInstance.Set("stamusconfig", f.file.Path)
	for key, value := range paramsValues {
		conf.viperInstance.Set(key, value)
	}
	// If latest, set stamusconfig value to version
	path := removeEmptyStrings(strings.Split(f.file.Path, "/"))
	if path[len(path)-1] == "latest" {
		// Get version
		version, err := os.ReadFile(filepath.Join(f.file.Path, "version"))
		if err != nil {
			log.Println("Error reading version file", err)
			return err
		}
		// Set stamusconfig value to version
		var versionPath []string = append([]string{}, path...)
		copy(versionPath, path)
		versionPath[len(versionPath)-1] = string(version)
		conf.viperInstance.Set("stamusconfig", "/"+filepath.Join(versionPath...))
	} else {
		conf.viperInstance.Set("stamusconfig", f.file.Path)
	}

	// Write the new config file
	err = conf.viperInstance.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config file", err)
		return err
	}

	return nil
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func (f *Config) DeleteFolder() error {
	return os.RemoveAll(f.file.Path)
}
