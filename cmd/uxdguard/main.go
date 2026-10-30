package main

import (
	"os"

	cli "github.com/arjunuxd/uxdguard/internal/cli"
)

var version string

func main() {
	if version != "" {
		cli.SetVersion(version)
	} else {
		cli.SetVersion("1.1.0")
	}
	os.Exit(cli.Execute(os.Args[1:]))
}
