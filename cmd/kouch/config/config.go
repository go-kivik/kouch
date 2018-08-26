package config

import (
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/spf13/cobra"
)

const (
	dynamicContextName = "$dynamic$"
)

func init() {
	registry.Register([]string{"config"}, configCmd)
}

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Display merged kouchconfig settings or a specified kouchconfig file",
		RunE:  viewConfig(),
	}
}

func viewConfig() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := kouch.GetContext(cmd)
		return util.CopyAll(kouch.Output(ctx), kouch.Conf(ctx).Dump())
	}
}
