package cmd

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

func TestRegister(t *testing.T) {
	defer lockRegistry()()
	fn := func(_ log.Logger, _ *viper.Viper) *cobra.Command {
		return nil
	}
	Register(fn)
	expected := []CommandInitFunc{fn}
	if d := diff.Interface(expected, initFuncs); d != nil {
		t.Error(d)
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
