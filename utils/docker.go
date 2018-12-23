package utils

import (
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types"
)

type docker struct {
	Client *client.Client
}

func NewDockerClient(connect *client.Client) docker {
	return docker{connect}
}

func (d *docker) GetActiveContainerList() (*[]types.Container, error) {
	var containers []types.Container
	images, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{Quiet: true})
	if err != nil {
		return nil, nil
	}
	for _, image := range images {
		containers = append(containers, image)
	}
	return &containers, nil
}

func (d *docker) GetInactiveContainerList() (*[]types.Container, error) {
	var containers []types.Container
	images, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, nil
	}
	for _, image := range images {
		if (image.State != "running") {
			containers = append(containers, image)
		}
	}
	return &containers, nil
}


