package test

import (
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
