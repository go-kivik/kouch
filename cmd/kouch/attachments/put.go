package attachments

import (
	"context"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

const (
	flagContentType      = "content-type"
	flagGuessContentType = "guess-content-type"
)

const defaultContentType = "application/octet-stream"

func init() {
	registry.Register([]string{"put"}, putAttCmd)
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
	cmd.Flags().BoolP(kouch.FlagAutoRev, kouch.FlagShortAutoRev, false, "Fetch the current rev before update. Use with caution!")

	cmd.Flags().String(flagContentType, "", "Attachment MIME type.")
	cmd.Flags().Bool(flagGuessContentType, false, "Attempt to guess the content type from the file. Falls back to 'application/octet-stream'.")

	return cmd
}

func putAttachmentCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	opts, err := putAttachmentOpts(ctx, cmd)
	if err != nil {
		return err
	}
	return putAttachment(ctx, opts)
}

func putAttachmentOpts(ctx context.Context, cmd *cobra.Command) (*kouch.Options, error) {
	o, err := commonOpts(ctx, cmd)
	if err != nil {
		return nil, err
	}

	autoRev, err := cmd.Flags().GetBool(kouch.FlagAutoRev)
	if err != nil {
		return nil, err
	}
	if autoRev {
		rev, e := util.FetchRev(ctx, o)
		if e != nil {
			return nil, e
		}
		o.Query().Set("rev", rev)
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

func putAttachment(ctx context.Context, o *kouch.Options) error {
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodPut, util.AttPath(o), o, kouch.HeadDumper(ctx), kouch.Output(ctx))
}
