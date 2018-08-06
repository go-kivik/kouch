package config

import (
	"os"

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
