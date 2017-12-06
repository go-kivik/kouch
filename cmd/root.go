package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd *cobra.Command

type cmdInitFunc func(conf *viper.Viper) *cobra.Command

var initFuncs []cmdInitFunc

func registerCommand(fn cmdInitFunc) {
	initFuncs = append(initFuncs, fn)
}

// Run initializes the root command, adds subordinate commands, then executes.
func Run(version string, conf *viper.Viper) {
	var (
		cfgFile string
		verbose bool
		server  string
	)

	rootCmd = &cobra.Command{
		Use:     "kouch",
		Short:   "kouch is a command-line tool for interacting with CouchDB",
		Version: version,
		PreRun: func(_ *cobra.Command, _ []string) {
			initConfig(conf, cfgFile)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kouch.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	conf.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "server address, optionally with schema, port, and auth credentials")
	conf.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))

	for _, fn := range initFuncs {
		rootCmd.AddCommand(fn(conf))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(conf *viper.Viper, cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		conf.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kouch" (without extension).
		conf.AddConfigPath(home)
		conf.SetConfigName(".kouch")
	}

	conf.SetEnvPrefix("KIVIK")
	conf.AutomaticEnv() // read environment variables that match

	// If a config file is found, read it.
	if err := conf.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", conf.ConfigFileUsed())
	}
}
