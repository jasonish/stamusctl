package compose

import (
	"fmt"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/spf13/cobra"
)

func isValideRestartModeDocker(restart string) bool {
	if restart != "no" &&
		restart != "always" &&
		restart != "on-failure" &&
		restart != "unless-stopped" {
		return false
	}
	return true
}

func isValideRamSizeDocker(ram string) bool {
	return ram != "" && ram != "0" && ram != "0m" && ram != "0g	" && ram != "0k" && ram != "0t" && ram != "0p"
}

func ValidateInputFlag(params Parameters) error {
	if !isValideRestartModeDocker(params.RestartMode) {
		return fmt.Errorf("please provid a valid value for --restart. %s is not valid", params.RestartMode)
	}
	if !isValideRamSizeDocker(params.ElasticMemory) {
		return fmt.Errorf("please provide a valid value for --es-memory")
	}
	if !isValideRamSizeDocker(params.LogstashMemory) {
		return fmt.Errorf("please provide a valid value for --ls-memory")
	}
	return nil
}

func GenerateComposeFileFromCli(cmd *cobra.Command, params *Parameters, nonInteractive bool) string {
	if _, err := CheckVersions(); err != nil {
		logging.Sugar.Fatal(err.Error())
	}

	if !nonInteractive {
		Ask(cmd, params)
	}

	if params.InterfacesList == "" {
		logging.Sugar.Fatal("please provide a valid network interface.")
	}

	if params.DebugMode {
		params.NginxExec = "nginx"
	} else {
		params.NginxExec = "nginx-debug"
	}
	manifest := GenerateComposeFile(*params)

	return manifest
}
