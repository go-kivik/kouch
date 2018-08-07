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
)

func init() {
	registry.Register([]string{"get"}, func(cx *kouch.CmdContext) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "uuids",
			Short: "Returns one or more server-generated UUIDs",
			Long: `Returns one or more Universally Unique Identifiers (UUIDs) from the
CouchDB server.`,
			RunE: uuidsCmd(cx),
		}
		cmd.Flags().IntP("count", "C", 1, "Number of UUIDs to return")
		return cmd
	})
}

func uuidsCmd(cx *kouch.CmdContext) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			return err
		}
		ctx, err := cx.Conf.DefaultCtx()
		if err != nil {
			return err
		}
		result, err := getUUIDs(ctx.Root, count)
		if err != nil {
			return err
		}
		return cx.Outputer.Output(cx.Output, result)
	}
}

func getUUIDs(url string, count int) (io.ReadCloser, error) {
	c, err := chttp.New(context.TODO(), url)
	if err != nil {
		return nil, err
	}
	res, err := c.DoReq(context.TODO(), http.MethodGet, fmt.Sprintf("/_uuids?count=%d", count), nil)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
