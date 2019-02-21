package utils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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
		if image.State != "running" {
			containers = append(containers, image)
		}
	}
	return &containers, nil
}

func (d *docker) GetImageList() (*[]types.ImageSummary, error) {
	var imageLists []types.ImageSummary
	images, err := d.Client.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		imageLists = append(imageLists, image)
	}
	return &imageLists, nil
}
