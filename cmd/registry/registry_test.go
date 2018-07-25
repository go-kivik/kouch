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
		initFuncs = []CommandInitFunc{}
		registryLock.Unlock()
	}
}

func TestAddSubcommands(t *testing.T) {
	defer lockRegistry()()
	initCount := 0
	Register(func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{}
	})
	Register(func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{}
	})
	Register(func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		initCount++
		return &cobra.Command{}
	})
	AddSubcommands(&cobra.Command{}, nil, nil)
	if initCount != 3 {
		t.Errorf("Expected 3 initializations, got %d", initCount)
	}
}

func TestRegister(t *testing.T) {
	type regTest struct {
		name      string
		fn        CommandInitFunc
		expected  interface{}
		recovered interface{}
	}
	tests := []regTest{
		func() regTest {
			fn := func(_ log.Logger, _ *viper.Viper) *cobra.Command {
				return nil
			}
			return regTest{
				name:     "simple",
				fn:       fn,
				expected: []CommandInitFunc{fn},
			}
		}(),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer lockRegistry()()
			recovered := func() (r interface{}) {
				defer func() { r = recover() }()
				Register(test.fn)
				return nil
			}()
			if d := diff.Interface(test.recovered, recovered); d != nil {
				t.Error(d)
			}
			if d := diff.Interface(test.expected, initFuncs); d != nil {
				t.Error(d)
			}
		})
	}
}
