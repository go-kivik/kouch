package registry

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/log"
)

// CommandInitFunc returns a cobra sub command.
type CommandInitFunc func(log log.Logger, conf *viper.Viper) *cobra.Command

type subCommand struct {
	children  map[string]*subCommand
	initFuncs []CommandInitFunc
}

var initMU sync.Mutex
var rootCommand = newSubCommand()

func newSubCommand() *subCommand {
	return &subCommand{
		children:  make(map[string]*subCommand),
		initFuncs: []CommandInitFunc{},
	}
}

// Register registers a sub-command.
func Register(parent []string, fn CommandInitFunc) {
	initMU.Lock()
	defer initMU.Unlock()
	cmd := rootCommand
	for _, p := range parent {
		if _, ok := cmd.children[p]; !ok {
			cmd.children[p] = newSubCommand()
		}
		cmd = cmd.children[p]
	}
	cmd.initFuncs = append(cmd.initFuncs, fn)
}

// AddSubcommands initializes and adds all registered subcommands to cmd.
func AddSubcommands(cmd *cobra.Command, l log.Logger, conf *viper.Viper) {
	initMU.Lock()
	defer initMU.Unlock()
	if err := addSubcommands(cmd, l, conf, nil, rootCommand); err != nil {
		panic(err.Error())
	}
}

func addSubcommands(cmd *cobra.Command, l log.Logger, conf *viper.Viper, path []string, cmdMap *subCommand) error {
	children := make(map[string]*cobra.Command)
	for _, fn := range cmdMap.initFuncs {
		subCmd := fn(l, conf)
		var cmdName string
		if u := subCmd.Use; u != "" {
			cmdName = strings.Fields(subCmd.Use)[0]
		}
		children[cmdName] = subCmd
		cmd.AddCommand(subCmd)
	}
	for name, childCmd := range cmdMap.children {
		child, ok := children[name]
		if !ok {
			return fmt.Errorf("Subcommand '%s %s' not registered", strings.Join(path, " "), name)
		}
		if err := addSubcommands(child, l, conf, append(path, name), childCmd); err != nil {
			return err
		}
	}
	return nil
}
