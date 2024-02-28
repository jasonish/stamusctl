package compose

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func getInterfaceFormFS() ([]string, error) {
	cmd := exec.Command("ls", "/sys/class/net")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logging.Sugar.Infow("cannot fetch version.", "error", err)
		return nil, err
	}

	output := stdout.String()
	logging.Sugar.Debugw("detected interfaces.", "interfaces", output)
	splited := strings.Split(output, " ")
	return splited, nil
}

func getInterface(netInterface *string) {
	interfaces, err := RetrieveValideInterfaceFromDockerContainer()
	if err != nil {
		interfaces, _ = getInterfaceFormFS()
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

func getRestart(restart *string) {
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

func getElasticPath(elasticPath *string) {
	root, _ := GetDockerRootPath()
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

func getDataPath(elasticPath *string) {
	path, err := os.Getwd()
	prompt := promptui.Prompt{
		Label:   "container data volume path",
		Default: path,
	}

	result, err := prompt.Run()

	if err != nil {
		logging.Sugar.Error("Prompt for container data path.", err)
	}

	*elasticPath = result
	logging.Sugar.Debugw("selected container data path.", "result", result)
}

func getRegistry(registry *string) {
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

func Ask(cmd *cobra.Command, netInterface, restart, elasticPath, dataPath, registry, token *string) {
	if cmd.Flags().Changed("restart") == false {
		getInterface(netInterface)
	}

	if cmd.Flags().Changed("restart") == false {
		getRestart(restart)
	}

	if cmd.Flags().Changed("es-datapath") == false {
		getElasticPath(elasticPath)
	}

	if cmd.Flags().Changed("container-datapath") == false {
		getDataPath(dataPath)
	}

	if cmd.Flags().Changed("registry") == false {
		getRegistry(registry)
	}

	if cmd.Flags().Changed("token") == false {
		*token, _ = GenerateSciriusSecretToken()

	}
}
