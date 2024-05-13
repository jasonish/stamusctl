package models

import (
	"fmt"
	"sort"
	"strconv"
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
		if !param.Validate() {
			return fmt.Errorf("Invalid value for %s", key)
		}
	}
	return nil
}

func (p *Parameters) AskAll() error {
	// Preprocess optional parameters
	err := p.ProcessOptionnalParams(true)
	if err != nil {
		return err
	}
	// Ask for all remaining parameters
	for _, param := range *p {
		if param.Type != "optional" {
			err := param.AskUser()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Parameters) SetToDefaults() error {
	// Preprocess optional parameters with default values
	err := p.ProcessOptionnalParams(false)
	if err != nil {
		return err
	}
	// Set all parameters to their default values
	for _, param := range *p {
		if param.Type != "optional" {
			param.SetToDefault()
		}
	}
	return nil
}

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
			// Remove all concerned parameters, except the optional one
			for paramKey, _ := range *p {
				if strings.HasPrefix(paramKey, optionalParam) && paramKey != optionalParam {
					delete(*p, paramKey)
				}
			}
			// Remove all concerned optional parameters
			remain := []string{}
			for _, key := range optionalParams {
				if !strings.HasPrefix(key, optionalParam) && key != optionalParam {
					remain = append(remain, key)
				}
			}
			optionalParams = remain
		} else {
			delete(*p, optionalParam)
		}
	}
	return nil
}

// Sets paramaters values to given values
func (p *Parameters) SetValues(values map[string]*Variable) {
	for key, value := range values {
		if (*p)[key] != nil {
			if (*p)[key].ValidateFunc != nil && !(*p)[key].ValidateFunc(*value) {
				fmt.Println("Invalid value for", key)
			} else {
				(*p)[key].Variable = *value
			}
		}
	}
}

func (p *Parameters) SetLooseValues(values map[string]string) {
	p.ProcessOptionnalParams(false)
	for key, value := range values {
		if (*p)[key] != nil {
			switch (*p)[key].Type {
			case "string":
				(*p)[key].Variable = CreateVariableString(value)
			case "bool":
				if value == "true" || value == "false" {
					(*p)[key].Variable = CreateVariableBool(value == "true")
				} else {
					fmt.Println("Invalid value for", key)
				}
			case "int":
				// Convert string to int
				asInt, err := strconv.Atoi(value)
				if err != nil {
					fmt.Println("Error converting string to int:", err)
				} else {
					asIntVar := CreateVariableInt(asInt)
					if (*p)[key].ValidateFunc != nil && !(*p)[key].ValidateFunc(asIntVar) {
						fmt.Println("Invalid value for", key)
					} else {
						(*p)[key].Variable = CreateVariableInt(asInt)
					}
				}
			case "optional":
				fmt.Println("Changing optional parameter", key, "is not supported")
				fmt.Println("Use `stamus-ctl compose init` to change optional blocks")
			}
		} else {
			fmt.Println("You cannot set value for", key)
		}
	}
}
