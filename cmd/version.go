package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Version struct {
	Version string
}

var version *Version

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "ss",
	Run: showVersion,
}


func init() {

}

func showVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", version.Version)
}