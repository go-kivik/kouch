package config

import (
	"path"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{}, configCmd)
}

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Modify kouchconfig files",
		Long: `Modify kouchconfig files using subcommands.

The loading order follows these rules:

  1. If the --` + kouch.FlagConfigFile + ` flag is set, that file is loaded.  The flag may only be set once and no merging takes place.
  2. Otherwise, ` + path.Join("${HOME}", homeDir) + `/config is used and no merging takes place.`,
	}
}
