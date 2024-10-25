package handlers

import (
	// Core
	"fmt"
	"io"
	"os"
	"path/filepath"

	// Internal
	"stamus-ctl/internal/models"
)

type ReadPcapParams struct {
	Config   string
	PcapPath string
}

func PcapHandler(params ReadPcapParams) error {
	// Load existing config
	confFile, err := models.CreateFileInstance(params.Config, "values.yaml")
	if err != nil {
		return err
	}
	existingConfig, err := models.LoadConfigFrom(confFile, false)
	if err != nil {
		return err
	}

	// Get existing config value
	value, ok := existingConfig.GetParams().GetVariablesValues()["suricata.pcapreplay.hostpath"]
	if !ok {
		return fmt.Errorf("pcap_path parameter not found in existing config")
	}
	asString := value.AsString()

	// Copy pcap file to the host path
	err = CopyFile(params.PcapPath, asString)
	if err != nil {
		return err
	}
	fmt.Println("Pcap file is being read by the configuration")

	return nil
}

func CopyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(dst)
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Create destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
