package pkg

import "github.com/docker/docker/api/types"

// Init
type InitRequest struct {
	IsDefault  bool              `json:"default"`     // Use default settings, default is false
	Folder     string            `json:"folder"`      // Folder where to save configuration files, default is "tmp"
	Project    string            `json:"project"`     // Project name, default is "selks"
	Values     map[string]string `json:"values"`      // Values to set, key is the name of the value, value is the value
	Version    string            `json:"version"`     // Target version, default is latest
	ValuesPath string            `json:"values_path"` // Path to a values.yaml file
	FromFile   map[string]string `json:"from_file"`   // Values keys and paths to files containing the content used as value
}

// Config
type SetRequest struct {
	Reload     bool              `json:"reload"`      // Reload the configuration, don't keep arbitrary parameters
	Project    string            `json:"project"`     // Project name, default is "tmp"
	Values     map[string]string `json:"values"`      // Values to set, key is the name of the value, value is the value
	Apply      bool              `json:"apply"`       // Apply the new configuration, relaunch it, default is false
	ValuesPath string            `json:"values_path"` // Path to a values.yaml file
	FromFile   map[string]string `json:"from_file"`   // Values keys and paths to files containing the content used as value
}
type GetRequest struct {
	Project string   `json:"project"` // Project name, default is "tmp"
	Values  []string `json:"values"`  // Values to retrieve, default is all
	Content bool     `json:"content"` // Get content or values, default is false
}

// Update
type UpdateRequest struct {
	Version string            `json:"version"` // Version to update to, default is latest
	Project string            `json:"project"` // Project name, default is tmp
	Values  map[string]string `json:"values"`  // Values to set, key is the name of the value, value is the value
}

// Wrapper
type UpRequest struct {
	Project string `json:"project"` // Project name, default is tmp
}
type PsResponse struct {
	Containers []types.Container `json:"containers"`
}
