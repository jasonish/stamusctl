package utils

import (
	// Common
	"log"
	"os"
	"stamus-ctl/internal/models"
	"strings"
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
	log.Println("Validating path", path)
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
