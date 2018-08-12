package root

import (
	"github.com/spf13/cobra"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/config"
	"github.com/go-kivik/kouch/io"

	// The individual sub-commands
	_ "github.com/go-kivik/kouch/cmd/kouch/attachments"
	_ "github.com/go-kivik/kouch/cmd/kouch/config"
	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/uuids"
)

const version = "0.0.1"

// global config flags
const (
	flagVerbose = "verbose"
)

// Run is the entry point, which executes the root command.
func Run() {
	cmd := rootCmd(version)
	if err := cmd.Execute(); err != nil {
		kouch.Exit(err)
	}
}

// Run initializes the root command, adds subordinate commands, then executes.
func rootCmd(version string) *cobra.Command {
	cx := &kouch.CmdContext{}

	cmd := &cobra.Command{
		Use:           "kouch",
		Short:         "kouch is a command-line tool for interacting with CouchDB",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := io.RedirStderr(cmd.Flags()); err != nil {
				return err
			}
			var err error
			if cx.Verbose, err = cmd.Flags().GetBool(flagVerbose); err != nil {
				return err
			}
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

	cmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Make the operation more talkative")

	io.AddFlags(cmd.PersistentFlags())
	config.AddFlags(cmd.PersistentFlags())

	registry.AddSubcommands(cx, cmd)
	return cmd
}
