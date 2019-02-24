package ui

import (
	"doku/utils"
	"fmt"

	"github.com/docker/docker/client"

	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

func ImageListView(g *gocui.Gui, maxX int, maxY int) error {

	v, err := g.SetView("Image", 0, 0, maxX-1, maxY/2)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Wrap = true
	v.Frame = true
	v.Title = v.Name()
	v.FgColor = gocui.AttrBold | gocui.ColorRed
	v.Clear()
	linea := pad.Right("REPOSITORYDA", maxX/3, " ") + pad.Right("TAG", maxX/5, " ") + pad.Right("IMAGE ID", maxX/5, " ") + pad.Right("SIZE", maxX/6, " ")
	fmt.Fprintln(v, linea)
	//line := FormatImageLine(v, "REPOSITORY", "TAG", "IMAGE ID", "SIZE", maxX)
	//line := pad.Right("REPOSITORY", 20, " ") + pad.Right("TAG", 10, " ") + pad.Right("IMAGE ID", 10, " ") + pad.Right("SIZE", 10, " ")
	//fmt.Fprintln(v, line)

	v, err = g.SetView("Image List", 0, 1, maxX-1, maxY/2)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}

	v.Title = v.Name()
	v.Frame = false
	v.Wrap = true
	v.Highlight = true
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	dockerHandler := utils.NewDockerClient(cli)
	go ShowContainerListWithAutoRefresh(g, &dockerHandler)
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	//if _, err = setCurrentViewOnTop(g, v.Name()); err != nil {
	//	return err
	//}
	return nil

}
