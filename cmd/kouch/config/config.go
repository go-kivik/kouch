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
	registry.Register([]string{"config"}, func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "view",
			Short: "Display merged kouchconfig settings or a specified kouchconfig file",
			RunE:  viewConfig(),
		}
		return cmd
	})
}

func viewConfig() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := kouch.GetContext(cmd)
		return kouch.Outputer(ctx).Output(os.Stdout, kouch.Conf(ctx).Dump())
	}
}
