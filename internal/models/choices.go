package models

import (
	"io/ioutil"
	"log"
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
