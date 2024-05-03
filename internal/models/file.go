package models

import (
	// Common

	"fmt"
	"io/ioutil"
	"os"
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
func createFileFromPath(path string) (file, error) {
	// Extract the file properties
	pathSplited := strings.Split(path, "/")
	nameSplited := strings.Split(pathSplited[len(pathSplited)-1], ".")
	// Validate all
	fmt.Println(path)
	if len(nameSplited) != 2 {
		return file{}, fmt.Errorf("path %s is not a valid file name", path)
	}
	err := isValidPath(path)
	if err != nil {
		return file{}, err
	}
	// Return file instance
	return file{
		Path: strings.Join(pathSplited[:len(pathSplited)-1], "/"),
		Name: nameSplited[0],
		Type: nameSplited[1],
	}, nil
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

// Empirical function to check if a path is valid
func isValidPath(path string) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(path, d, 0644); err == nil {
		os.Remove(path) // And delete it
		return nil
	}
	// Return error if not possible
	return fmt.Errorf("path %s is not valid", path)
}
