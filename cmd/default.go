package cmd

import (
	"log"
	"os"

	"doku/ui"

	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
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
	g.SetManagerFunc(layout)

	SetKeyBindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		g.Close()
		log.Println(err)
		os.Exit(1)
	}

}

func SetKeyBindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ui.CursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ui.CursorUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.NextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("Container", 'r', gocui.ModNone, ui.ContainerStart); err != nil {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	ui.ImageListView(g, maxX, maxY)
	ui.ContainerListView(g, maxX, maxY)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
