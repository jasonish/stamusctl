package models

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Parameters is a map of parameters, where the key is its place in the configuration file
type Parameters map[string]*Parameter

// Adds the parameters to the given parameters
func (p *Parameters) AddAsParameters(paramsList ...*Parameters) *Parameters {
	for _, params := range paramsList {
		for key, param := range *params {
			(*p)[key] = param
		}
	}
	return p
}

// Adds the parameter to the given parameters
func (p *Parameters) AddAsParameter(configName string, param *Parameter) {
	(*p)[configName] = param
}

// Adds the parameters as flags to the command
func (p *Parameters) AddAsFlags(cmd *cobra.Command, persistent bool) {
	for _, param := range *p {
		param.AddAsFlag(cmd, persistent)
	}
}

// Validates the parameters using their respective validation functions
// Returns the name of the parameter that failed validation or an empty string if all parameters are valid
func (p *Parameters) ValidateAll() error {
	for key, param := range *p {
		if !param.isValid() {
			return fmt.Errorf("Invalid value for %s", key)
		}
	}
	return nil
}

// Asks the user for all parameters
func (p *Parameters) AskAll() error {
	// Preprocess optional parameters
	err := p.ProcessOptionnalParams(true)
	if err != nil {
		return err
	}

	// Ask for all remaining parameters
	for _, key := range p.GetOrdered() {
		param := (*p)[key]
		if param.Type != "optional" {
			err := param.AskUser()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Asks the user for all not set parameters
func (p *Parameters) AskMissing() error {
	// Preprocess optional parameters
	err := p.ProcessOptionnalParams(false)
	if err != nil {
		return err
	}

	// Ask for all remaining parameters
	for _, key := range p.GetOrdered() {
		param := (*p)[key]
		if param.Type != "optional" && param.Variable.IsNil() {
			err := param.AskUser()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Set all parameters to their default values if they are not set
func (p *Parameters) SetToDefaults() error {
	// Preprocess optional parameters with default values
	err := p.ProcessOptionnalParams(false)
	if err != nil {
		return err
	}
	// Set all parameters to their default values
	for _, param := range *p {
		if param.Type != "optional" && param.Variable.IsNil() {
			param.SetToDefault()
		}
	}
	return nil
}

func (p *Parameters) GetValues(keys ...string) map[string]string {
	values := make(map[string]string)
	for key, param := range *p {
		// if keys are provided, only return values for keys that start with the provided keys
		if len(keys) > 0 {
			for _, k := range keys {
				if strings.HasPrefix(key, k) {
					values[key] = param.Variable.asString()
				}
			}
		} else {
			values[key] = param.Variable.asString()
		}
	}
	return values
}

func (p *Parameters) GetVariablesValues(keys ...string) map[string]*Variable {
	values := make(map[string]*Variable)
	for key, param := range *p {
		// if keys are provided, only return values for keys that start with the provided keys
		if len(keys) > 0 {
			for _, k := range keys {
				if strings.HasPrefix(key, k) {
					values[key] = &param.Variable
				}
			}
		} else {
			values[key] = &param.Variable
		}
	}
	return values
}

// Returns an ordered slices of the parameters keys
func (p *Parameters) GetOrdered() []string {
	keys := make([]string, 0, len(*p))
	for key := range *p {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Process optional parameters
// If interactive is true, ask the user for the optional parameters
func (p *Parameters) ProcessOptionnalParams(interactive bool) error {
	// Filter optional parameters
	optionalParams := []string{}
	for key, param := range *p {
		if param.Type == "optional" {
			optionalParams = append(optionalParams, key)
		}
	}
	// Sort by specificity
	sort.Slice(optionalParams, func(i, j int) bool {
		return len(strings.Split(optionalParams[i], ".")) < len(strings.Split(optionalParams[j], "."))
	})
	// Ask for optional parameters, filtering optional parameters and concerned parameters from instance
	for len(optionalParams) != 0 {
		// Get first element and remove it
		optionalParam := optionalParams[0]
		optionalParams = optionalParams[1:]
		// Get the optionnal parameter value
		param := (*p)[optionalParam]
		if interactive {
			err := param.AskUser()
			if err != nil {
				return err
			}
		} else {
			param.SetToDefault()
		}
		// Clean if false
		if !*param.Variable.Bool {
			p.cleanOptionatedParams(optionalParam)
			optionalParams = filterRemainingOptionalParams(optionalParams, optionalParam)
		} else {
			delete(*p, optionalParam)
		}
	}
	return nil
}

// Remove all concerned optional parameters
func (p *Parameters) cleanOptionatedParams(optionalParam string) {
	for paramKey := range *p {
		if strings.HasPrefix(paramKey, optionalParam) && paramKey != optionalParam {
			delete(*p, paramKey)
		}
	}
}

// Remove all concerned optional parameters
func filterRemainingOptionalParams(optionalParams []string, optionalParam string) []string {
	remain := []string{}
	for _, key := range optionalParams {
		if !strings.HasPrefix(key, optionalParam) && key != optionalParam {
			remain = append(remain, key)
		}
	}
	return remain
}

func (p *Parameters) MergeValues(toMerge *Parameters) *Parameters {
	for key, value := range *toMerge {
		if (*p)[key] != nil {
			(*p)[key].Variable = value.Variable
		}
	}
	return p
}

// Sets paramaters values to given values
func (p *Parameters) SetValues(values map[string]*Variable) {
	for key, value := range values {
		if (*p)[key] == nil {
			continue
		}
		if !(*p)[key].ValidateFunc(*value) {
			fmt.Println("Invalid value for", key)
		} else {
			(*p)[key].Variable = *value
		}
	}
}

func (p *Parameters) SetLooseValues(values map[string]string) error {
	for key, value := range values {
		if (*p)[key] != nil {
			(*p)[key].SetLooseValue(key, value)
		} else {
			fmt.Println("Invalid parameter", key)
		}
	}

	return nil
}

// Set values from a file (values.yaml)
func (p *Parameters) SetValuesFromFile(valuesPath string) error {
	if valuesPath != "" {
		file, err := CreateFileInstanceFromPath(valuesPath)
		if err != nil {
			return err
		}
		valuesConf, err := LoadConfigFrom(file, false)
		if err != nil {
			return err
		}
		p.MergeValues(valuesConf.GetParams())
	}
	return nil
}

// Set specific values from files content
func (p *Parameters) SetValuesFromFiles(fromFiles string) error {
	if fromFiles == "" {
		return nil
	}
	// For each fromFile
	args := strings.Split(fromFiles, " ")
	values := make(map[string]*Variable)
	for _, arg := range args {
		// Split argument
		split := strings.Split(arg, "=")
		if len(split) != 2 {
			return fmt.Errorf("Invalid argument: %s. Must be parameter.subparameter=./folder/file", arg)
		}
		// Get file content
		content, err := os.ReadFile(split[1])
		if err != nil {
			return err
		}
		// Set value of parameter
		temp := CreateVariableString(string(content))
		values[split[0]] = &temp
	}
	p.SetValues(values)
	return nil
}
