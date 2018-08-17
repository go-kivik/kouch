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
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

// Get-doc specific flags
const (
	flagIncludeAttachments     = "attachments"
	flagIncludeAttEncoding     = "att-encoding"
	flagAttsSince              = "attachments-since"
	flagIncludeConflicts       = "conflicts"
	flagIncludeDeletedConfligs = "deleted-conflicts"
	flagForceLatest            = "latest"
	flagIncludeLocalSeq        = "local-seq"
	flagMeta                   = "meta"
	flagOpenRevs               = "open-revs"
	flagRev                    = "rev"
	flagRevs                   = "revs"
	flagRevsInfo               = "revs-info"
)

func init() {
	registry.Register([]string{"get"}, docCmd())
}

func docCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "document [target]",
		Aliases: []string{"doc"},
		Short:   "Fetches a single document.",
		Long: "Fetches a single document.\n\n" +
			target.HelpText(target.Document),
		RunE: documentCmd,
	}
	cmd.Flags().String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}.")
	cmd.Flags().String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}.")
	cmd.Flags().String(kouch.FlagIfNoneMatch, "", "Optionally fetch the document, only if the current rev does not match the one provided")
	return cmd
}

func documentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := getDocumentOpts(cmd, args)
	if err != nil {
		return err
	}
	result, err := getDocument(opts)
	if err != nil {
		return err
	}
	return kouch.Outputer(ctx).Output(kouch.Output(ctx), result)
}

type opts struct {
	*kouch.Target
	ifNoneMatch string
}

func getDocumentOpts(cmd *cobra.Command, _ []string) (*opts, error) {
	ctx := kouch.GetContext(cmd)
	opts := &opts{
		Target: &kouch.Target{},
	}
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

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if opts.Root == "" {
			opts.Root = defCtx.Root
		}
	}
	var err error
	opts.ifNoneMatch, err = cmd.Flags().GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func getDocument(o *opts) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), o.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document))
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
}

func validateTarget(t *kouch.Target) error {
	if t.Filename != "" {
		panic("non-nil filename")
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
