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
	addCommonFlags(cmd.Flags())
	cmd.PersistentFlags().BoolP(kouch.FlagHead, kouch.FlagShortHead, false, "Fetch the headers only.")
	cmd.Flags().String(kouch.FlagIfNoneMatch, "", "Optionally fetch the attachment, only if the MD5 digest does not match the one provided")
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

func getAttachmentOpts(cmd *cobra.Command, args []string) (*kouch.Options, error) {
	o, err := commonOpts(cmd, args)
	if err != nil {
		return nil, err
	}
	if e := o.SetHead(cmd.Flags()); e != nil {
		return nil, e
	}
	o.Options.IfNoneMatch, err = cmd.Flags().GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func getAttachment(o *kouch.Options) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), o.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document), url.QueryEscape(o.Filename))
	method := http.MethodGet
	if o.Head {
		method = http.MethodHead
	}
	res, err := c.DoReq(context.TODO(), method, path, o.Options)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	if o.Head {
		_ = res.Body.Close()
		r, w := io.Pipe()
		go func() {
			err := res.Header.Write(w)
			w.CloseWithError(err)
		}()
		return r, nil
	}
	return res.Body, nil
}
