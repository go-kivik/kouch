package kouch

import (
	"github.com/spf13/pflag"
)

type initFlagSet func(*pflag.FlagSet)

func flagSet(init ...initFlagSet) *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.PanicOnError)
	for _, fn := range init {
		fn(fs)
	}
	return fs
}
