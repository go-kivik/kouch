package attachments

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"get"}, getAttCmd())
}

func getAttCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attachment [target]",
		Aliases: []string{"att"},
		Short:   "Fetches a file attachment.",
		Long: "Fetches a file attachment.\n\n" +
			target.HelpText(target.Attachment),
		RunE: getAttachmentCmd,
	}
	cmd.Flags().String(kouch.FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	cmd.Flags().String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}/{filename}.")
	cmd.Flags().String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}/{filename}")
	cmd.Flags().String(kouch.FlagIfNoneMatch, "", "Optionally fetch the attachment, only if the MD5 digest does not match the one provided")
	cmd.Flags().StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves attachment from document of specified revision.")
	return cmd
}

func getAttachmentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := getAttachmentOpts(cmd, args)
	if err != nil {
		return err
	}
	resp, err := getAttachment(opts)
	if err != nil {
		return err
	}
	defer resp.Close()
	_, err = io.Copy(kouch.Output(ctx), resp)
	return err
}

func getAttachmentOpts(cmd *cobra.Command, args []string) (*opts, error) {
	o, err := commonOpts(cmd, args)
	if err != nil {
		return nil, err
	}
	o.Options.IfNoneMatch, err = cmd.Flags().GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func getAttachment(o *opts) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), o.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document), url.QueryEscape(o.Filename))
	query := &url.Values{}
	if o.rev != "" {
		query.Add("rev", o.rev)
	}
	if eq := query.Encode(); eq != "" {
		path = path + "?" + eq
	}
	res, err := c.DoReq(context.TODO(), http.MethodGet, path, o.Options)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
