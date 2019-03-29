package utils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Docker interface {
	GetActiveContainerList() (*[]types.Container, error)
	GetImageList() (*[]types.ImageSummary, error)
}
type Docker struct {
	Client *client.Client
}

func NewDockerClient(connect *client.Client) Docker {
	return Docker{connect}
}

func (d *Docker) GetActiveContainerList() (*[]types.Container, error) {
	var containers []types.Container
	images, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, nil
	}
	for _, image := range images {
		containers = append(containers, image)
	}
	return &containers, nil
}

func (d *Docker) GetImageList() (*[]types.ImageSummary, error) {
	var imageLists []types.ImageSummary
	images, err := d.Client.ImageList(context.Background(), types.ImageListOptions{All: false})
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		imageLists = append(imageLists, image)
	}
	return &imageLists, nil
}
