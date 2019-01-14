package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doku [Command]",
	Short: "Docker container managemant tool",
	Args:  cobra.MaximumNArgs(1),
	Run:   defaultCmd,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show doku Version")
}
