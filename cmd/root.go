package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
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
	var (
		cfgFile        string
		stdout, stderr string
		clobber        bool
	)
	l := log.New()

	rootCmd = &cobra.Command{
		Use:     "kouch",
		Short:   "kouch is a command-line tool for interacting with CouchDB",
		Version: version,
	}
	cobra.OnInitialize(func() {
		initConfig(l, conf, cfgFile)
		if err := ValidateConfig(conf); err != nil {
			l.Errorln(err)
			os.Exit(ExitFailedToInitialize)
		}
		if stdout != "" {
			out, err := log.OpenLogFile(stdout, clobber)
			if err != nil {
				l.Errorln(err)
				os.Exit(ExitWriteError)
			}
			l.SetStdout(out)
		}
		if stderr != "" {
			out, err := log.OpenLogFile(stderr, clobber)
			if err != nil {
				l.Errorln(err)
				os.Exit(ExitWriteError)
			}
			l.SetStderr(out)
		}
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kouch.yaml)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	conf.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().StringP("server", "s", "", "server address, optionally with schema, port, and auth credentials")
	conf.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	rootCmd.PersistentFlags().StringVarP(&stdout, "output", "o", "", "output file to use instead of stdout")
	rootCmd.PersistentFlags().StringVarP(&stderr, "stderr", "", "", "redirect output to stderr to the specified file instead")
	rootCmd.PersistentFlags().BoolVarP(&clobber, "force", "F", false, "overwrite output files specified by --output and --stderr, if they exist")
	rootCmd.PersistentFlags().StringP("format", "f", "raw", "output format: 'raw' or pretty-formatted 'json'")
	conf.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))

	for _, fn := range initFuncs {
		rootCmd.AddCommand(fn(l, conf))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(ExitUnknownFailure)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(log log.Logger, conf *viper.Viper, cfgFile string) {
	conf.SetEnvPrefix("KIVIK")
	conf.AutomaticEnv() // read environment variables that match

	// If a config file is found, read it.
	if err := readConfigFile(conf, cfgFile); err != nil {
		fmt.Println(err)
		os.Exit(ExitFailedToInitialize)
	}
	log.SetVerbose(conf.GetBool("verbose"))
	log.Debugln("Using config file:", conf.ConfigFileUsed())
}

func readConfigFile(conf *viper.Viper, cfgFile string) error {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".kouch" (without extension).
		conf.AddConfigPath(home)
		conf.SetConfigName(".kouch")
		err = conf.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil
		}
		return err
	}
	conf.SetConfigFile(cfgFile)
	return conf.ReadInConfig()
}
