package models

import (
	// Common

	"strings"
	// External
)

// Used to setup viper instances
type fileProperties struct {
	Path string
	Name string
	Type string
}

// Ised to get the properties from path
func extractProperties(path string) fileProperties {
	pathSplited := strings.Split(path, "/")
	nameSplited := strings.Split(pathSplited[len(pathSplited)-1], ".")
	return fileProperties{
		Path: strings.Join(pathSplited[:len(pathSplited)-1], "/"),
		Name: nameSplited[0],
		Type: nameSplited[1],
	}
}
