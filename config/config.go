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
	if e := yaml.NewDecoder(r).Decode(&conf); e != nil {
		return nil, e
	}
	conf.File = file
	return conf, nil
}

// ReadConfig reads the config from files, env, and/or command-line arguments.
func ReadConfig(cmd *cobra.Command) (*kouch.Config, error) {
	conf, err := fileConf(cmd)
	if err != nil {
		return nil, err
	}
	root, err := cmd.Flags().GetString(flagServerRoot)
	if err != nil {
		return nil, err
	}
	if root != "" {
		conf.DefaultContext = dynamicContextName
		conf.Contexts = append(conf.Contexts, kouch.NamedContext{
			Name: dynamicContextName,
			Context: &kouch.Context{
				Root: root,
			},
		})
	}
	return conf, nil
}

func fileConf(cmd *cobra.Command) (*kouch.Config, error) {
	cfgFile, err := cmd.Flags().GetString(flagConfigFile)
	if err != nil {
		return nil, err
	}
	if cfgFile != "" {
		return readConfigFile(cfgFile)
	}
	home := Home()
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
	cmd.PersistentFlags().String(flagConfigFile, "", "Path to the kouchconfig file to use for CLI requests")
	cmd.PersistentFlags().StringP(flagServerRoot, "r", "", "The default context's root URL")
}
