package main

import (
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"

	// The root command
	_ "github.com/go-kivik/kouch/cmd/kouch/root"

	// Top-level sub-commands
	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/put"

	// The individual sub-commands
	_ "github.com/go-kivik/kouch/cmd/kouch/attachments"
	_ "github.com/go-kivik/kouch/cmd/kouch/config"
	_ "github.com/go-kivik/kouch/cmd/kouch/documents"
	_ "github.com/go-kivik/kouch/cmd/kouch/uuids"
)

func main() {
	Run()
}

// Run executes the root command.
func Run() {
	cmd := registry.Root()
	if err := cmd.Execute(); err != nil {
		kouch.Exit(err)
	}
}
