package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kouch/log"
)

var rootCmd *cobra.Command

type cmdInitFunc func(log log.Logger, conf *viper.Viper) *cobra.Command

var initFuncs []cmdInitFunc

func registerCommand(fn cmdInitFunc) {
	initFuncs = append(initFuncs, fn)
}

// Run initializes the root command, adds subordinate commands, then executes.
func Run(version string, conf *viper.Viper) {
	l := log.New()

	rootCmd = &cobra.Command{
		Use:     "kouch",
		Short:   "kouch is a command-line tool for interacting with CouchDB",
		Version: version,
	}
	cobra.OnInitialize(func() {
		if err := ValidateConfig(conf); err != nil {
			l.Errorln(err)
			os.Exit(ExitFailedToInitialize)
		}
	})

	for _, fn := range initFuncs {
		rootCmd.AddCommand(fn(l, conf))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(ExitUnknownFailure)
	}
}
