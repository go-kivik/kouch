package registry

import (
	"sync"
	"testing"

	"github.com/flimzy/diff"
	"github.com/spf13/cobra"
)

var registryLock sync.Mutex

func lockRegistry() func() {
	registryLock.Lock()
	return func() {
		rootCommand = newSubCommand()
		registryLock.Unlock()
	}
}

func TestAddSubcommandsPanic(t *testing.T) {
	defer lockRegistry()()
	Register(nil, &cobra.Command{Use: "foo"})
	Register([]string{"foo", "bar", "baz"}, &cobra.Command{Use: "bar"})
	recovered := func() (r interface{}) {
		defer func() { r = recover() }()
		AddSubcommands(&cobra.Command{})
		return nil
	}()
	expected := "Subcommand 'foo bar' not registered"
	if d := diff.Interface(expected, recovered); d != nil {
		t.Error(d)
	}
}

func TestRegister(t *testing.T) {
	type regTest struct {
		name     string
		init     func()
		parent   []string
		cmd      *cobra.Command
		expected interface{}
	}
	tests := []regTest{
		{
			name: "simple",
			cmd:  nil,
			expected: &subCommand{
				children: map[string]*subCommand{},
				cmds:     []*cobra.Command{nil},
			},
		},
		{
			name:   "subcommand with no parent",
			parent: []string{"foo", "bar"},
			cmd:    nil,
			expected: &subCommand{
				children: map[string]*subCommand{
					"foo": &subCommand{
						children: map[string]*subCommand{
							"bar": &subCommand{
								children: map[string]*subCommand{},
								cmds:     []*cobra.Command{nil},
							},
						},
						cmds: []*cobra.Command{},
					},
				},
				cmds: []*cobra.Command{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer lockRegistry()()
			if test.init != nil {
				test.init()
			}
			Register(test.parent, test.cmd)
			if d := diff.Interface(test.expected, rootCommand); d != nil {
				t.Error(d)
			}
		})
	}
}
