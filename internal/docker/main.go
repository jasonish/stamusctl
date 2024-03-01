package docker

import (
	"context"

	"github.com/docker/docker/client"
)

var (
	ctx = context.Background()
	cli *client.Client
)

func init() {

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli = docker

	if err != nil {
		panic(err)
	}

}
