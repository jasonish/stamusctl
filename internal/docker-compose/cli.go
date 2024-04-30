package compose

import (
	"errors"
	"os"
	"strings"

	"stamus-ctl/internal/logging"
	"stamus-ctl/internal/utils"

	"github.com/manifoldco/promptui"
)

func getInterfaceCli(netInterface *string) {
	interfaces, err := RetrieveValideInterfacesFromDockerContainer()
	if err != nil {
		interfaces, _ = utils.GetInterfaceFormFS()
	}

	interfacesSelect := promptui.Select{
		Label: "select network interface",
		Items: interfaces,
	}
	_, result, err := interfacesSelect.Run()

	if err != nil {
		logging.Sugar.Error("Prompt for network interface failed.", err)
	}

	*netInterface = result
	logging.Sugar.Debugw("selected interface.", "interface", netInterface)
}

func getRestartCli(restart *string) {
	prompt := promptui.Prompt{
		Label:     "Do you want the containers to restart automatically on startup",
		IsConfirm: true,
		Default:   "y",
	}
	validate := func(s string) error {
		if len(s) == 1 && strings.Contains("YyNn", s) || prompt.Default != "" && len(s) == 0 {
			return nil
		}
		return errors.New("invalid input")
	}
	prompt.Validate = validate

	_, err := prompt.Run()
	confirmed := !errors.Is(err, promptui.ErrAbort)
	if err != nil && confirmed {
		logging.Sugar.Error("Prompt for restart failed.", err)
		return
	}

	if confirmed {
		*restart = "unless-stopped"
	} else {
		*restart = "never"
	}
	logging.Sugar.Debugw("selected restart.", "restart", *restart)
}

func getElasticPathCli(elasticPath *string) {
	root := utils.IgnoreError(os.Getwd()) + "/elastic-data"
	prompt := promptui.Prompt{
		Label:   "elasticsearch database volume path",
		Default: root,
	}

	result, err := prompt.Run()

	if err != nil {
		logging.Sugar.Error("Prompt for elastic data path.", err)
	}

	*elasticPath = result
	logging.Sugar.Debugw("selected elastic data path.", "result", result)
}

func getRegistryCli(registry *string) {
	prompt := promptui.Prompt{
		Label:   "Image registry",
		Default: "ghcr.io/stamusnetworks",
	}

	result, err := prompt.Run()

	if err != nil {
		logging.Sugar.Error("Prompt for container data path.", err)
	}

	*registry = result
	logging.Sugar.Debugw("selected container data path.", "result", result)
}
