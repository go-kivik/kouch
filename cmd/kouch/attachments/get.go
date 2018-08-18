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
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"get"}, attCmd())
}

func attCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attachment [target]",
		Aliases: []string{"att"},
		Short:   "Fetches a file attachment.",
		Long: "Fetches a file attachment.\n\n" +
			target.HelpText(target.Attachment),
		RunE: attachmentCmd,
	}
	cmd.Flags().String(kouch.FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	cmd.Flags().String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}/{filename}.")
	cmd.Flags().String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}/{filename}")
	cmd.Flags().String(kouch.FlagIfNoneMatch, "", "Optionally fetch the attachment, only if the MD5 digest does not match the one provided")
	cmd.Flags().StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves attachment from document of specified revision.")
	return cmd
}

func attachmentCmd(cmd *cobra.Command, args []string) error {
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

type opts struct {
	*kouch.Target
	rev         string
	ifNoneMatch string
}

func getAttachmentOpts(cmd *cobra.Command, _ []string) (*opts, error) {
	ctx := kouch.GetContext(cmd)
	o := &opts{
		Target: &kouch.Target{},
	}
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Attachment, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := o.Target.FilenameFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}

	var err error
	o.ifNoneMatch, err = cmd.Flags().GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}
	o.rev, err = cmd.Flags().GetString(kouch.FlagRev)
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
	res, err := c.DoReq(context.TODO(), http.MethodGet, path, &chttp.Options{
		IfNoneMatch: o.ifNoneMatch,
	})
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}

func validateTarget(t *kouch.Target) error {
	if t.Filename == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No filename provided")
	}
	if t.Document == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No document ID provided")
	}
	if t.Database == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No database name provided")
	}
	if t.Root == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No root URL provided")
	}
	return nil
}
