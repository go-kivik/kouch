package config

import (
	"os"
	"path"

	"github.com/go-kivik/kouch"
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
func ReadConfig() (*kouch.Config, error) {
	home := kouch.Home()
	if home != "" {
		conf, err := readConfigFile(path.Join(home, "config"))
		if err == nil || !os.IsNotExist(err) {
			return conf, err
		}
	}
	return &kouch.Config{}, nil
}
