package root

import (
	"github.com/spf13/cobra"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/config"
	"github.com/go-kivik/kouch/io"
	"github.com/go-kivik/kouch/registry"
)

const version = "0.0.1"

func init() {
	registry.RegisterRoot(rootCmd)
}

func rootCmd(cx *kouch.CmdContext) *cobra.Command {
	return &cobra.Command{
		Use:           "kouch",
		Short:         "kouch is a command-line tool for interacting with CouchDB",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			output, err := io.SelectOutput(cmd)
			if err != nil {
				return err
			}
			cx.Output = output
			outputer, err := io.SelectOutputProcessor(cmd)
			if err != nil {
				return err
			}
			cx.Outputer = outputer
			conf, err := config.ReadConfig(cmd)
			if err != nil {
				return err
			}
			cx.Conf = conf
			return nil
		},
	}
}
