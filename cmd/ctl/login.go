package ctl

import (
	// Common
	"context"
	"fmt"

	// External
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	// Custom
	parameters "stamus-ctl/internal/handlers"
	"stamus-ctl/internal/models"
	"stamus-ctl/internal/stamus"
)

func loginCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to a registry",
		RunE:  loginHandler,
	}
	// Flags
	parameters.Registry.AddAsFlag(cmd, false)
	parameters.Username.AddAsFlag(cmd, false)
	parameters.Password.AddAsFlag(cmd, false)
	return cmd
}

func loginHandler(cmd *cobra.Command, args []string) error {
	// Extract parameters
	registry, err := parameters.Registry.GetValue()
	if err != nil {
		return err
	}
	username, err := parameters.Username.GetValue()
	if err != nil {
		return err
	}
	password, err := parameters.Password.GetValue()
	if err != nil {
		return err
	}

	// Call handler
	params := models.RegistryInfo{
		Registry: registry.(string),
		Username: username.(string),
		Password: password.(string),
	}
	return LoginHandler(params)
}

func LoginHandler(registryInfo models.RegistryInfo) error {
	// Validate parameters from flags
	err := registryInfo.ValidateAllRegistry()
	if err != nil {
		return err
	}

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return err
	}

	// Log in to registry
	authConfig := registry.AuthConfig{
		ServerAddress: registryInfo.Registry,
		Username:      registryInfo.Username,
		Password:      registryInfo.Password,
	}
	_, err = cli.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		return err
	}

	// Save credentials
	stamus.SaveLogin(registryInfo)

	fmt.Println("Login successful")

	return nil
}
