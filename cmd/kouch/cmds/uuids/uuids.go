package uuids

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

func init() {
	registry.Register([]string{"get"}, func(cx *kouch.Context) *cobra.Command {
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

func uuidsCmd(cx *kouch.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			return err
		}
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}
		result, err := getUUIDs(url, count)
		if err != nil {
			return err
		}
		return cx.Outputer.Output(os.Stdout, result)
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
