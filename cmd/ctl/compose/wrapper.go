package compose

import (
	// Common
	// External

	"os"
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

var composeFlags = ComposeFlags{
	"up": createComposeFlags(
		[]string{"file"},
		[]string{"detach"},
	),
	"down": createComposeFlags(
		[]string{"file"},
		[]string{"volumes", "remove-orphans"},
	),
}

// Commands
func wrappedCmd() []*cobra.Command {
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

	// Filter commands
	for _, c := range cmdDocker.Commands() {
		command := strings.Split(c.Use, " ")[0]
		if composeFlags.Contains(command) {
			flags := composeFlags[command].ExtractFlags(cmdDocker.Flags(), c.Flags())
			c.ResetFlags()
			c.Flags().AddFlagSet(flags)

			if c.Flags().Lookup("file") != nil {
				c.Flags().Lookup("file").Value.Set("tmp/docker-compose.yaml")
				c.Flags().Lookup("file").DefValue = "tmp/docker-compose.yaml"
			}

			cmds = append(cmds, c)
		}
	}

	// cmd.AddCommand(cmdDocker)

	return cmds
}

type Flags struct {
	Root []string
	Leaf []string
}

func (f *Flags) ExtractFlags(root *pflag.FlagSet, leaf *pflag.FlagSet) *pflag.FlagSet {
	var toReturn pflag.FlagSet
	for _, flag := range f.Root {
		if root.Lookup(flag) != nil {
			toReturn.AddFlag(root.Lookup(flag))
		}
	}
	for _, flag := range f.Leaf {
		if leaf.Lookup(flag) != nil {
			toReturn.AddFlag(leaf.Lookup(flag))
		}
	}
	return &toReturn
}

type ComposeFlags map[string]*Flags

func createComposeFlags(root []string, leaf []string) *Flags {
	return &Flags{
		Root: root,
		Leaf: leaf,
	}
}

func (c *ComposeFlags) Contains(command string) bool {
	_, ok := (*c)[command]
	return ok
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
