package utils

import (
	"fmt"
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

func (d *docker) GetImageList() error {
	images, err := d.Client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil
	}
	for _, image := range images {
		fmt.Println(image.ID)
		fmt.Println(image.RepoTags)
	}
	return nil
}

