package attachments

import (
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

type opts struct {
	*kouch.Target
	*chttp.Options
	rev string
}

func newOpts() *opts {
	return &opts{
		Target:  &kouch.Target{},
		Options: &chttp.Options{},
	}
}

func validateTarget(t *kouch.Target) error {
	if t.Filename == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No filename provided")
	}
	if t.Document == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No document ID provided")
	}
	if t.Database == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No database name provided")
	}
	if t.Root == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No root URL provided")
	}
	return nil
}

func commonOpts(cmd *cobra.Command, _ []string) (*opts, error) {
	ctx := kouch.GetContext(cmd)
	o := newOpts()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Attachment, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := o.Target.FilenameFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}

	var err error
	o.rev, err = cmd.Flags().GetString(kouch.FlagRev)
	if err != nil {
		return nil, err
	}

	return o, nil
}
