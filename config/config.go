package config

import (
	"io"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

// readConfigFile reads the config file found at file.
func readConfigFile(file string) (*kouch.Config, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	conf := new(kouch.Config)
	if e := yaml.NewDecoder(r).Decode(&conf); e != nil {
		if e == io.EOF {
			return conf, nil
		}
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
	root, err := cmd.Flags().GetString(kouch.FlagServerRoot)
	if err != nil {
		return nil, err
	}
	user, err := cmd.Flags().GetString(kouch.FlagUser)
	if err != nil {
		return nil, err
	}
	var password string
	if !cmd.Flags().Changed(kouch.FlagPassword) {
		parts := append(strings.SplitN(user, ":", 2), "")
		user, password = parts[0], parts[1]
	} else {
		password, err = cmd.Flags().GetString(kouch.FlagPassword)
		if err != nil {
			return nil, err
		}
	}
	if root != "" || user != "" || password != "" {
		conf.DefaultContext = dynamicContextName
		conf.Contexts = append(conf.Contexts, kouch.NamedContext{
			Name: dynamicContextName,
			Context: &kouch.Context{
				Root:     root,
				User:     user,
				Password: password,
			},
		})
	}
	context, err := cmd.Flags().GetString(flagContext)
	if err != nil {
		return nil, err
	}
	if context != "" {
		conf.DefaultContext = context
	}
	return conf, nil
}

func fileConf(cmd *cobra.Command) (*kouch.Config, error) {
	cfgFile, err := cmd.Flags().GetString(kouch.FlagConfigFile)
	if err != nil {
		return nil, err
	}
	if cfgFile != "" {
		return readConfigFile(cfgFile)
	}
	home := Home()
	if home != "" {
		conf, err := readConfigFile(path.Join(home, "config"))
		if err == nil || !isNotExist(err) {
			return conf, err
		}
	}
	return &kouch.Config{}, nil
}

func isNotExist(err error) bool {
	if os.IsNotExist(err) {
		return true
	}
	if pe, ok := err.(*os.PathError); ok {
		if errno, ok := pe.Err.(syscall.Errno); ok {
			if errno == syscall.ENOTDIR {
				return true
			}
		}
	}
	return false
}

// AddFlags adds command line flags for global config options.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(kouch.FlagConfigFile, "", "Path to the kouchconfig file to use for CLI requests")
	flags.StringP(kouch.FlagServerRoot, kouch.FlagShortServerRoot, "", "The root URL")
	flags.String(flagContext, "", "The named context to use")
	flags.StringP(kouch.FlagUser, kouch.FlagShortUser, "", "Specify the username, and possibly password, to user for server authentication. If the password is not set with the "+kouch.FlagShortPassword+"/"+kouch.FlagPassword+" option, then the first colon in this option will be considered a separator for the username and password. To specificy a username with a colon, you must provide a password as a separate option.")
	flags.StringP(kouch.FlagPassword, kouch.FlagShortPassword, "", "Specify the password for server authentication.")
}
