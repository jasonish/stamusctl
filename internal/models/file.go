package models

import (
	// Common

	"strings"
	// External
)

// Used to setup viper instances
type file struct {
	Path string
	Name string
	Type string
}

// Used to get the file as properties from path
func createFileFromPath(path string) file {
	pathSplited := strings.Split(path, "/")
	nameSplited := strings.Split(pathSplited[len(pathSplited)-1], ".")
	return file{
		Path: strings.Join(pathSplited[:len(pathSplited)-1], "/"),
		Name: nameSplited[0],
		Type: nameSplited[1],
	}
}

// Used create a file from path and name
func CreateFile(path string, fileName string) file {
	nameSplited := strings.Split(fileName, ".")
	return file{
		Path: path,
		Name: nameSplited[0],
		Type: nameSplited[1],
	}
}

func (f *file) completePath() string {
	return f.Path + "/" + f.Name + "." + f.Type
}
