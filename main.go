package main

import (
	"github.com/jon20/doku/cmd"
)

var (
	version = "0.1"
)

func main() {
	cmd.SetVersion(&cmd.Version{
		Version: version,
	})
	cmd.Execute()
}