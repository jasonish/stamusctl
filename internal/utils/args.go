package utils

import (
	"fmt"
	"strings"
)

func ExtractArgs(args []string) map[string]string {
	paramsArgs := make(map[string]string)
	for _, arg := range args {
		splited := strings.Split(arg, "=")
		if len(splited) != 2 {
			fmt.Println("Error: invalid argument", arg)
		} else {
			paramsArgs[splited[0]] = splited[1]
		}
	}
	return paramsArgs
}
