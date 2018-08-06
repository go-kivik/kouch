package config

import (
	"os"
	"path"

	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// readConfigFile reads the config file found at file.
func readConfigFile(file string) (*kouch.Config, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var conf *kouch.Config
	err = yaml.NewDecoder(r).Decode(&conf)
	return conf, err
}

// ReadConfig reads the config from files, env, and/or command-line arguments.
func ReadConfig(cmd *cobra.Command) (*kouch.Config, error) {
	cfgFile, err := cmd.Flags().GetString(kouch.FlagConfigFile)
	if err != nil {
		return nil, err
	}
	if cfgFile != "" {
		return readConfigFile(cfgFile)
	}
	home := kouch.Home()
	if home != "" {
		conf, err := readConfigFile(path.Join(home, "config"))
		if err == nil || !os.IsNotExist(err) {
			return conf, err
		}
	}
	return &kouch.Config{}, nil
}

// AddFlags adds command line flags for global config options.
func AddFlags(cmd *cobra.Command) {
	cmd.Flags().String(kouch.FlagConfigFile, "", "Path to the kouchconfig file to use for CLI requests.")
	cmd.PersistentFlags().StringP("url", "u", "", "The default context's root URL")
}
