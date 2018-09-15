package util

import (
	"context"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/pflag"
)

// CommonOptions parses options common to most or all commands.
func CommonOptions(ctx context.Context, scope kouch.TargetScope, flags *pflag.FlagSet) (*kouch.Options, error) {
	o := kouch.NewOptions()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(scope, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := o.Target.FilenameFromFlags(flags); err != nil {
		return nil, err
	}
	if err := o.Target.DocumentFromFlags(flags); err != nil {
		return nil, err
	}
	if err := o.Target.DatabaseFromFlags(flags); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}

	if e := o.SetParamString(flags, kouch.FlagRev); e != nil {
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
