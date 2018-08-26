package documents

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
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
	return putDocument(ctx, opts)
}

func putDocument(ctx context.Context, o *kouch.Options) error {
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	path := fmt.Sprintf("/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document))
	return util.ChttpDo(ctx, http.MethodPut, path, o, kouch.HeadDumper(ctx), kouch.Output(ctx))
}
