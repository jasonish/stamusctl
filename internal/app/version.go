package app

import (
	"os"
	"path/filepath"
)

var (
	Arch          = ""
	Commit        = ""
	Version       = ""
	Name          = ""
	binaryNameEnv = "STAMUSCTL_NAME"
)

func init() {
	Name = filepath.Base(os.Args[0])
	if val := os.Getenv(binaryNameEnv); val != "" {
		Name = val
	}
}
