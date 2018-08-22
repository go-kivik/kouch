package documents

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"put"}, putDocCmd)
}

func putDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "document [target]",
		Aliases: []string{"doc"},
		Short:   "Create or update a single document.",
		Long: "Fetches a single document.\n\n" +
			target.HelpText(target.Document),
		RunE: putDocumentCmd,
	}
	f := cmd.Flags()
	f.String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}.")
	f.String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}.")
	f.StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves document of specified revision.")
	f.Bool(kouch.FlagFullCommit, false, "Overrides server’s commit policy.")

	f.Bool(flagBatch, false, "Store document in batch mode.")
	f.Bool(flagNewEdits, true, "When disabled, prevents insertion of conflicting documents.")
	return cmd
}

func putDocumentOpts(cmd *cobra.Command, _ []string) (*kouch.Options, error) {
	ctx := kouch.GetContext(cmd)
	o := kouch.NewOptions()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Document, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := o.Target.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	var err error
	o.Options.FullCommit, err = cmd.Flags().GetBool(kouch.FlagFullCommit)
	if err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}
	if e := o.SetParamString(cmd.Flags(), kouch.FlagRev); e != nil {
		return nil, e
	}
	if e := setBatch(o, cmd.Flags()); e != nil {
		return nil, e
	}
	if e := o.SetParamBool(cmd.Flags(), flagNewEdits); e != nil {
		return nil, e
	}

	return o, nil
}

func putDocumentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := putDocumentOpts(cmd, args)
	if err != nil {
		return err
	}
	result, err := putDocument(ctx, opts)
	if err != nil {
		return err
	}
	return kouch.Outputer(ctx).Output(kouch.Output(ctx), result)
}

func putDocument(ctx context.Context, o *kouch.Options) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), o.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document))
	res, err := c.DoReq(ctx, http.MethodPut, path, o.Options)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
