package compose

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/Masterminds/semver/v3"
)

func GetExecDockerVersion(executable string) (*semver.Version, error) {
	cmd := exec.Command(executable, "version", "--format", "{{.Server.Version}}")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logging.Sugar.Errorw("cannot fetch version.", "error", err, "exec", executable)
		return nil, fmt.Errorf("cannot %s fetch version", executable)
	}

	output := stdout.String()
	splited := strings.Split(output, " ")
	extracted := strings.Trim(splited[len(splited)-1], "\n")
	version, err := semver.NewVersion(extracted)
	if err != nil {
		logging.Sugar.Errorw("cannot parse version.", "error", err, "exec", extracted)
		return nil, fmt.Errorf("cannot parse %s version", executable)
	}

	logging.Sugar.Debugw("detected version.", "version", version, "executable", executable)

	return version, nil
}

func RetrieveValideInterfaceFromDockerContainer() ([]string, error) {

	images, _ := GetInstalledImages()
	var alreadyHasBusybox = true
	if images != nil {
		alreadyHasBusybox = slices.Contains(images, "busybox")
	}

	cmd := exec.Command("docker", "run", "--net=host", "--rm", "busybox", "ls", "/sys/class/net")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logging.Sugar.Infow("cannot fetch interfaces.", "error", err)
		return nil, err
	}

	output := stdout.String()
	logging.Sugar.Debugw("detected interfaces.", "interfaces", output)

	if !alreadyHasBusybox {
		logging.Sugar.Debug("busybox image was not previously installed, deleting.")
		DeleteDockerImage("busybox:latest")
	}

	return strings.Split(output, "\n"), nil
}

func GetInstalledImages() ([]string, error) {
	cmd := exec.Command("docker", "image", "ls", "--all")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logging.Sugar.Warnw("cannot fetch images.", "error", err)
		return nil, err
	}

	output := stdout.String()
	imagesFullList := strings.Split(output, "\n")
	var imagesName []string
	for i := 1; i != len(imagesFullList); i++ {
		imageData := imagesFullList[i]
		splited := strings.Split(imageData, " ")
		imageName := splited[0]

		imagesName = append(imagesName, imageName)
	}
	logging.Sugar.Debugw("detected images.", "images", imagesName)
	return imagesName, nil
}

func DeleteDockerImage(name string) (bool, error) {
	cmd := exec.Command("docker", "rmi", name)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logging.Sugar.Warnw("cannot delete image.", "error", err, "image", name)
		return false, err
	}

	return true, nil
}

func GetDockerRootPath() (string, error) {
	cmd := exec.Command("docker", "system", "info", "--format", "{{.DockerRootDir}}")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logging.Sugar.Errorw("cannot fetch docker root dir.", "error", err)
		return "", err
	}

	output := stdout.String()

	logging.Sugar.Debugw("detected docker root dir.", "dir", output)

	return strings.Trim(output, "\n"), nil
}
