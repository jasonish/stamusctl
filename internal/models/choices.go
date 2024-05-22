package models

import (
	"io/ioutil"
	"log"
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/docker"
	"stamus-ctl/internal/logging"
	"strings"
	"unicode"
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
	if app.Mode == "prod" {
		return getInterfacesBusybox()
	} else {
		return getInterfacesHost()
	}
}

func getInterfacesHost() ([]Variable, error) {
	// Define the directory where network interfaces are listed
	netDir := "/sys/class/net"

	// Read the directory contents
	files, err := ioutil.ReadDir(netDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", netDir, err)
	}

	// Loop through the files and print the names
	interfaces := []Variable{}
	for _, file := range files {
		interfaces = append(interfaces, CreateVariableString(file.Name()))
	}
	return interfaces, nil
}

func getInterfacesBusybox() ([]Variable, error) {
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
