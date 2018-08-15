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
	kio "github.com/go-kivik/kouch/io"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"get"}, attCmd())
}

func attCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attachment [target]",
		Aliases: []string{"att"},
		Short:   "Fetches a file attachment",
		Long: `Fetches a file attachment, and sends the content to --` + kio.FlagOutputFile + `.

Target may be of the following formats:

  - {filename} -- The filename only. Alternately, the filename may be passed with the --` + kouch.FlagFilename + ` option, particularly for filenames with slashes.
  - {id}/{filename} -- The document ID and filename.
  - /{db}/{id}/{filename} -- With leading slash, the database name, document ID, and filename.
  - http://host.com/{db}/{id}/{filename} -- A fully qualified URL, may include auth credentials.
`,
		RunE: attachmentCmd,
	}
	cmd.Flags().String(kouch.FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	cmd.Flags().String(kouch.FlagDocID, "", "The document ID. May be provided with the target in the format {id}/{filename}.")
	cmd.Flags().String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}/{filename}")
	return cmd
}

type getAttOpts struct {
	kouch.Target
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

func getAttachmentOpts(cmd *cobra.Command, args []string) (*getAttOpts, error) {
	ctx := kouch.GetContext(cmd)
	opts := &getAttOpts{}
	if len(args) > 0 {
		if len(args) > 1 {
			return nil, &errors.ExitError{
				Err:      errors.New("Too many targets provided"),
				ExitCode: chttp.ExitFailedToInitialize,
			}
		}
		var err error
		target, err := kouch.ParseAttachmentTarget(args[0])
		if err != nil {
			return nil, err
		}
		opts = &getAttOpts{*target}
	}

	if err := opts.Target.FilenameFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := opts.Target.DocIDFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := opts.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if opts.Root == "" {
			opts.Root = defCtx.Root
		}
	}

	return opts, nil
}

func getAttachment(opts *getAttOpts) (io.ReadCloser, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), opts.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s/%s", url.QueryEscape(opts.Database), chttp.EncodeDocID(opts.DocID), url.QueryEscape(opts.Filename))
	res, err := c.DoReq(context.TODO(), http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (o *getAttOpts) validate() error {
	if o.Filename == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No filename provided")
	}
	if o.DocID == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No document ID provided")
	}
	if o.Database == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No database name provided")
	}
	if o.Root == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No root URL provided")
	}
	return nil
}
