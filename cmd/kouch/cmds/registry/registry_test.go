package registry

import (
	"sync"
	"testing"

	"github.com/flimzy/diff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/log"
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
	Register(nil, func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{Use: "foo"}
	})
	Register(nil, func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{Use: "bar"}
	})
	Register(nil, func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{Use: "baz"}
	})
	Register([]string{"foo"}, func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{}
	})
	AddSubcommands(&cobra.Command{}, nil, nil)
	if expected := 4; initCount != expected {
		t.Errorf("Expected %d initializations, got %d", expected, initCount)
	}
}

func TestAddSubcommandsPanic(t *testing.T) {
	defer lockRegistry()()
	Register(nil, func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		return &cobra.Command{Use: "foo"}
	})
	Register([]string{"foo", "bar", "baz"}, func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		return &cobra.Command{Use: "bar"}
	})
	recovered := func() (r interface{}) {
		defer func() { r = recover() }()
		AddSubcommands(&cobra.Command{}, nil, nil)
		return nil
	}()
	expected := "Subcommand 'foo bar' not registered"
	if d := diff.Interface(expected, recovered); d != nil {
		t.Error(d)
	}
}

func TestRegister(t *testing.T) {
	nilFn := func(_ log.Logger, _ *viper.Viper) *cobra.Command {
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
				children:  map[string]*subCommand{},
				initFuncs: []CommandInitFunc{nilFn},
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
								children:  map[string]*subCommand{},
								initFuncs: []CommandInitFunc{nilFn},
							},
						},
						initFuncs: []CommandInitFunc{},
					},
				},
				initFuncs: []CommandInitFunc{},
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