package cmd

import (
	"github.com/spf13/cobra"
	"github.com/marcusolsson/tui-go"
	"github.com/jon20/doku/utils"
	"github.com/docker/docker/client"
	"fmt"
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
	box := tui.NewHBox(
		tui.NewLabel("tui-go"),
	)

	newui, err := tui.New(box)
	if err != nil {
		panic(err)
	}

	ui := setKetBinding(newui)
	if err := ui.Run(); err != nil {
		panic(err)
	}

}

func setKetBinding(ui tui.UI)(tui.UI) {
	ui.SetKeybinding("q", func() { ui.Quit() })

	return ui
}