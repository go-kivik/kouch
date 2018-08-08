package registry

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/go-kivik/kouch"
)

// CommandInitFunc returns a cobra sub command.
type CommandInitFunc func(cx *kouch.CmdContext) *cobra.Command

type subCommand struct {
	children  map[string]*subCommand
	initFuncs []CommandInitFunc
}

var initMU sync.Mutex
var cmdTree = newSubCommand()
var rootInit CommandInitFunc

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
	cmd := cmdTree
	for _, p := range parent {
		if _, ok := cmd.children[p]; !ok {
			cmd.children[p] = newSubCommand()
		}
		cmd = cmd.children[p]
	}
	cmd.initFuncs = append(cmd.initFuncs, fn)
}

// RegisterRoot registers the root command.
func RegisterRoot(fn CommandInitFunc) {
	if rootInit != nil {
		panic("Root command already registered")
	}
	rootInit = fn
}

// Root returns the initialized, configured root command.
func Root(cx *kouch.CmdContext) *cobra.Command {
	cmd := rootInit(cx)

	for _, fn := range flagInitFuncs {
		fn(cmd.PersistentFlags())
	}

	AddSubcommands(cx, cmd)
	return cmd
}

// AddSubcommands initializes and adds all registered subcommands to cmd.
func AddSubcommands(cx *kouch.CmdContext, cmd *cobra.Command) {
	initMU.Lock()
	defer initMU.Unlock()
	if err := addSubcommands(cx, cmd, nil, cmdTree); err != nil {
		panic(err.Error())
	}
}

func addSubcommands(cx *kouch.CmdContext, cmd *cobra.Command, path []string, cmdMap *subCommand) error {
	children := make(map[string]*cobra.Command)
	for _, fn := range cmdMap.initFuncs {
		subCmd := fn(cx)
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
		if err := addSubcommands(cx, child, append(path, name), childCmd); err != nil {
			return err
		}
	}
	return nil
}

type FlagInitFunc func(*pflag.FlagSet)

var flagInitFuncs []FlagInitFunc

func RegisterFlags(fn FlagInitFunc) {
	flagInitFuncs = append(flagInitFuncs, fn)
}
