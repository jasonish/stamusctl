package compose

import (
	"fmt"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/spf13/cobra"
)

func ValidateInputFlag(params Parameters) error {
	if params.RestartMode != "no" &&
		params.RestartMode != "always" &&
		params.RestartMode != "on-failure" &&
		params.RestartMode != "unless-stopped" {
		return fmt.Errorf("please provid a valid value for --restart. %s is not valid", params.RestartMode)
	}
	return nil
}

func GenerateComposeFileFromCli(cmd *cobra.Command, params Parameters, nonInteractive bool) string {
	if _, err := CheckVersions(); err != nil {
		logging.Sugar.Fatal(err.Error())
	}

	if !nonInteractive {
		Ask(cmd, &params)
	}

	if params.InterfacesList == "" {
		logging.Sugar.Fatal("please provide a valid network interface.")
	}

	if params.DebugMode {
		params.NginxExec = "nginx"
	} else {
		params.NginxExec = "nginx-debug"
	}
	manifest := GenerateComposeFile(params)

	return manifest
}
