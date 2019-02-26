package ui

import (
	"doku/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/docker/docker/client"
	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

func ShowImageListWithAutoRefresh(g *gocui.Gui) {
	t := time.NewTicker(time.Duration(1 * time.Second))
	for {
		select {
		case <-t.C:
			ImagesRefresh(g)
			//go ContainerListTitleResize(g)
		}
	}
}

func ImagesRefresh(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		maxX, _ := g.Size()
		v, err := g.View("Image List")
		if err != nil {
			return err
		}
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}
		defer cli.Close()
		dockerHandler := utils.NewDockerClient(cli)
		images, err := dockerHandler.GetImageList()
		if err != nil {
			return err
		}
		v.Clear()
		for _, item := range *images {
			splitline := strings.Split(item.RepoTags[0], ":")
			size := strconv.FormatInt(item.Size, 10)
			fmt.Fprintln(v, splitline)
			line := FormatImageLine(v, splitline[0], splitline[0], splitline[0], size, maxX)
			fmt.Fprintln(v, line)
		}
		return nil
	})
}

func ShowContainerListWithAutoRefresh(g *gocui.Gui) {
	t := time.NewTicker(time.Duration(1 * time.Second))
	for {
		select {
		case <-t.C:
			ContainerListRefresh(g)
			//go ContainerListTitleResize(g)
		}
	}
}

func ContainerListRefresh(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		maxX, _ := g.Size()
		v, err := g.View("Container List")
		if err != nil {
			return err
		}
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}
		defer cli.Close()
		dockerHandler := utils.NewDockerClient(cli)
		containers, err := dockerHandler.GetActiveContainerList()
		if err != nil {
			return err
		}
		v.Clear()
		for _, item := range *containers {
			line := FormatImageLine(v, item.Names[0], item.Image, item.State, item.Status, maxX)
			fmt.Fprintln(v, line)
		}
		return nil
	})
}

func FormatImageLine(v *gocui.View, repository string, tag string, imageID string, size string, maxX int) string {
	// 30 30 10
	line := pad.Right(repository, maxX/3, " ") + pad.Right(tag, maxX/5, " ") + pad.Right(imageID, maxX/5, " ") + pad.Right(size, maxX/6, " ")
	return line
}

func LimitStringLine(words string, maxLength int) string {
	if utf8.RuneCountInString(words) > maxLength {
		return words[:maxLength]
	}
	return ""
}
