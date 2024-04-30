package docker

import (
	"bytes"
	"fmt"

	"stamus-ctl/internal/logging"

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
		fmt.Sprintf("Pulling %s. Please wait", name),
		fmt.Sprintf("Pulling %s done\n", name),
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
