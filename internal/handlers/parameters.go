package handlers

import "stamus-ctl/internal/models"

// Init
var OutputPath = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
	Usage:     "Declare the folder where to save configuration files",
}
var IsDefaultParam = models.Parameter{
	Name:      "default",
	Shorthand: "d",
	Type:      "bool",
	Default:   models.CreateVariableBool(false),
	Usage:     "Set to default settings",
}

// Config
var ConfigPath = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
	Usage:     "Declare the folder where the configuration is saved",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
}
var Format = models.Parameter{
	Name:    "format",
	Usage:   "Format of the output (go template)",
	Type:    "string",
	Default: models.CreateVariableString("{{.}}"),
}

var Reload = models.Parameter{
	Name:    "reload",
	Usage:   "Reload the configuration, don't keep arbitrary parameters",
	Type:    "bool",
	Default: models.CreateVariableBool(false),
}

// Update
var Config = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
	Usage:     "Configuration to update",
}
var Registry = models.Parameter{
	Name:  "registry",
	Type:  "string",
	Usage: "Registry to use",
}
var Username = models.Parameter{
	Name:  "user",
	Type:  "string",
	Usage: "Registry username",
}
var Password = models.Parameter{
	Name:  "pass",
	Type:  "string",
	Usage: "Registry password",
}
var Version = models.Parameter{
	Name:      "version",
	Shorthand: "v",
	Type:      "string",
	Usage:     "Target version",
	Default:   models.CreateVariableString("latest"),
}
