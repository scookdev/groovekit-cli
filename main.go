// Package main is the entry point for the GrooveKit CLI
package main

import "github.com/scookdev/groovekit-cli/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.Execute()
}
