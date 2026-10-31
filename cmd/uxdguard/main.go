package main

import (
	"os"

	cli "github.com/arjunuxd/uxdguard/internal/cli"
	_ "github.com/arjunuxd/uxdguard/internal/external"
)

var (
	version   string
	commit    string
	buildDate string
)

func main() {
	cli.SetVersion(version)
	cli.SetBuildMetadata(commit, buildDate)
	os.Exit(cli.Execute(os.Args[1:]))
}
