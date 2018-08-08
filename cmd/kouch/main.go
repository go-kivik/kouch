package main

import (
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/registry"

	"github.com/go-kivik/kouch/log"

	// The individual sub-commands
	_ "github.com/go-kivik/kouch/cmd/kouch/config"
	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
	_ "github.com/go-kivik/kouch/cmd/kouch/uuids"
)

func main() {
	cx := &kouch.CmdContext{
		Logger: log.New(),
	}
	registry.Root(cx).Execute()
}
