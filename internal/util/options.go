package util

import (
	"context"

	"github.com/go-kivik/kouch"
	"github.com/spf13/pflag"
)

// CommonOptions parses options common to most or all commands.
func CommonOptions(ctx context.Context, scope kouch.TargetScope, flags *pflag.FlagSet) (*kouch.Options, error) {
	o := kouch.NewOptions()
	var err error
	o.Target, err = kouch.NewTarget(ctx, scope, flags)
	if err != nil {
		return nil, err
	}

	if e := o.SetParam(flags, kouch.FlagRev); e != nil {
		return nil, e
	}
	if e := setAutoRev(ctx, o, flags); e != nil {
		return nil, e
	}

	return o, nil
}

func setAutoRev(ctx context.Context, o *kouch.Options, flags *pflag.FlagSet) error {
	if flag := flags.Lookup(kouch.FlagAutoRev); flag == nil {
		return nil
	}
	autoRev, err := flags.GetBool(kouch.FlagAutoRev)
	if err != nil {
		return err
	}
	if !autoRev {
		return nil
	}
	rev, e := FetchRev(ctx, o)
	if e != nil {
		return e
	}
	o.Query().Set("rev", rev)
	return nil
}
