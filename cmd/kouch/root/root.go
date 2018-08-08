package root

import (
	"github.com/spf13/cobra"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/config"
	"github.com/go-kivik/kouch/io"
	"github.com/go-kivik/kouch/log"
	"github.com/go-kivik/kouch/registry"

	// The individual sub-commands
	_ "github.com/go-kivik/kouch/cmd/kouch/config"
	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/uuids"
)

const version = "0.0.1"

// Run is the entry point, which executes the root command.
func Run() {
	l := log.New()

	cmd := rootCmd(l, version)
	if err := cmd.Execute(); err != nil {
		kouch.Exit(err)
	}
}

// Run initializes the root command, adds subordinate commands, then executes.
func rootCmd(l log.Logger, version string) *cobra.Command {
	cx := &kouch.CmdContext{
		Logger: l,
	}

	cmd := &cobra.Command{
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

	io.AddFlags(cmd)
	config.AddFlags(cmd)

	registry.AddSubcommands(cx, cmd)
	return cmd
}
