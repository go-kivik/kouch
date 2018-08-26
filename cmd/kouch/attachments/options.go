package attachments

import (
	"context"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/pflag"
)

func addCommonFlags(flags *pflag.FlagSet) {
	flags.String(kouch.FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	flags.String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}/{filename}.")
	flags.String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}/{filename}")
	flags.StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves attachment from document of specified revision.")
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

func commonOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o := kouch.NewOptions()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Attachment, tgt)
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

	return o, nil
}
