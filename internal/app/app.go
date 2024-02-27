package app

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	Name = ""
)

const (
	binaryNameEnv = "STAMUSCTL_NAME"
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
}
