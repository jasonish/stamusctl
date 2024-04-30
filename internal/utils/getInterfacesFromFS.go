package utils

import (
	"bytes"
	"os/exec"
	"strings"

	"stamus-ctl/internal/logging"
)

func GetInterfaceFormFS() ([]string, error) {
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
