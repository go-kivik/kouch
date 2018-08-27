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
	"github.com/spf13/pflag"
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
	o, err := putAttachmentOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodPut, util.AttPath(o), o)
}

func putAttachmentOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, target.Attachment, flags)
	if err != nil {
		return nil, err
	}

	o.Options.Body = kouch.Input(ctx)
	var ct string
	ct, err = flags.GetString(flagContentType)
	if err != nil {
		return nil, err
	}
	if ct == "" {
		guess, err := flags.GetBool(flagGuessContentType)
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
