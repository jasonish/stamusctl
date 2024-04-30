package models

import (
	"strings"
	"unicode"

	"stamus-ctl/internal/docker"
	"stamus-ctl/internal/logging"
)

func GetChoices(name string) []Variable {
	switch name {
	case "restart":
		return []Variable{
			CreateVariableString("no"),
			CreateVariableString("always"),
			CreateVariableString("on-failure"),
			CreateVariableString("unless-stopped"),
		}
	case "nginx":
		return []Variable{
			CreateVariableString("nginx"),
			CreateVariableString("nginx-exec"),
		}
	case "interfaces":
		return getInterfaces()
	default:
		return nil
	}
}

func getInterfaces() []Variable {
	s := logging.NewSpinner(
		"Identifying interfaces",
		"Did identify interfaces\n",
	)

	// alreadyHasBusybox, _ := docker.PullImageIfNotExisted("busybox")
	docker.PullImageIfNotExisted("busybox")

	// log.Println("alreadyHasBusybox", alreadyHasBusybox)

	output, _ := docker.RunContainer("busybox", []string{
		"ls",
		"/sys/class/net",
	}, nil, "host")

	// log.Println("output", output)

	// if !alreadyHasBusybox {
	// 	// logging.Sugar.Debug("busybox image was not previously installed, deleting.")
	// 	docker.DeleteDockerImageByName("busybox")
	// }

	interfaces := strings.Split(output, "\n")
	// log.Println("interfaces", interfaces)
	interfaces = interfaces[:len(interfaces)-1]
	for i, in := range interfaces {
		in = strings.TrimFunc(in, unicode.IsControl)
		interfaces[i] = in
	}
	logging.Sugar.Debugw("detected interfaces.", "interfaces", interfaces)

	interfacesVariables := []Variable{}
	for _, in := range interfaces {
		if in != "" {
			interfacesVariables = append(interfacesVariables, CreateVariableString(in))
		}
	}

	logging.SpinnerStop(s)
	return interfacesVariables
}
