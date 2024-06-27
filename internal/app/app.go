package app

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	Name   = ""
	Mode   = "prod"
	Folder = "/"
)

const (
	binaryNameEnv = "STAMUSCTL_NAME"

	CtlName = "stamusctl"
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
	Name = filepath.Base(os.Args[0])
	if val := os.Getenv(binaryNameEnv); val != "" {
		Name = val
	}

	if val := os.Getenv("BUILD_MODE"); val != "" {
		Mode = val
	}

	if val := os.Getenv("STAMUSCTL_FOLDER"); val != "" {
		Folder = val
	} else {
		Folder = os.Getenv("HOME") + "/.stamus"
	}
}
