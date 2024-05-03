package models

import (
	// Common

	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	// External
	"github.com/Masterminds/sprig/v3"
	cp "github.com/otiai10/copy"
	"github.com/spf13/viper"

	// Custom
	"stamus-ctl/internal/app"
)

type Config struct {
	path          string
	parameters    *Parameters
	viperInstance *viper.Viper
}

func NewConfigFrom(path string) (*Config, error) {
	conf := Config{
		path: path,
	}
	err := conf.instanciateViper(path)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func LoadConfigFrom(path string) (*Config, error) {
	// Load the config
	configured, err := NewConfigFrom(path)
	if err != nil {
		return nil, err
	}
	// Extract config data
	values := configured.ExtractValues()
	stamusConfPathPointer := values["stamusconfig"]
	stamusConfPath := *stamusConfPathPointer.String
	// Load origin config
	originConf, err := NewConfigFrom(stamusConfPath)
	if err != nil {
		return nil, err
	}
	originConf.ExtractParams()
	// Merge
	originConf.parameters.SetValues(values)
	return originConf, nil
}

func (f *Config) extracKeys() map[string]bool {
	// Extract parameters list
	parametersList := make(map[string]bool)
	for _, key := range f.viperInstance.AllKeys() {
		// Extract the parameter name
		parameterAsArray := strings.Split(key, ".")
		parameter := strings.Join(parameterAsArray[:len(parameterAsArray)-1], ".")
		parametersList[parameter] = true
	}
	return parametersList
}

func (f *Config) ExtractParams() *Parameters {
	// Extract parameters list
	parametersList := f.extracKeys()
	// Extract parameters
	var parameters Parameters = make(Parameters)
	for parameter, _ := range parametersList {
		// Extract the parameter
		currentParam := Parameter{
			Name:         f.getStringParamValue(parameter, "name"),
			Shorthand:    f.getStringParamValue(parameter, "shorthand"),
			Type:         f.getStringParamValue(parameter, "type"),
			Usage:        f.getStringParamValue(parameter, "usage"),
			ValidateFunc: GetValidateFunc(f.getStringParamValue(parameter, "validate")),
		}
		// Extract variables
		switch currentParam.Type {
		case "string":
			currentParam.Default = CreateVariableString(f.getStringParamValue(parameter, "default"))
			currentParam.Choices = GetChoices(f.getStringParamValue(parameter, "choices"))
		case "bool":
			currentParam.Default = CreateVariableBool(f.getBoolParamValue(parameter, "default"))
			currentParam.Choices = GetChoices(f.getStringParamValue(parameter, "choices"))
		case "int":
			currentParam.Default = CreateVariableInt(f.getIntParamValue(parameter, "default"))
			currentParam.Choices = GetChoices(f.getStringParamValue(parameter, "choices"))
		}
		// Add the parameter to the list
		parameters.AddAsParameter(parameter, &currentParam)
	}
	f.parameters = &parameters
	return &parameters
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

func (f *Config) instanciateViper(path string) error {
	// Extract the properties from the path
	properties := extractProperties(path + "/config.yaml")
	// Create a new viper instance
	f.viperInstance = viper.New()
	// General configuration
	f.viperInstance.SetEnvPrefix(app.Name)
	f.viperInstance.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	f.viperInstance.AutomaticEnv()
	// Specific configuration
	f.viperInstance.SetConfigName(properties.Name)
	f.viperInstance.SetConfigType(properties.Type)
	f.viperInstance.AddConfigPath(properties.Path)
	// Read the config file
	err := f.viperInstance.ReadInConfig()
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}
	return nil
}

func (f *Config) GetProjectParams() *Parameters {
	return f.parameters
}

// Copy everything from the f.path to the destination path
func (f *Config) CopyToPath(dest string) error {
	return cp.Copy(f.path, dest)
}

func (f *Config) SaveConfigTo(dest string) error {
	// Get flat map of parameters
	var data = map[string]any{}
	for key, param := range *f.parameters {
		data[key] = param.GetValue()
	}
	data = nestMap(data)
	// Process templates
	err := processTemplates(f.path, dest, data)
	if err != nil {
		return err
	}
	// Save parameters values to config file
	f.saveParamsTo(dest)

	return nil
}

// Nests a flat map into a nested map
func nestMap(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	// (Impressive) Reccursive stuff happening here
	for key, value := range input {
		parts := strings.Split(key, ".")

		if len(parts) > 1 {
			subMap := nestMap(map[string]interface{}{strings.Join(parts[1:], "."): value})

			if existingMap, ok := output[parts[0]].(map[string]interface{}); ok {
				for k, v := range subMap {
					existingMap[k] = v
				}
			} else {
				output[parts[0]] = subMap
			}
		} else {
			output[key] = value
		}
	}

	return output
}

func processTemplates(inputFolder string, outputFolder string, data map[string]interface{}) error {
	// Walk the source directory and process templates
	err := filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(inputFolder, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputFolder, rel)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		tmpl, err := template.New(filepath.Base(path)).Funcs(sprig.FuncMap()).ParseFiles(path)
		if err != nil {
			fmt.Println("Error parsing template", path, err)
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		return tmpl.Execute(destFile, data)

	})
	if err != nil {
		return err
	}
	log.Println("Templates processed", err)
	return nil
}

// Save parameters values to config file
func (f *Config) saveParamsTo(dest string) error {
	//Clear the file
	err := os.Remove(dest + "/config.yaml")
	if err != nil {
		fmt.Println("Error removing config file", err)
		return err
	}
	//ReCreate the file
	file, err := os.Create(dest + "/config.yaml")
	if err != nil {
		fmt.Println("Error creating config file", err)
		return err
	}
	defer file.Close()

	//Get current config parameters values
	f.instanciateViper(dest)
	paramsValues := make(map[string]any)
	for key, param := range *f.parameters {
		paramsValues[key] = param.GetValue()
	}
	// Set the new values
	f.viperInstance.Set("stamusconfig", f.path)
	for key, value := range paramsValues {
		f.viperInstance.Set(key, value)
	}
	// Write the new config file
	err = f.viperInstance.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config file", err)
		return err
	}

	return nil
}
