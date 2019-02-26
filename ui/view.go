package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Image List View creates docker image list components
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
	line := FormatImageLine(v, "REPOSITORY", "TAG", "IMAGE ID", "SIZE", maxX)
	fmt.Fprintln(v, line)

	v, err = g.SetView("Image List", 0, 1, maxX-1, maxY/2)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}

	v.Title = v.Name()
	v.Frame = false
	v.Wrap = true
	v.Highlight = true
	go ShowContainerListWithAutoRefresh(g)
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	if _, err = SetCurrentViewOnTop(g, v.Name()); err != nil {
		return err
	}
	return nil

}

// Image List View creates docker container list components
func ContainerListView(g *gocui.Gui, maxX int, maxY int) error {

	v, err := g.SetView("Container List", 0, maxY/2, maxX-1, maxY-1)

	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Title = v.Name()
	v.Wrap = true
	v.Autoscroll = true

	return nil
}
