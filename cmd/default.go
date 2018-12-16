package cmd

import (
	"github.com/spf13/cobra"
	"github.com/marcusolsson/tui-go"
)

func defaultCmd(cmd *cobra.Command, args []string) {
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