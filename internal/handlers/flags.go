package handlers

import (
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/stamus"
	"stamus-ctl/internal/utils"
)

// Init
var Config = models.Parameter{
	Name:         "config",
	Shorthand:    "c",
	Type:         "string",
	Default:      models.CreateVariableString(app.DefaultConfigName),
	Usage:        "Configuration path",
	ValidateFunc: utils.ValidatePath,
}
var IsDefaultParam = models.Parameter{
	Name:      "default",
	Shorthand: "d",
	Type:      "bool",
	Default:   models.CreateVariableBool(false),
	Usage:     "Set to default settings",
}
var Values = models.Parameter{
	Name:      "values",
	Shorthand: "v",
	Type:      "string",
	Default:   models.CreateVariableString(""),
	Usage:     "Values file to use",
}
var Template = models.Parameter{
	Name:      "template",
	Shorthand: "t",
	Type:      "string",
	Default:   models.CreateVariableString(""),
	Usage:     "Template folder to use",
	Hidden:    true,
}

// Config
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
var Apply = models.Parameter{
	Name:    "apply",
	Usage:   "Apply the new configuration",
	Type:    "bool",
	Default: models.CreateVariableBool(false),
}
var FromFile = models.Parameter{
	Name:    "fromFile",
	Usage:   "Uses the content of a file as parameter values",
	Type:    "string",
	Default: models.CreateVariableString(""),
}

// Update
var Version = models.Parameter{
	Name:    "version",
	Type:    "string",
	Usage:   "Target version",
	Default: models.CreateVariableString("latest"),
}

// Registry
var Registry = models.Parameter{
	Name:    "registry",
	Type:    "string",
	Usage:   "Registry to use",
	Default: models.CreateVariableString("docker.io/library/"),
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

func GetConfigFolderPath() (string, error) {
	var configPath string
	if app.IsCtl() {
		val, err := Config.GetValue()
		if err != nil {
			return "", err
		}
		configPath = val.(string)
	} else {
		conf, err := stamus.GetCurrent()
		if err != nil {
			return "", err
		}
		configPath = app.GetConfigsFolder(conf)
	}
	return configPath, nil
}
