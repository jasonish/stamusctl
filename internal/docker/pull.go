package docker

import (
	"bytes"
	"fmt"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/docker/docker/api/types"
)

func PullImageIfNotExisted(name string) (bool, error) {
	logger := logging.Sugar.With("name", name)
	alreadyHere, err := IsImageAlreadyInstalled(name)
	if err != nil {
		logger.Debugw("image failed to test", "error", err)
		return true, err
	}
	if alreadyHere {
		logger.Debugw("image found")
		return true, nil
	}

	logger.Debugw("image not found")

	s := logging.NewSpinner(
		fmt.Sprintf("pulling %s. please wait", name),
		fmt.Sprintf("pulling %s done", name),
	)
	reader, err := cli.ImagePull(ctx, "docker.io/library/"+name, types.ImagePullOptions{})

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	logging.SpinnerStop(s)
	if err != nil {
		logger.Debugw("image failed to dl", "error", err)
		return false, err
	}

	logger.Debugw("image dl", "error", err)
	return false, nil
}
