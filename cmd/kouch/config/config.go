package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/spf13/cobra"
)

const (
	dynamicContextName = "$dynamic$"
)

func init() {
	registry.Register([]string{}, func(_ *kouch.CmdContext) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "config",
			Short: "Modify kouchconfig files",
			Long: `Modify kouchconfig files using subcommands.

The loading order follows these rules:

  1. If the --` + kouch.FlagConfigFile + ` flag is set, that file is loaded.  The flag may only be set once and no merging takes place.
  2. Otherwise, ` + path.Join("${HOME}", kouch.HomeDir) + `/config is used and no merging takes place.`,
		}
		return cmd
	})

	registry.Register([]string{"config"}, func(cx *kouch.CmdContext) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "view",
			Short: "Display merged kouchconfig settings or a specified kouchconfig file",
			RunE:  viewConfig(cx),
		}
		return cmd
	})
}

func viewConfig(cx *kouch.CmdContext) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		return cx.Outputer.Output(os.Stdout, ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)))
	}
}
