package main

import (
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/cmd"
)

const version = "0.0.1"

func main() {
	conf := viper.New()
	conf.SetDefault("server", "http://localhost:5984/")
	cmd.Run(version, conf)
}
