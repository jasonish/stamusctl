package compose

import (
	// Common
	"archive/tar"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// External
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"

	// Custom
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/utils"
)

// Flags
var config = models.Parameter{
	Name:      "folder",
	Shorthand: "f",
	Type:      "string",
	Default:   models.CreateVariableString("tmp"),
	Usage:     "Folder configuration to update",
}
var registryFlag = models.Parameter{
	Name:  "registry",
	Type:  "string",
	Usage: "Registry to use",
}
var userFlag = models.Parameter{
	Name:  "user",
	Type:  "string",
	Usage: "Registry username",
}
var passwordFlag = models.Parameter{
	Name:  "pass",
	Type:  "string",
	Usage: "Registry password",
}
var versionFlag = models.Parameter{
	Name:      "version",
	Shorthand: "v",
	Type:      "string",
	Usage:     "Target version",
	Default:   models.CreateVariableString("latest"),
}

// Commands
func updateCmd() *cobra.Command {
	// Create cmd
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update compose configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateHandler(cmd, args)
		},
	}
	// Add flags
	registryFlag.AddAsFlag(cmd, false)
	userFlag.AddAsFlag(cmd, false)
	passwordFlag.AddAsFlag(cmd, false)
	versionFlag.AddAsFlag(cmd, false)
	config.AddAsFlag(cmd, false)
	return cmd
}

func updateHandler(cmd *cobra.Command, args []string) error {
	// Validate parameters from flags
	if *registryFlag.Variable.String == "" {
		err := registryFlag.AskUser()
		if err != nil {
			return err
		}
	}
	if *userFlag.Variable.String == "" {
		err := userFlag.AskUser()
		if err != nil {
			return err
		}
	}
	if *passwordFlag.Variable.String == "" {
		err := passwordFlag.AskUser()
		if err != nil {
			return err
		}
	}

	// Get values from flags
	var registryVal, usernameVal, passwordVal string
	registryVal = *registryFlag.Variable.String
	usernameVal = *userFlag.Variable.String
	passwordVal = *passwordFlag.Variable.String

	// Create docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Create auth config
	imageUrl := registryVal + "/selks:" + *versionFlag.Variable.String
	authConfig := registry.AuthConfig{
		Username: usernameVal,
		Password: passwordVal,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	// Pull image
	fmt.Println("Getting configuration")
	out, err := cli.ImagePull(ctx, imageUrl, types.ImagePullOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return err
	}
	defer out.Close()

	// Parse progress details
	type ImagePullResponse struct {
		Progress string `json:"progress"`
		Status   string `json:"status"`
	}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		var pullResp ImagePullResponse
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &pullResp); err != nil {
			fmt.Fprintf(os.Stderr, "\rError unmarshalling progress detail: %v", err)
			continue // Skip lines that can't be unmarshalled
		}

		if pullResp.Progress != "" {
			fmt.Printf("\r%s %s", pullResp.Status, pullResp.Progress)
		}
	}
	fmt.Printf("\rGot configuration                                                                                 ")
	fmt.Println()

	// Run container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageUrl,
		Cmd:   []string{"sleep 60"},
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	// Kill container
	defer func() {
		if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			fmt.Printf("Failed to remove container: %v\n", err)
		}
	}()

	// Extract conf from container
	srcPaths := []string{"/data", "/sbin"}                    // Source path inside the container
	destPath := filepath.Join(app.TemplatesFolder + "selks/") // Destination path on the host
	// Remove existing configuration
	if err := os.RemoveAll(filepath.Join(destPath, "latest")); err != nil {
		return err
	}
	// Copy files from container
	for _, srcPath := range srcPaths {
		if err := copyFromContainer(cli, ctx, resp.ID, srcPath, destPath); err != nil {
			return err
		}
	}
	// Move files to correct locations
	originPath := filepath.Join(destPath, "data/")
	latestPath := filepath.Join(destPath, "latest/")
	if err := os.Rename(originPath, latestPath); err != nil {
		return err
	}
	// Copy templates latest to templates version
	version, err := os.ReadFile(filepath.Join(latestPath, "version"))
	if err != nil {
		return err
	}
	err = cp.Copy(latestPath, filepath.Join(destPath, string(version)))
	if err != nil {
		return err
	}
	fmt.Println("Configuration extracted")

	// Execute update script
	prerunPath := filepath.Join(destPath, "sbin/pre-run")
	postrunPath := filepath.Join(destPath, "sbin/post-run")
	prerun := exec.Command(prerunPath)
	postrun := exec.Command(postrunPath)
	// Display output to terminal
	runOutput := new(strings.Builder)
	prerun.Stdout = runOutput
	prerun.Stderr = os.Stderr
	// Change execution rights
	os.Chmod(prerunPath, 0755)
	os.Chmod(postrunPath, 0755)
	// Run pre-run script
	if err := prerun.Run(); err != nil {
		return err
	}

	// Save output
	outputFile, err := os.Create(filepath.Join(*config.Variable.String, "values.yaml"))
	if err != nil {
		return err
	}
	defer outputFile.Close()
	if _, err := outputFile.WriteString(runOutput.String()); err != nil {
		return err
	}

	// Load existing config
	confFile, err := models.CreateFileInstance(*config.Variable.String, "values.yaml")
	if err != nil {
		return err
	}
	existingConfig, err := models.LoadConfigFrom(confFile, false)
	if err != nil {
		return err
	}

	// Create new config
	newConfFile, err := models.CreateFileInstance(latestPath, "config.yaml")
	if err != nil {
		return err
	}
	newConfig, err := models.NewConfigFrom(newConfFile)
	if err != nil {
		return err
	}
	_, _, err = newConfig.ExtractParams(true)
	if err != nil {
		return err
	}

	// Extract and set values from args and existing config
	paramsArgs := utils.ExtractArgs(args)
	newConfig.GetParams().SetValues(existingConfig.GetParams().GetVariablesValues())
	newConfig.GetParams().SetLooseValues(paramsArgs)
	newConfig.SetArbitrary(paramsArgs)
	newConfig.GetParams().ProcessOptionnalParams(false)

	// Ask for missing parameters
	err = newConfig.GetParams().AskMissing()
	if err != nil {
		return err
	}

	// Save the configuration
	err = newConfig.SaveConfigTo(confFile)
	if err != nil {
		return err
	}

	// Run post-run script
	postrunOutput := new(strings.Builder)
	postrun.Stdout = postrunOutput
	postrun.Stderr = os.Stderr
	// Run pre-run script
	if err := postrun.Run(); err != nil {
		return err
	}
	fmt.Println("")

	return nil
}

func copyFromContainer(cli *client.Client, ctx context.Context, containerID, srcPath, destPath string) error {
	reader, _, err := cli.CopyFromContainer(ctx, containerID, srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}
