package main

import (
	"os"

	cli "github.com/arjunuxd/uxdguard/internal/cli"
)

var version string

func main() {
	if version != "" {
		cli.SetVersion(version)
	}
	os.Exit(cli.Execute(os.Args[1:]))
}
