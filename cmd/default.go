package cmd

import (
	"fmt"
	"log"

	"github.com/docker/docker/client"
	"github.com/jon20/doku/utils"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

type sample struct {
	art string
	ter string
}

var samples = []sample{
	{art: "aa", ter: "yes"},
}

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
	/*
		app := tview.NewApplication()
		list := tview.NewList().
			AddItem("List item 1", "Some explanatory text", 'a', nil).
			AddItem("List item 2", "Some explanatory text", 'b', nil).
			AddItem("List item 3", "Some explanatory text", 'c', nil).
			AddItem("List item 4", "Some explanatory text", 'd', nil).
			AddItem("Quit", "Press to exit", 'q', func() {
				app.Stop()
			})
		if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
			panic(err)
		}
	*/
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
