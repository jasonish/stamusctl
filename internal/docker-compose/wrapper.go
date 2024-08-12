package compose

import (
	"os"
	"path/filepath"
	"stamus-ctl/internal/models"
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

// Constants
var ComposeFlags = models.ComposeFlags{
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

// Variables
var ComposeCmds map[string]*cobra.Command = make(map[string]*cobra.Command)

func GetComposeCmd(cmd string) *cobra.Command {
	WrappedCmd(ComposeFlags)
	return ComposeCmds[cmd]
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
		if composeFlags.Contains(command) && ComposeCmds[command] == nil {
			ComposeCmds[command] = c
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

// Modify the file flag to be hidden and add a folder flag
func modifyFileFlag(c *cobra.Command, command string) {
	c.Flags().Lookup("file").Hidden = true
	// Save the command
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
		fileFlag := GetComposeCmd(command).Flags().Lookup("file")
		fileFlag.Value.Set(flagValue)
		fileFlag.DefValue = flagValue
		// Run existing command
		return runE(cmd, args)
	}
}
