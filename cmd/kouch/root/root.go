package root

import (
	"context"

	"github.com/spf13/cobra"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/config"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/io"

	// The individual sub-commands
	_ "github.com/go-kivik/kouch/cmd/kouch/attachments"
	_ "github.com/go-kivik/kouch/cmd/kouch/config"
	_ "github.com/go-kivik/kouch/cmd/kouch/documents"
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
	cmd := &cobra.Command{
		Use:           "kouch [options] [command] [target]",
		Short:         "kouch is a command-line tool for interacting with CouchDB",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := kouch.GetContext(cmd)
			ctx, err := setTarget(ctx, args)
			if err != nil {
				return err
			}
			if e := io.RedirStderr(cmd.Flags()); e != nil {
				return e
			}
			ctx, err = verbose(ctx, cmd)
			if err != nil {
				return err
			}
			output, err := io.SelectOutput(cmd)
			if err != nil {
				return err
			}
			ctx = kouch.SetOutput(ctx, output)
			outputer, err := io.SelectOutputProcessor(cmd)
			if err != nil {
				return err
			}
			ctx = kouch.SetOutputer(ctx, outputer)
			conf, err := config.ReadConfig(cmd)
			if err != nil {
				return err
			}
			ctx = kouch.SetConf(ctx, conf)

			input, err := io.SelectInput(cmd)
			if err != nil {
				return err
			}
			ctx = kouch.SetInput(ctx, input)

			kouch.SetContext(ctx, cmd)
			return nil
		},
	}

	cmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Make the operation more talkative")

	io.AddFlags(cmd.PersistentFlags())
	config.AddFlags(cmd.PersistentFlags())

	registry.AddSubcommands(cmd)
	return cmd
}

func setTarget(ctx context.Context, args []string) (context.Context, error) {
	if len(args) == 0 {
		return ctx, nil
	}
	if len(args) > 1 {
		return nil, errors.NewExitError(chttp.ExitFailedToInitialize, "Too many targets provided")
	}
	return kouch.SetTarget(ctx, args[0]), nil
}
