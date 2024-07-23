package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	root "stamus-ctl/cmd/daemon/run"
)

func newRequest(method string, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	// Create router
	router := root.SetupRouter(func(string) {})
	// Create a new request
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, newBody(body))
	if err != nil {
		return nil, err
	}
	// Serve the request
	router.ServeHTTP(w, req)

	return w, nil
}

// Generic input body for POST requests
func newBody(body interface{}) io.Reader {
	bodyJson, _ := json.Marshal(body)
	return bytes.NewReader(bodyJson)
}

// compareDirs compares the content of two directories
func compareDirs(t *testing.T, dir1, dir2 string) {
	folder1, err := getFolderContent(dir1)
	assert.NoError(t, err, fmt.Sprintf("failed to read directory %s with error %s", dir1, err))
	folder2, err := getFolderContent(dir2)
	assert.NoError(t, err, fmt.Sprintf("failed to read directory %s with error %s", dir2, err))

	err = compareFolderContent(folder1, folder2)
	assert.NoError(t, err, fmt.Sprintf("directories have different content with error %s", err))
	err = compareFolderContent(folder2, folder1)
	assert.NoError(t, err, fmt.Sprintf("directories have different content with error %s", err))
}

func getFolderContent(folder string) (map[string]string, error) {
	fileMap := make(map[string]string)
	err := readFolder(folder, folder, fileMap)
	if err != nil {
		return nil, err
	}
	return fileMap, nil
}

func readFolder(basePath, currentPath string, fileMap map[string]string) error {
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		relativePath := filepath.Join(currentPath, file.Name())
		if file.IsDir() {
			err := readFolder(basePath, relativePath, fileMap)
			if err != nil {
				return err
			}
		} else {
			content, err := ioutil.ReadFile(relativePath)
			if err != nil {
				return err
			}
			// Generate the key as a relative path from the base folder
			key, err := filepath.Rel(basePath, relativePath)
			if err != nil {
				return err
			}
			fileMap[key] = string(content)
		}
	}
	return nil
}

func compareFolderContent(folder1 map[string]string, folder2 map[string]string) error {
	if len(folder1) != len(folder2) {
		return fmt.Errorf("directories have different number of files")
	}
	for name, content1 := range folder1 {
		content2, ok := folder2[name]
		if !ok {
			return fmt.Errorf("file %s is missing in directory", name)
		}
		if content1 != content2 {
			return fmt.Errorf("file content mismatch for %s", name)
		}
	}
	return nil
}
