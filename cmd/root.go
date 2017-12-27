package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kouch/log"
)

const version = "0.0.1"

// Run is the entry point, which executes the root command.
func Run() {
	l := log.New()

	cmd := rootCmd(l, viper.New(), version)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(ExitUnknownFailure)
	}
}

func onInit(l log.Logger, conf *viper.Viper) func() {
	return func() {
		if err := ValidateConfig(conf); err != nil {
			l.Errorln(err)
			os.Exit(ExitFailedToInitialize)
		}
	}
}

// Run initializes the root command, adds subordinate commands, then executes.
func rootCmd(l log.Logger, conf *viper.Viper, version string) *cobra.Command {
	var rootCmd *cobra.Command

	rootCmd = &cobra.Command{
		Use:     "kouch",
		Short:   "kouch is a command-line tool for interacting with CouchDB",
		Version: version,
	}

	cobra.OnInitialize(onInit(l, conf))
	AddSubcommands(rootCmd, l, conf)
	return rootCmd
}
