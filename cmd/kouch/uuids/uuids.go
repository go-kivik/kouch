package uuids

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
)

func init() {
	registry.Register([]string{"get"}, uuidsCmd)
}

func uuidsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuids [target]",
		Short: "Returns one or more server-generated UUIDs",
		Long: "Returns Universally Unique Identifiers (UUIDs) from the CouchDB server.\n\n" +
			target.HelpText(target.Root),
		RunE: getUUIDsCmd,
	}
	cmd.Flags().IntP("count", "C", 1, "Number of UUIDs to return")
	return cmd
}

func getUUIDsOpts(cmd *cobra.Command, args []string) (*kouch.Options, error) {
	ctx := kouch.GetContext(cmd)
	o := kouch.NewOptions()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Root, tgt)
		if err != nil {
			return nil, err
		}
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}

	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		return nil, err
	}
	if count != 1 {
		o.Options.Query = url.Values{"count": []string{strconv.Itoa(count)}}
	}
	return o, nil
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
	defer result.Close()
	_, err = io.Copy(kouch.Output(ctx), result)
	kouch.Output(ctx).Close()
	return err
}

func getUUIDs(ctx context.Context, o *kouch.Options) (io.ReadCloser, error) {
	c, err := chttp.New(ctx, o.Target.Root)
	if err != nil {
		return nil, err
	}
	res, err := c.DoReq(ctx, http.MethodGet, "/_uuids", o.Options)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
