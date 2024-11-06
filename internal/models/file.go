package models

import (
	// Common

	"fmt"
	"os"
	"strings"
	// External
)

// Used to setup viper instances
type File struct {
	Path string
	Name string
	Type string
}

// Used to get the file as properties from path
func CreateFileInstanceFromPath(path string) (File, error) {
	// Extract the file properties
	pathSplited := strings.Split(path, "/")
	if len(pathSplited) < 2 {
		pathSplited = []string{".", pathSplited[0]}
	}
	nameSplited := strings.Split(pathSplited[len(pathSplited)-1], ".")
	// Validate name
	if len(nameSplited) < 2 {
		return File{}, fmt.Errorf("path %s is not a valid file name", path)
	}
	// File
	file := File{
		Path: strings.Join(pathSplited[:len(pathSplited)-1], "/"),
		Name: strings.Join(nameSplited[:len(nameSplited)-1], "."),
		Type: nameSplited[len(nameSplited)-1],
	}
	// Validate all
	err := file.isValidPath()
	if err != nil {
		return File{}, err
	}
	// Return file instance
	return file, nil
}

// Used create a file from path and name
func CreateFileInstance(path string, fileName string) (File, error) {
	// Extract the file properties
	nameSplited := strings.Split(fileName, ".")
	if len(nameSplited) != 2 {
		return File{}, fmt.Errorf("path %s is not a valid file name", path)
	}
	// File
	file := File{
		Path: path,
		Name: nameSplited[0],
		Type: nameSplited[1],
	}

	// Validate all
	err := file.isValidPath()
	if err != nil {
		return File{}, err
	}

	// Return file instance
	return file, nil
}

func (f *File) completePath() string {
	return f.Path + "/" + f.Name + "." + f.Type
}

// Empirical function to check if a path is valid
func (f *File) isValidPath() error {
	// Check if file already exists
	if _, err := os.Stat(f.completePath()); err == nil {
		return nil
	}

	// Check parts
	if f.Path == "" {
		f.Path = "."
	}
	if f.Name == "" || f.Type == "" {
		return fmt.Errorf("type %s is not valid", f.Type)
	}

	// Return error if not possible
	return nil
}
