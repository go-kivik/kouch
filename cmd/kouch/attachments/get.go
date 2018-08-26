package attachments

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/go-kivik/kouch/io"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	registry.Register([]string{"get"}, getAttCmd)
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

func getAttachmentCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := getAttachmentOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	return getAttachment(ctx, opts)
}

func getAttachmentOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := commonOpts(ctx, flags)
	if err != nil {
		return nil, err
	}
	o.Options.IfNoneMatch, err = flags.GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func getAttachment(ctx context.Context, o *kouch.Options) error {
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	path := fmt.Sprintf("/%s/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document), url.QueryEscape(o.Filename))
	ctx = kouch.SetOutput(ctx, io.Underlying(kouch.Output(ctx)))
	return util.ChttpDo(ctx, http.MethodGet, path, o)
}
