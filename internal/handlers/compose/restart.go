package handlers

import (
	// Common
	"context"
	"fmt"
	"sync"

	// Internal
	"stamus-ctl/internal/app"
	"stamus-ctl/pkg/mocker"

	// External
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func HandleConfigRestart(configPath string) error {
	if app.Mode.IsTest() {
		mocker.Mocked.Restart(configPath)
		return nil
	}
	return handleConfigRestart(configPath)
}

// HandleConfigRestart restarts the containers defined in the container composition file
func handleConfigRestart(configPath string) error {
	err := HandleDown(configPath, false, false)
	if err != nil {
		return err
	}
	return HandleUp(configPath)
}

func HandleContainersRestart(containers []string) error {
	if app.Mode.IsTest() {
		mocker.Mocked.RestartContainers(containers)
		return nil
	}
	return handleContainersRestart(containers)
}

// Given a list of container IDs, restarts them
func handleContainersRestart(containers []string) error {
	// Create docker client
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer apiClient.Close()
	// Sync
	wg := sync.WaitGroup{}
	wg.Add(len(containers))
	returned := make(chan error)
	defer close(returned)
	// Restart containers
	for _, containerID := range containers {
		go func(containerID string) {
			defer wg.Done()
			err := RestartContainer(containerID)
			if err != nil {
				returned <- err
			}
		}(containerID)
	}
	// Resync
	wg.Wait()
	if len(returned) != 0 {
		var toReturn error
		for err := range returned {
			toReturn = fmt.Errorf("%s\n%s", toReturn, err)
		}
		return toReturn
	}
	return nil
}

// RestartContainer restarts a container given its ID
func RestartContainer(containerID string) error {
	// Create docker client
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer apiClient.Close()
	// Restart container
	return apiClient.ContainerRestart(context.Background(), containerID, container.StopOptions{})
}
