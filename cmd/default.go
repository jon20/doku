package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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
		g.Close()
		log.Println(err)
		os.Exit(1)
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
		line := FormatImageLine(v, "REPOSITORY", "TAG", "IMAGE ID", "SIZE", maxX)
		//line := pad.Right("REPOSITORY", 20, " ") + pad.Right("TAG", 10, " ") + pad.Right("IMAGE ID", 10, " ") + pad.Right("SIZE", 10, " ")
		fmt.Fprintln(v, line)
	}
	// View: image
	if v, err := g.SetView("Image List", 0, 1, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = v.Name()
		v.Frame = false
		v.Wrap = true
		v.Highlight = true
		go ShowContainerListWithAutoRefresh(g)
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

func ShowContainerListWithAutoRefresh(g *gocui.Gui) {
	go ImagesRefresh(g)
	go ContainerListTitleResize(g)
	t := time.NewTicker(time.Duration(200 * time.Millisecond))
	for {
		select {
		case <-t.C:
			go ImagesRefresh(g)
			go ContainerListTitleResize(g)
		}
	}
}

func ContainerListTitleResize(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		maxX, _ := g.Size()
		v, err := g.View("Image")
		if err != nil {
			return err
		}
		v.Clear()
		line := FormatImageLine(v, "REPOSITORY", "TAG", "IMAGE ID", "SIZE", maxX)
		fmt.Fprintln(v, line)
		return nil
	})
}

func ImagesRefresh(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		maxX, _ := g.Size()
		v, err := g.View("Image List")
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}
		v.Clear()
		dockerHandler := utils.NewDockerClient(cli)
		images, err := dockerHandler.GetImageList()
		if err != nil {
			return err
		}
		for _, item := range *images {
			splitline := strings.Split(item.RepoTags[0], ":")
			size := strconv.FormatInt(item.Size, 10)
			line := FormatImageLine(v, splitline[0], splitline[0], splitline[0], size, maxX)
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
