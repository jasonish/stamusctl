package docker

import (
	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/docker/docker/api/types"
)

func DeleteDockerImageByName(name string) (bool, error) {

	id, err := GetImageIdFromName(name)

	if err != nil {
		logging.Sugar.Warnw("image id not found", "error", err)
		return false, err
	}

	if _, err = cli.ImageRemove(ctx, id, types.ImageRemoveOptions{}); err != nil {
		return false, err
	}

	return true, nil
}
