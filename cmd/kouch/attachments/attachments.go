package attachments

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/errors"
	kio "github.com/go-kivik/kouch/io"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flags for necessary arguments
const (
	FlagFilename = "filename"
	FlagDocID    = "id"
	FlagDatabase = "database"
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

  - {filename} -- The filename only. Alternately, the filename may be passed with the --` + FlagFilename + ` option, particularly for filenames with slashes.
  - {id}/{filename} -- The document ID and filename.
  - /{db}/{id}/{filename} -- With leading slash, the database name, document ID, and filename.
  - http://host.com/{db}/{id}/{filename} -- A fully qualified URL, may include auth credentials.
`,
		RunE: attachmentCmd,
	}
	cmd.Flags().String(FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	cmd.Flags().String(FlagDocID, "", "The document ID. May be provided with the target in the format {id}/{filename}.")
	cmd.Flags().String(FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}/{filename}")
	return cmd
}

type getAttOpts struct {
	url      string
	root     string
	db       string
	id       string
	filename string
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
		opts, err = parseTarget(args[0])
		if err != nil {
			return nil, err
		}
	}

	if err := opts.filenameFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := opts.idFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := opts.dbFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if opts.root == "" {
			opts.root = defCtx.Root
		}
	}

	return opts, nil
}

func (o *getAttOpts) filenameFromFlags(flags *pflag.FlagSet) error {
	fn, err := flags.GetString(FlagFilename)
	if err != nil {
		return err
	}
	if fn == "" {
		return nil
	}
	if o.filename != "" {
		return &errors.ExitError{
			Err:      errors.New("Must not use --" + FlagFilename + " and pass separate filename"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	o.filename = fn
	return nil
}

func (o *getAttOpts) idFromFlags(flags *pflag.FlagSet) error {
	id, err := flags.GetString(FlagDocID)
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	if o.id != "" {
		return &errors.ExitError{
			Err:      errors.New("Must not use --" + FlagDocID + " and pass doc ID as part of the target"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	o.id = id
	return nil
}

func (o *getAttOpts) dbFromFlags(flags *pflag.FlagSet) error {
	db, err := flags.GetString(FlagDatabase)
	if err != nil {
		return err
	}
	if db == "" {
		return nil
	}
	if o.db != "" {
		return &errors.ExitError{
			Err:      errors.New("Must not use --" + FlagDatabase + " and pass database as part of the target"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	o.db = db
	return nil
}

func getAttachment(opts *getAttOpts) (io.ReadCloser, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), opts.root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s/%s", url.QueryEscape(opts.db), chttp.EncodeDocID(opts.id), url.QueryEscape(opts.filename))
	res, err := c.DoReq(context.TODO(), http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}

func parseTarget(target string) (*getAttOpts, error) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		url, err := url.Parse(target)
		if err != nil {
			return nil, &errors.ExitError{Err: err, ExitCode: chttp.ExitStatusURLMalformed}
		}
		opts, err := parseTarget(url.Path)
		opts.root = fmt.Sprintf("%s://%s/", url.Scheme, url.Host)
		return opts, err
	}
	if strings.HasPrefix(target, "/") {
		parts := strings.SplitN(target, "/", 4)
		if len(parts) < 4 {
			return nil, errors.NewExitError(chttp.ExitStatusURLMalformed, "invalid target")
		}
		return &getAttOpts{
			db:       parts[1],
			id:       parts[2],
			filename: parts[3],
		}, nil
	}
	if strings.Contains(target, "/") {
		parts := strings.SplitN(target, "/", 2)
		return &getAttOpts{
			id:       parts[0],
			filename: parts[1],
		}, nil
	}
	return &getAttOpts{filename: target}, nil
}

func (o *getAttOpts) validate() error {
	if o.filename == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No filename provided")
	}
	if o.id == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No document ID provided")
	}
	if o.db == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No database name provided")
	}
	if o.root == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No root URL provided")
	}
	return nil
}
