package uuids

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/registry"
)

type getUUIDsCtx struct {
	*kouch.CmdContext
}

func init() {
	registry.Register([]string{"get"}, getUUIDsCmd)
}

func getUUIDsCmd(cx *kouch.CmdContext) *cobra.Command {
	g := &getUUIDsCtx{cx}
	cmd := &cobra.Command{
		Use:   "uuids",
		Short: "Returns one or more server-generated UUIDs",
		Long: `Returns one or more Universally Unique Identifiers (UUIDs) from the
CouchDB server.`,
		RunE: g.getUUIDs,
	}
	cmd.Flags().IntP("count", "C", 1, "Number of UUIDs to return")
	return cmd
}

type getUUIDsOpts struct {
	Count int
	Root  string
}

func (cx *getUUIDsCtx) getUUIDsOpts(cmd *cobra.Command, _ []string) (*getUUIDsOpts, error) {
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		return nil, err
	}
	ctx, err := cx.Conf.DefaultCtx()
	if err != nil {
		return nil, err
	}
	return &getUUIDsOpts{
		Count: count,
		Root:  ctx.Root,
	}, nil
}

func (cx *getUUIDsCtx) getUUIDs(cmd *cobra.Command, args []string) error {
	opts, err := cx.getUUIDsOpts(cmd, args)
	if err != nil {
		return err
	}
	result, err := getUUIDs(opts)
	if err != nil {
		return err
	}
	return cx.Outputer.Output(cx.Output, result)
}

func getUUIDs(opts *getUUIDsOpts) (io.ReadCloser, error) {
	c, err := chttp.New(context.TODO(), opts.Root)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("/_uuids?count=%d", opts.Count)
	res, err := c.DoReq(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
