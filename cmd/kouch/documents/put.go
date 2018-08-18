package documents

import (
	"io"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"put"}, putDocCmd())
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
	return cmd
}

func putDocumentOpts(cmd *cobra.Command, _ []string) (*opts, error) {
	ctx := kouch.GetContext(cmd)
	opts := newOpts()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		opts.Target, err = target.Parse(target.Document, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := opts.Target.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := opts.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	var err error
	opts.fullCommit, err = cmd.Flags().GetBool(kouch.FlagFullCommit)
	if err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if opts.Root == "" {
			opts.Root = defCtx.Root
		}
	}
	if e := opts.setRev(cmd.Flags()); e != nil {
		return nil, e
	}

	return opts, nil
}

func putDocumentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := getDocumentOpts(cmd, args)
	if err != nil {
		return err
	}
	result, err := putDocument(opts)
	if err != nil {
		return err
	}
	return kouch.Outputer(ctx).Output(kouch.Output(ctx), result)
}

func putDocument(o *opts) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	return nil, nil
	/*
		c, err := chttp.New(context.TODO(), o.Root)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document))
		query := o.Values
		if eq := query.Encode(); eq != "" {
			path = path + "?" + eq
		}
		res, err := c.DoReq(context.TODO(), http.MethodGet, path, &chttp.Options{
			IfNoneMatch: o.ifNoneMatch,
		})
		if err != nil {
			return nil, err
		}
		if err = chttp.ResponseError(res); err != nil {
			return nil, err
		}
		return res.Body, nil
	*/
}
