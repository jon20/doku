package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func defaultCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Default command")
}