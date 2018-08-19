package attachments

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

const (
	flagContentType      = "content-type"
	flagGuessContentType = "guess-content-type"
)

const defaultContentType = "application/octet-stream"

func init() {
	registry.Register([]string{"put"}, putAttCmd())
}

func putAttCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attachment [target]",
		Aliases: []string{"att"},
		Short:   "Upload an attachment.",
		Long: "Upload the supplied content as an attachment to the specified document\n\n" +
			target.HelpText(target.Attachment),
		RunE: putAttachmentCmd,
	}
	addCommonFlags(cmd.Flags())

	cmd.Flags().String(flagContentType, "", "Attachment MIME type.")
	cmd.Flags().Bool(flagGuessContentType, false, "Attempt to guess the content type from the file. Falls back to 'application/octet-stream'.")

	return cmd
}

func putAttachmentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := putAttachmentOpts(cmd, args)
	if err != nil {
		return err
	}
	resp, err := putAttachment(ctx, opts)
	if err != nil {
		return err
	}
	defer resp.Close()
	_, err = io.Copy(kouch.Output(ctx), resp)
	return err
}

func putAttachmentOpts(cmd *cobra.Command, args []string) (*kouch.Options, error) {
	o, err := commonOpts(cmd, args)
	if err != nil {
		return nil, err
	}
	o.Options.Body = kouch.Input(kouch.GetContext(cmd))
	var ct string
	ct, err = cmd.Flags().GetString(flagContentType)
	if err != nil {
		return nil, err
	}
	if ct == "" {
		guess, err := cmd.Flags().GetBool(flagGuessContentType)
		if err != nil {
			return nil, err
		}
		if guess {
			ct = mime.TypeByExtension(filepath.Ext(o.Target.Filename))
			if ct == "" {
				ct = defaultContentType
			}
		}
	}
	o.Options.ContentType = ct
	return o, nil
}

func putAttachment(ctx context.Context, o *kouch.Options) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), o.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document), url.QueryEscape(o.Filename))
	res, err := c.DoReq(context.TODO(), http.MethodPut, path, o.Options)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}
