package models

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Parameter is equivalent to a flag
type Parameter struct {
	Name         string
	Shorthand    string
	Usage        string
	Type         string
	Variable     Variable
	Default      Variable
	Choices      []Variable
	ValidateFunc func(Variable) bool
}

// Variable struct is used to store values of different types
type Variable struct {
	String *string
	Bool   *bool
	Int    *int
}

func (p *Parameter) GetValue() any {
	switch p.Type {
	case "string":
		return *p.Variable.String
	case "bool", "optional":
		return *p.Variable.Bool
	case "int":
		return *p.Variable.Int
	default:
		return nil
	}
}

// Adds the parameter as a flag to the command
func (p *Parameter) AddAsFlag(cmd *cobra.Command, persistent bool) {
	switch p.Type {
	case "string":
		if p.Default.String == nil {
			p.Default = CreateVariableString("")
		}
		p.Variable = p.Default

		if p.Shorthand == "" {
			if persistent {
				cmd.PersistentFlags().StringVar(p.Variable.String, p.Name, *p.Default.String, p.Usage)
			} else {
				cmd.Flags().StringVar(p.Variable.String, p.Name, *p.Default.String, p.Usage)
			}
		} else {
			if persistent {
				cmd.PersistentFlags().StringVarP(p.Variable.String, p.Name, p.Shorthand, *p.Default.String, p.Usage)
			} else {
				cmd.Flags().StringVarP(p.Variable.String, p.Name, p.Shorthand, *p.Default.String, p.Usage)
			}
		}
	case "bool":
		if p.Default.Bool == nil {
			p.Default = CreateVariableBool(false)
		}
		p.Variable = p.Default
		if p.Shorthand == "" {
			if persistent {
				cmd.PersistentFlags().BoolVar(p.Variable.Bool, p.Name, *p.Default.Bool, p.Usage)
			} else {
				cmd.Flags().BoolVar(p.Variable.Bool, p.Name, *p.Default.Bool, p.Usage)
			}
		} else {
			if persistent {
				cmd.PersistentFlags().BoolVarP(p.Variable.Bool, p.Name, p.Shorthand, *p.Default.Bool, p.Usage)
			} else {
				cmd.Flags().BoolVarP(p.Variable.Bool, p.Name, p.Shorthand, *p.Default.Bool, p.Usage)
			}
		}
	case "int":
		if p.Default.Int == nil {
			p.Default = CreateVariableInt(0)
		}
		p.Variable = p.Default
		if p.Shorthand == "" {
			if persistent {
				cmd.PersistentFlags().IntVar(p.Variable.Int, p.Name, *p.Default.Int, p.Usage)
			} else {
				cmd.Flags().IntVar(p.Variable.Int, p.Name, *p.Default.Int, p.Usage)
			}
		} else {
			if persistent {
				cmd.PersistentFlags().IntVarP(p.Variable.Int, p.Name, p.Shorthand, *p.Default.Int, p.Usage)
			} else {
				cmd.Flags().IntVarP(p.Variable.Int, p.Name, p.Shorthand, *p.Default.Int, p.Usage)
			}
		}
	}
}

// Validates the variable with the given function
// If choices are provided, the variable must be in the list of choices
func (p *Parameter) Validate() bool {
	if p.ValidateFunc != nil {
		return p.ValidateFunc(p.Variable)
	}
	if p.Choices != nil && len(p.Choices) > 0 {
		switch p.Type {
		case "string":
			asStrings := []string{}
			for _, choice := range p.Choices {
				asStrings = append(asStrings, *choice.String)
			}
			return slices.Contains(asStrings, *p.Variable.String)
		case "int":
			asInts := []int{}
			for _, choice := range p.Choices {
				asInts = append(asInts, *choice.Int)
			}
			return slices.Contains(asInts, *p.Variable.Int)
		}
	}
	return true
}

// Validates the type of the variable
func (p *Parameter) ValidateType() bool {
	switch p.Type {
	case "string":
		return p.Variable.String != nil
	case "bool", "optional":
		return p.Variable.Bool != nil
	case "int":
		return p.Variable.Int != nil
	default:
		return false
	}
}

func (p *Parameter) AskUser() error {
	switch p.Type {
	case "string":
		if p.Choices != nil && len(p.Choices) > 0 {
			choices := []string{}
			for _, choice := range p.Choices {
				choices = append(choices, *choice.String)
			}
			result, err := selectPrompt(p, choices)
			if err != nil {
				return err
			}
			p.Variable = CreateVariableString(result)
		} else {
			defaultValue := *p.Default.String
			result, err := textPrompt(p, defaultValue)
			if err != nil {
				return err
			}
			p.Variable = CreateVariableString(result)
		}
	case "bool", "optional":
		result, err := selectPrompt(p, []string{"true", "false"})
		if err != nil {
			return err
		}
		p.Variable = CreateVariableBool(result == "true")
	case "int":
		defaultValue := strconv.Itoa(*p.Default.Int)
		result, err := textPrompt(p, defaultValue)
		if err != nil {
			return err
		}
		asInt, _ := strconv.Atoi(result)
		p.Variable = CreateVariableInt(asInt)
	}
	return nil
}

func CreateVariableString(value string) Variable {
	return Variable{String: &value}
}
func CreateVariableBool(value bool) Variable {
	return Variable{Bool: &value}
}
func CreateVariableInt(value int) Variable {
	return Variable{Int: &value}
}

func CreateParameterString(name string, shorthand string, variable Variable, defaultValue Variable, usage string) Parameter {
	return Parameter{
		Name:      name,
		Shorthand: shorthand,
		Type:      "string",
		Variable:  variable,
		Default:   defaultValue,
		Usage:     usage,
	}
}

func textPrompt(param *Parameter, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   param.Usage,
		Default: defaultValue,
		Validate: func(input string) error {
			switch param.Type {
			case "string":
				if param.ValidateFunc != nil && !param.ValidateFunc(CreateVariableString(input)) {
					return fmt.Errorf("invalid input")
				}
			case "int":
				asInt, err := strconv.Atoi(input)
				if err != nil {
					return err
				}
				if param.ValidateFunc != nil && !param.ValidateFunc(CreateVariableInt(asInt)) {
					return fmt.Errorf("invalid input")
				}
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("Prompt cancelled")
	}
	return result, nil
}

func selectPrompt(p *Parameter, choices []string) (string, error) {
	prompt := promptui.Select{
		Label: p.Usage,
		Items: choices,
		Templates: &promptui.SelectTemplates{
			Selected: fmt.Sprintf("%s {{ . | green }} %s ", promptui.IconGood, p.Usage),
		},
	}
	_, result, err := prompt.Run()

	if err != nil {
		return "", fmt.Errorf("Prompt cancelled")
	}
	return result, nil
}
