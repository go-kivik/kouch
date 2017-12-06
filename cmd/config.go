package cmd

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("server", "http://localhost:5984/")
}
