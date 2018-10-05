package test

import (
	"context"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// OptionsTest is a generic Options Test, reusable for any command
type OptionsTest struct {
	Conf     *kouch.Config
	Args     []string
	Expected *kouch.Options
	Err      string
	Status   int
}

// Options returns an options-testing function for use with a testy.Table.
//
// cmdFn must return a new, relevant, *cobra.Command
// optsFn is the Options function to be tested. It takes a context and flags.
func Options(cmdFn func() *cobra.Command, optsFn func(context.Context, *pflag.FlagSet) (*kouch.Options, error)) func(*testing.T, OptionsTest) {
	return func(t *testing.T, test OptionsTest) {
		cmd := cmdFn()
		if err := cmd.ParseFlags(test.Args); err != nil {
			t.Fatal(err)
		}
		ctx := kouch.GetContext(cmd)
		conf := test.Conf
		if conf == nil {
			conf = &kouch.Config{}
		}
		ctx = kouch.SetConf(ctx, conf)
		if flags := cmd.Flags().Args(); len(flags) > 0 {
			ctx = kouch.SetTarget(ctx, flags[0])
		}
		kouch.SetContext(ctx, cmd)
		opts, err := optsFn(ctx, cmd.Flags())
		testy.ExitStatusError(t, test.Err, test.Status, err)
		if d := diff.Interface(test.Expected, opts); d != nil {
			t.Error(d)
		}
	}
}
