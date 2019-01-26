package cmd

import (
	"fmt"
	"log"

	"github.com/docker/docker/client"
	"github.com/jon20/doku/utils"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
	"github.com/willf/pad"
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
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}

}
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		nextLine, err := v.Line(cy + 1)
		if err != nil {
			return nil
		}
		if nextLine == "" {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
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
	defer g.Close()

	g.SetManagerFunc(layout)

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

type hello struct {
	ID   string `tag:"ID"`
	Name string `tag:"Name"`
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("Image", 0, 0, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
		v.Wrap = true
		v.Frame = true
		v.Title = v.Name()
		v.FgColor = gocui.AttrBold | gocui.ColorRed
		line := pad.Right("name", 10, " ") + pad.Right("image", 10, " ")
		fmt.Fprintln(v, line)
	}
	if v, err := g.SetView("Image List", 0, 1, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = v.Name()
		v.Frame = false
		v.Wrap = true
		v.Highlight = true
		cli, _ := client.NewEnvClient()
		a := utils.NewDockerClient(cli)
		con, _ := a.GetImageList()
		fmt.Sprintf("%s:%s", "goo", "aaa")
		for _, item := range *con {
			fmt.Fprintln(v, hello{ID: item.RepoTags[0], Name: item.ParentID})
		}
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		if _, err = setCurrentViewOnTop(g, v.Name()); err != nil {
			return err
		}
	}

	if v, err := g.SetView("Container List", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = v.Name()
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

func ShowContainerList(v *gocui.View) {
	fmt.Println(v)
}
