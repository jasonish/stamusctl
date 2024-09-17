package app

import (
	// Common
	"os"
	"path/filepath"
	"runtime"

	// External
	"github.com/adrg/xdg"
)

// Variables
var (
	Name              = ""
	Mode              = ModeStruct("prod")
	ConfigFolder      = "/"
	ConfigsFolder     = "/"
	TemplatesFolder   = "/"
	DefaultSelksPath  = "/"
	LatestSelksPath   = "/"
	DefaultConfigName = "tmp"
	StamusAppName     = ""
)

// Constants
const (
	binaryNameEnv = "STAMUSCTL_NAME"
	CtlName       = "stamusctl"
)

func CatchException() {
	if err := recover(); err != nil {
		switch err.(type) {
		case *runtime.Error:
			panic(err)
		default:
		}
	}
}

func init() {
	// Binary name
	Name = filepath.Base(os.Args[0])
	if val := os.Getenv(binaryNameEnv); val != "" {
		Name = val
	}

	// Mode
	if val := os.Getenv("BUILD_MODE"); val != "" {
		Mode.set(val)
	}
	if val := os.Getenv("STAMUS_APP_NAME"); val != "" {
		StamusAppName = val
	}

	// Folders
	if val := os.Getenv("STAMUS_CONFIG_FOLDER"); val != "" {
		ConfigFolder = val
	} else {
		ConfigFolder = xdg.ConfigHome + "/.stamus/"
	}
	if val := os.Getenv("STAMUS_TEMPLATES_FOLDER"); val != "" {
		TemplatesFolder = val
	} else {
		TemplatesFolder = xdg.UserDirs.Templates + "/.stamus/"
	}

	// Derived paths
	DefaultSelksPath = TemplatesFolder + "selks/embedded/"
	LatestSelksPath = TemplatesFolder + "selks/latest/"
	ConfigsFolder = ConfigFolder + "configs/"

}

func GetConfigsFolder(name string) string {
	return filepath.Join(ConfigsFolder, name)
}

func IsCtl() bool {
	return StamusAppName == CtlName
}
