package utils

import (
	"os"
	"strings"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
)

func readFile(procFile string) string {
	data, err := os.ReadFile(procFile)
	if err != nil {
		logging.Sugar.Warnw("cannot read /proc/cpuinfo.", "error", err)
		return "sse4_2"
	}

	return string(data)

}

func findSupportInData(data string) bool {
	return strings.Contains(data, "sse4_2")
}

func GetSSESupport() bool {
	fileContent := readFile("/proc/cpuinfo")
	support := findSupportInData(fileContent)

	logging.Sugar.Debugw("support of sse4.2.", "support", support)

	return support
}
