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

	_, err := docker.PullImageIfNotExisted("busybox")
	if err != nil {
		logging.SpinnerStop(s)
		return nil, err
	}

	output, _ := docker.RunContainer("busybox", []string{
		"ls",
		"/sys/class/net",
	}, nil, "host")

	interfaces := strings.Split(output, "\n")
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
