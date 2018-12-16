package cmd

import (
	"github.com/spf13/cobra"
	"github.com/marcusolsson/tui-go"
)

func defaultCmd(cmd *cobra.Command, args []string) {
	box := tui.NewHBox(
		tui.NewLabel("tui-go"),
	)

	ui, err := tui.New(box)
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func(){ ui.Quit() })
	if err := ui.Run(); err != nil {
		panic(err)
	}

}