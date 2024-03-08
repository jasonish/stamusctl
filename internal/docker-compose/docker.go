package compose

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/docker"
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

func RetrieveValideInterfacesFromDockerContainer() ([]string, error) {

	alreadyHasBusybox, _ := docker.PullImageIfNotExisted("busybox")

	output, _ := docker.RunContainer("busybox", []string{
		"ls",
		"/sys/class/net",
	}, nil, "host")

	if !alreadyHasBusybox {
		logging.Sugar.Debug("busybox image was not previously installed, deleting.")
		docker.DeleteDockerImageByName("busybox")
	}

	interfaces := strings.Split(output, "\n")
	interfaces = interfaces[:len(interfaces)-1]
	for i, in := range interfaces {
		in = strings.TrimFunc(in, unicode.IsControl)
		interfaces[i] = in
	}
	logging.Sugar.Debugw("detected interfaces.", "interfaces", interfaces)

	return interfaces, nil
}

func GenerateSSLWithDocker(sslPath string) error {
	logging.Sugar.Debugw("Generating ssl cert.", "path", sslPath)

	err := os.MkdirAll(sslPath, os.ModePerm)
	if err != nil {
		logging.Sugar.Errorw("cannot create cert containing folder.", "error", err)
	}

	alreadyHasNginx, err := docker.PullImageIfNotExisted("nginx")
	if err != nil {
		logging.Sugar.Warnw("nginx pull failed", "error", err)
		return err
	}

	_, err = docker.RunContainer("nginx", []string{
		"openssl",
		"req", "-new", "-nodes", "-x509",
		"-subj", "/C=FR/ST=IDF/L=Paris/O=Stamus/CN=SELKS",
		"-days", "3650",
		"-keyout", "/etc/nginx/ssl/scirius.key",
		"-out", "/etc/nginx/ssl/scirius.crt",
		"-extensions", "v3_ca",
	}, []string{
		sslPath + ":/etc/nginx/ssl",
	}, "")

	if err != nil {
		logging.Sugar.Infow("cannot generate cert.", "error", err)
		return err
	}

	if !alreadyHasNginx {
		logging.Sugar.Debug("nginx image was not previously installed, deleting.")
		docker.DeleteDockerImageByName("nginx")
	}

	logging.Sugar.Debug("cert created.", "path", sslPath)
	return nil
}
