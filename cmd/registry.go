package cmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/log"
)

// CommandInitFunc returns a cobra sub command.
type CommandInitFunc func(log log.Logger, conf *viper.Viper) *cobra.Command

var initMU sync.Mutex
var initFuncs []CommandInitFunc

// Register registers a sub-command.
func Register(fn CommandInitFunc) {
	initMU.Lock()
	defer initMU.Unlock()
	initFuncs = append(initFuncs, fn)
}

// AddSubcommands initializes and adds all registered subcommands to cmd.
func AddSubcommands(cmd *cobra.Command, l log.Logger, conf *viper.Viper) {
	initMU.Lock()
	defer initMU.Unlock()
	for _, fn := range initFuncs {
		cmd.AddCommand(fn(l, conf))
	}
}
