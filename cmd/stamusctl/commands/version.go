package cmd

import (
	"fmt"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/app"
	"github.com/spf13/cobra"
)

func printVersion() {
	fmt.Printf("version: %s\narch: %s\ncommit: %s\n", app.Version, app.Arch, app.Commit)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {

}
