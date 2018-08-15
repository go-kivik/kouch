package uuids

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/errors"
)

func init() {
	registry.Register([]string{"get"}, uuidsCmd())
}

func uuidsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuids [target]",
		Short: "Returns one or more server-generated UUIDs",
		Long: `Returns one or more Universally Unique Identifiers (UUIDs) from the
CouchDB server.`,
		RunE: getUUIDsCmd,
	}
	cmd.Flags().IntP("count", "C", 1, "Number of UUIDs to return")
	return cmd
}

type opts struct {
	root  string
	count int
}

func getUUIDsOpts(cmd *cobra.Command, args []string) (*opts, error) {
	ctx := kouch.GetContext(cmd)
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		return nil, err
	}
	var root string
	if len(args) > 0 {
		if len(args) > 1 {
			return nil, errors.NewExitError(chttp.ExitFailedToInitialize, "Too many targets provided")
		}
		root = args[0]
	}
	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil && root == "" {
		root = defCtx.Root
	}
	return &opts{
		root:  root,
		count: count,
	}, nil
}

func getUUIDsCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := getUUIDsOpts(cmd, args)
	if err != nil {
		return err
	}
	result, err := getUUIDs(ctx, opts)
	if err != nil {
		return err
	}
	return kouch.Outputer(ctx).Output(kouch.Output(ctx), result)
}

func getUUIDs(ctx context.Context, opts *opts) (io.ReadCloser, error) {
	c, err := chttp.New(ctx, opts.root)
	if err != nil {
		return nil, err
	}
	res, err := c.DoReq(ctx, http.MethodGet, fmt.Sprintf("/_uuids?count=%d", opts.count), nil)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
