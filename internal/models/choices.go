package models

import (
	"strings"
	"unicode"

	"stamus-ctl/internal/docker"
	"stamus-ctl/internal/logging"
)

func GetChoices(name string) ([]Variable, error) {
	switch name {
	case "restart":
		return []Variable{
			CreateVariableString("no"),
			CreateVariableString("always"),
			CreateVariableString("on-failure"),
			CreateVariableString("unless-stopped"),
		}, nil
	case "nginx":
		return []Variable{
			CreateVariableString("nginx"),
			CreateVariableString("nginx-exec"),
		}, nil
	case "interfaces":
		return getInterfaces()
	default:
		return nil, nil
	}
}

func getInterfaces() ([]Variable, error) {
	s := logging.NewSpinner(
		"Identifying interfaces",
		"Did identify interfaces\n",
	)

	// alreadyHasBusybox, _ := docker.PullImageIfNotExisted("busybox")
	_, err := docker.PullImageIfNotExisted("busybox")
	if err != nil {
		logging.SpinnerStop(s)
		return nil, err
	}

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
	return interfacesVariables, nil
}
