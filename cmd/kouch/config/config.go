package config

import (
	"os"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/spf13/cobra"
)

const (
	dynamicContextName = "$dynamic$"
)

func init() {
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
		return cx.Outputer.Output(os.Stdout, cx.Conf.Dump())
	}
}
