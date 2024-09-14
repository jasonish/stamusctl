package utils

import (
	// Common

	"os"
	"strings"

	// Internal
	"stamus-ctl/internal/models"
)

// Check if the folder exists
func FolderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

var forbiddenChars = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", "$", "."}

func ValidatePath(path models.Variable) bool {
	if *path.String == "" {
		return false
	}
	for _, char := range forbiddenChars {
		if strings.Contains(*path.String, char) {
			return false
		}
	}
	return true
}
