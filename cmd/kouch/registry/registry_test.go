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

func TestAddSubcommands(t *testing.T) {
	defer lockRegistry()()
	initCount := 0
	Register(nil, func() *cobra.Command {
		initCount++
		return &cobra.Command{Use: "foo"}
	})
	Register(nil, func() *cobra.Command {
		initCount++
		return &cobra.Command{Use: "bar"}
	})
	Register(nil, func() *cobra.Command {
		initCount++
		return &cobra.Command{Use: "baz"}
	})
	Register([]string{"foo"}, func() *cobra.Command {
		initCount++
		return &cobra.Command{}
	})
	AddSubcommands(&cobra.Command{})
	if expected := 4; initCount != expected {
		t.Errorf("Expected %d initializations, got %d", expected, initCount)
	}
}

func TestAddSubcommandsPanic(t *testing.T) {
	defer lockRegistry()()
	Register(nil, func() *cobra.Command {
		return &cobra.Command{Use: "foo"}
	})
	Register([]string{"foo", "bar", "baz"}, func() *cobra.Command {
		return &cobra.Command{Use: "bar"}
	})
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
	nilFn := func() *cobra.Command {
		return nil
	}
	type regTest struct {
		name     string
		init     func()
		parent   []string
		fn       CommandInitFunc
		expected interface{}
	}
	tests := []regTest{
		{
			name: "simple",
			fn:   nilFn,
			expected: &subCommand{
				children: map[string]*subCommand{},
				cmds:     []*cobra.Command{nil},
			},
		},
		{
			name:   "subcommand with no parent",
			parent: []string{"foo", "bar"},
			fn:     nilFn,
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
			Register(test.parent, test.fn)
			if d := diff.Interface(test.expected, rootCommand); d != nil {
				t.Error(d)
			}
		})
	}
}
