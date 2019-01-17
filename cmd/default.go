package cmd

import (
	"fmt"
	"log"

	"github.com/docker/docker/client"
	"github.com/jon20/doku/utils"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

var (
	viewArr = []string{"v1", "v2", "v3", "v4"}
	active  = 0
)

func defaultCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		version, err := cmd.PersistentFlags().GetBool("version")
		if err == nil && version {
			showVersion(cmd, args)
			return
		}
	}
	cli, _ := client.NewEnvClient()
	a := utils.NewDockerClient(cli)
	con, _ := a.GetInactiveContainerList()
	fmt.Println(con)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(layout)

	setKeyBindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func setKeyBindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}

}
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}
func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("v2")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("Image", 0, 0, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = v.Name()
		v.Frame = true
		v.Wrap = true
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		fmt.Fprintln(v, "Item 4")
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		if _, err = setCurrentViewOnTop(g, "Image"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("v2", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v2"
		v.Wrap = true
		v.Autoscroll = true

	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}
