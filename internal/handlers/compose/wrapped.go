package handlers

import (
	"log"
	"os"
	"path/filepath"
	"stamus-ctl/internal/models"
	"strconv"
	"strings"

	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/cmd/compatibility"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/compose"
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
		[]string{"services", "quiet"},
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
				c.Flags().Lookup("file").Hidden = true
				// Save the command
				composeCmds[command] = c
				currentRunE := c.RunE
				// Modify cmd function
				c.RunE = func(cmd *cobra.Command, args []string) error {
					log.Println("cmd.Flags()", cmd.Flags())
					flagValue := filepath.Join(cmd.Flags().Lookup("folder").Value.String(), "/docker-compose.yaml")

					fileFlag := getComposeCmd(command).Flags().Lookup("file")
					fileFlag.Value.Set(flagValue)
					fileFlag.DefValue = flagValue

					return currentRunE(cmd, args)
				}
				// Add custom folder flag
				folderFlag := *pflag.NewFlagSet("folder", pflag.ContinueOnError)
				folderFlag.String("folder", "tmp", "Folder where the config is located")
				c.Flags().AddFlagSet(&folderFlag)
			}

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
