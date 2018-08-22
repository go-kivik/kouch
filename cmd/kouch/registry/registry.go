package registry

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type subCommand struct {
	children map[string]*subCommand
	cmds     []*cobra.Command
}

var initMU sync.Mutex
var rootCommand = newSubCommand()

func newSubCommand() *subCommand {
	return &subCommand{
		children: make(map[string]*subCommand),
		cmds:     []*cobra.Command{},
	}
}

var root *cobra.Command

// RegisterRoot registers the root command.
func RegisterRoot(cmd *cobra.Command) {
	if root != nil {
		panic("Root command already registered")
	}
	root = cmd
}

var configured bool

// Root initializes and returns the root command.
func Root() *cobra.Command {
	if !configured {
		AddSubcommands(root)
		configured = true
	}

	return root
}

// Register registers a sub-command.
func Register(parent []string, cmd *cobra.Command) {
	initMU.Lock()
	defer initMU.Unlock()
	rootCmd := rootCommand
	for _, p := range parent {
		if _, ok := rootCmd.children[p]; !ok {
			rootCmd.children[p] = newSubCommand()
		}
		rootCmd = rootCmd.children[p]
	}
	rootCmd.cmds = append(rootCmd.cmds, cmd)
}

// AddSubcommands initializes and adds all registered subcommands to cmd.
func AddSubcommands(cmd *cobra.Command) {
	initMU.Lock()
	defer initMU.Unlock()
	if err := addSubcommands(cmd, nil, rootCommand); err != nil {
		panic(err.Error())
	}
}

func addSubcommands(cmd *cobra.Command, path []string, cmdMap *subCommand) error {
	children := make(map[string]*cobra.Command)
	for _, subCmd := range cmdMap.cmds {
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
		if err := addSubcommands(child, append(path, name), childCmd); err != nil {
			return err
		}
	}
	return nil
}
