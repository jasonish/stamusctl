package compose

import (
	// Common
	// External

	"log"
	"os"
	"strings"

	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/cmd/compatibility"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// Constants
const DefaultSelksPath = ".configs/selks/embedded"

// Commands
func ComposeCmd() *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "compose",
		Short: "Create container compose file",
	}

	// Custom commands
	cmd.AddCommand(initCmd())
	cmd.AddCommand(configCmd())

	// Docker OS stuff
	if plugin.RunningStandalone() {
		os.Args = append([]string{"docker"}, compatibility.Convert(os.Args[2:])...)
		log.Println(os.Args)
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

	// Filter commands
	toKeep := []string{"up", "down"}
	for _, c := range cmdDocker.Commands() {
		if contains(toKeep, strings.Split(c.Use, " ")[0]) {
			c.ResetFlags()
			c.Flags().AddFlagSet(cmdDocker.Flags())
			cmd.AddCommand(c)
		}
	}
	cmd.AddCommand(cmdDocker)

	return cmd
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func dockerCliPostInitialize(dockerCli command.Cli) {
	// HACK(milas): remove once docker/cli#4574 is merged; for now,
	// set it in a rather roundabout way by grabbing the underlying
	// concrete client and manually invoking an option on it
	_ = dockerCli.Apply(func(cli *command.DockerCli) error {
		if mobyClient, ok := cli.Client().(*client.Client); ok {
			_ = client.WithUserAgent("compose/" + "v2")(mobyClient)
		}
		return nil
	})
}

// func pluginMain() {
// 	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
// 		// TODO(milas): this cast is safe but we should not need to do this,
// 		// 	we should expose the concrete service type so that we do not need
// 		// 	to rely on the `api.Service` interface internally
// 		backend := compose.NewComposeService(dockerCli).(commands.Backend)
// 		cmd := commands.RootCommand(dockerCli, backend)
// 		originalPreRunE := cmd.PersistentPreRunE
// 		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
// 			// initialize the dockerCli instance
// 			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
// 				return err
// 			}
// 			// compose-specific initialization
// 			// dockerCliPostInitialize(dockerCli)

// 			if err := cmdtrace.Setup(cmd, dockerCli, os.Args[1:]); err != nil {
// 				logrus.Debugf("failed to enable tracing: %v", err)
// 			}

// 			if originalPreRunE != nil {
// 				return originalPreRunE(cmd, args)
// 			}
// 			return nil
// 		}

// 		cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
// 			return cli.StatusError{
// 				StatusCode: compose.CommandSyntaxFailure.ExitCode,
// 				Status:     err.Error(),
// 			}
// 		})
// 		return cmd
// 	},
// 		manager.Metadata{
// 			SchemaVersion: "0.1.0",
// 			Vendor:        "Docker Inc.",
// 			Version:       "internal.Version",
// 		})
// }
