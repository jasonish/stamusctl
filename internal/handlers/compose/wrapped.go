package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"stamus-ctl/internal/models"
	"stamus-ctl/pkg"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/cmd/compatibility"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Variables
var composeFlags = models.ComposeFlags{
	"up": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"detach"},
	),
	"down": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"volumes", "remove-orphans"},
	),
	"ps": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"services", "quiet", "format"},
	),
	"logs": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"timestamps", "tail", "since", "until"},
	),
}
var composeCmds map[string]*cobra.Command = make(map[string]*cobra.Command)

func getComposeCmd(cmd string) *cobra.Command {
	return composeCmds[cmd]
}

// Handlers
func WrappedCmd(composeFlags models.ComposeFlags) ([]*cobra.Command, map[string]*cobra.Command) {
	// Docker stuff
	if plugin.RunningStandalone() && len(os.Args) > 2 && os.Args[1] == "compose" {
		os.Args = append([]string{"docker"}, compatibility.Convert(os.Args[2:])...)
	}
	// Create docker client
	cliOptions := func(cli *command.DockerCli) error {
		op := &flags.ClientOptions{}
		cli.Initialize(op)
		return nil
	}
	dockerCli, err := command.NewDockerCli(cliOptions)
	if err != nil {
		panic(err)
	}
	// Create docker command
	backend := compose.NewComposeService(dockerCli).(commands.Backend)
	cmdDocker := commands.RootCommand(dockerCli, backend)

	// Stuff to return
	cmds := []*cobra.Command{}
	mappedCmds := make(map[string]*cobra.Command)

	// Filter commands
	for _, c := range cmdDocker.Commands() {
		command := strings.Split(c.Use, " ")[0]
		if composeFlags.Contains(command) {
			// Filter flags
			flags := composeFlags[command].ExtractFlags(cmdDocker.Flags(), c.Flags())
			c.ResetFlags()
			c.Flags().AddFlagSet(flags)
			// Modify file flag
			if c.Flags().Lookup("file") != nil {
				modifyFileFlag(c, command)
			}
			// Save command
			cmds = append(cmds, c)
			mappedCmds[command] = c
		}
	}
	return cmds, mappedCmds
}

func HandleUp(configPath string) error {
	// Get command
	WrappedCmd(composeFlags)
	command := getComposeCmd("up")
	// Set flags
	command.Flags().Lookup("folder").DefValue = configPath
	command.Flags().Lookup("folder").Value.Set(configPath)
	command.Flags().Lookup("detach").DefValue = "true"
	command.Flags().Lookup("detach").Value.Set("true")
	// Create root command
	var cmd *cobra.Command = &cobra.Command{Use: "compose"}
	cmd.SetArgs([]string{"up"})
	cmd.AddCommand(command)
	// Run command
	return cmd.Execute()
}

func HandleDown(configPath string, removeOrphans bool, volumes bool) error {
	// Get command
	WrappedCmd(composeFlags)
	command := getComposeCmd("down")
	// Set flags
	command.Flags().Lookup("folder").DefValue = configPath
	command.Flags().Lookup("folder").Value.Set(configPath)
	command.Flags().Lookup("remove-orphans").DefValue = strconv.FormatBool(removeOrphans)
	command.Flags().Lookup("remove-orphans").Value.Set(strconv.FormatBool(removeOrphans))
	command.Flags().Lookup("volumes").DefValue = strconv.FormatBool(volumes)
	command.Flags().Lookup("volumes").Value.Set(strconv.FormatBool(volumes))
	// Create root command
	var cmd *cobra.Command = &cobra.Command{Use: "compose"}
	cmd.SetArgs([]string{"up"})
	cmd.AddCommand(command)
	// Run command
	return cmd.Execute()
}

func HandlePs() ([]types.Container, error) {
	// Create docker client
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer apiClient.Close()
	// Get containers
	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func HandleConfigRestart(configPath string) error {
	err := HandleDown(configPath, false, false)
	if err != nil {
		return err
	}
	return HandleUp(configPath)
}

func HandleContainersRestart(containers []string) error {
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

func HandleLogs(logParams pkg.LogsRequest) (pkg.LogsResponse, error) {
	// Get containers
	containers, err := HandlePs()
	if err != nil {
		return pkg.LogsResponse{}, err
	}
	// Filter containers
	if logParams.Containers != nil && len(logParams.Containers) > 0 {
		filteredContainers := []types.Container{}
		for _, container := range containers {
			for _, id := range logParams.Containers {
				if container.ID == id {
					filteredContainers = append(filteredContainers, container)
					break
				}
			}
		}
		containers = filteredContainers
	}
	// Create docker client
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return pkg.LogsResponse{}, err
	}
	defer apiClient.Close()
	// Get logs
	response := pkg.LogsResponse{}
	for _, current := range containers {
		reader, err := apiClient.ContainerLogs(context.Background(), current.ID, container.LogsOptions{
			Since:      logParams.Since,
			Until:      logParams.Until,
			Tail:       logParams.Tail,
			Timestamps: logParams.Timestamps,
			ShowStdout: true,
			ShowStderr: true,
			Details:    true,
		})
		if err != nil {
			return pkg.LogsResponse{}, err
		}
		defer reader.Close()
		// Read logs
		buffer, err := io.ReadAll(reader)
		if err != nil {
			return pkg.LogsResponse{}, err
		}
		// Clean logs
		rawLogs := string(buffer)
		rawLogs = strings.ReplaceAll(rawLogs, "\u0000", "")
		rawLogs = strings.ReplaceAll(rawLogs, "\u0001", "")
		rawLogs = strings.ReplaceAll(rawLogs, "\u0002", "")
		logs := strings.Split(rawLogs, "\n")
		// Create response
		response.Containers = append(response.Containers, pkg.ContainerLogs{
			Container: current,
			Logs:      logs,
		})
	}

	return response, nil
}

// Modify the file flag to be hidden and add a folder flag
func modifyFileFlag(c *cobra.Command, command string) {
	c.Flags().Lookup("file").Hidden = true
	// Save the command
	composeCmds[command] = c
	currentRunE := c.RunE
	// Modify cmd function
	c.RunE = makeCustomRunner(currentRunE, command)
	// Add custom folder flag
	folderFlag := *pflag.NewFlagSet("folder", pflag.ContinueOnError)
	folderFlag.String("folder", "tmp", "Folder where the config is located")
	c.Flags().AddFlagSet(&folderFlag)
}

// Return a custom runner for the command, that sets the file flag to the folder flag
func makeCustomRunner(
	runE func(cmd *cobra.Command, args []string) error,
	command string,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Get folder flag value
		flagValue := filepath.Join(cmd.Flags().Lookup("folder").Value.String(), "/docker-compose.yaml")
		// Set file flag
		fileFlag := getComposeCmd(command).Flags().Lookup("file")
		fileFlag.Value.Set(flagValue)
		fileFlag.DefValue = flagValue
		// Run existing command
		return runE(cmd, args)
	}
}
