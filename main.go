package main

import (
	"doku/cmd"
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
