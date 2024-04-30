package models

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// Parameters is a map of parameters, where the key is its place in the configuration file
type Parameters map[string]*Parameter

// Adds the parameters to the given parameters
func (p *Parameters) AddAsParameters(paramsList ...Parameters) *Parameters {
	for _, params := range paramsList {
		for key, param := range params {
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
func (p *Parameters) ValidateAll() string {
	for key, param := range *p {
		fmt.Println("Validating", key)
		if !param.Validate() {
			return key
		}
	}
	return ""
}

func (p *Parameters) AskAll() {
	for _, param := range *p {
		err := param.AskUser()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
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
			}
		}
	}
}
