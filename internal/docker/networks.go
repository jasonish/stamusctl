package docker

import (
	"errors"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/docker/docker/api/types"
)

func GetNetworkIdByName(name string) (string, error) {
	networks, _ := cli.NetworkList(ctx, types.NetworkListOptions{})
	for _, network := range networks {

		if network.Name == name {
			logging.Sugar.Debugw("network found", "network.ID", network.ID, "network.Name", network.Name, "name", name)
			return network.ID, nil
		}
	}

	logging.Sugar.Debugw("network not found", "networks", networks, "name", name)
	return "", errors.New("network not found")
}
