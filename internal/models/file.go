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
type File struct {
	Path string
	Name string
	Type string
}

// Used to get the file as properties from path
func CreateFileInstanceFromPath(path string) (File, error) {
	// Extract the file properties
	pathSplited := strings.Split(path, "/")
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
	// Attempt to create it
	var d []byte
	if err := os.MkdirAll(f.Path, 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(f.completePath(), d, 0644); err == nil {
		os.Remove(f.completePath()) // And delete it
		return nil
	}
	// Return error if not possible
	return fmt.Errorf("path %s is not valid", f.Path)
}
