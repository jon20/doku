package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Cont interface {
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
func (d *Docker) ContainerExecStart(containerID string, startCheck types.ExecStartCheck) error {
	err := d.Client.ContainerExecStart(context.Background(), containerID, startCheck)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docker) ContainerCreate(config *container.Config, hostConfig *container.HostConfig, networkConfig *network.NetworkingConfig, containerName string) error {
	create, err := d.Client.ContainerCreate(context.Background(), config, hostConfig, networkConfig, containerName)
	if err != nil {
		return err
	}
	fmt.Println(create)
	return nil
}

func (d *Docker) ContainerStart(containerID string, option types.ContainerStartOptions) error {
	err := d.Client.ContainerStart(context.Background(), containerID, option)
	if err != nil {
		return err
	}
	return nil
}
func (d *Docker) ContainerStop(containerID string, timeout *time.Duration) error {
	err := d.Client.ContainerStop(context.Background(), containerID, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docker) ContainerInspect(containerID string) (*types.ContainerJSON, error) {
	inspect, err := d.Client.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}
	return &inspect, nil
}
