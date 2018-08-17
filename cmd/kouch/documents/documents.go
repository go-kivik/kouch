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

func getDocumentOpts(cmd *cobra.Command, _ []string) (*kouch.Target, error) {
	ctx := kouch.GetContext(cmd)
	t := &kouch.Target{}
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		t, err = target.Parse(target.Document, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := t.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := t.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if t.Root == "" {
			t.Root = defCtx.Root
		}
	}

	return t, nil
}

func getDocument(t *kouch.Target) (io.ReadCloser, error) {
	if err := validateTarget(t); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), t.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s", url.QueryEscape(t.Database), chttp.EncodeDocID(t.Document))
	res, err := c.DoReq(context.TODO(), http.MethodGet, path, nil)
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
