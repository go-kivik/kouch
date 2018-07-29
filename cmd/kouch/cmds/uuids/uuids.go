package uuids

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/cmds/registry"
)

func init() {
	registry.Register([]string{"get"}, func(cx *kouch.Context) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "uuids",
			Short: "Returns one or more server-generated UUIDs",
			Long: `Returns one or more Universally Unique Identifiers (UUIDs) from the
CouchDB server.`,
			Run: uuidsCmd,
		}
		cmd.Flags().IntP("count", "C", 1, "Number of UUIDs to return")
		return cmd
	})
}

func uuidsCmd(cmd *cobra.Command, _ []string) {
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		panic(err.Error())
	}
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%d UUIDs coming right up, from %s\n", count, url)
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
