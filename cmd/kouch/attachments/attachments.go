package attachments

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/io"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flags for necessary arguments
const (
	FlagFilename = "filename"
	FlagDocID    = "id"
	FlagDatabase = "database"
)

type attCmdCtx struct {
	*kouch.CmdContext
}

func init() {
	registry.Register([]string{"get"}, attCmd)
}

func attCmd(cx *kouch.CmdContext) *cobra.Command {
	a := &attCmdCtx{cx}
	cmd := &cobra.Command{
		Use:     "attachment [target]",
		Aliases: []string{"att"},
		Short:   "Fetches a file attachment",
		Long:    `Fetches a file attachment, and sends the content to --` + io.FlagOutputFile + `.`,
		RunE:    a.attachmentCmd,
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

func (cx *attCmdCtx) attachmentCmd(cmd *cobra.Command, args []string) error {
	opts, err := cx.getAttachmentOpts(cmd, args)
	if err != nil {
		return err
	}
	return getAttachment(opts)
}

/*
   The single argument may be of one of the following formats:
   nil -- in which case the --filename argument is necessary
   {filename}  -- Assumes root url, db, and doc-id are provided with other flags
   {docid}/{filename} -- doc id with slashes must use --docid. If there are multiple slashes in the target, the first one is considered a separator between docID and filename, and subsequent ones are part of the filename name.
   /{db}/{docid}/{filename}
   http://url/{db}/{docid}/{filename}
*/

func (cx *attCmdCtx) getAttachmentOpts(cmd *cobra.Command, args []string) (*getAttOpts, error) {
	opts := &getAttOpts{}
	if len(args) < 1 || len(args) > 1 {
		return nil, &errors.ExitError{
			Err:      errors.New("Must provide exactly one filename"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	var err error
	if opts.root, opts.db, opts.id, opts.filename, err = parseTarget(args[0]); err != nil {
		return nil, err
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

	if opts.id == "" {
		return nil, &errors.ExitError{
			Err:      errors.New("No document ID provided"),
			ExitCode: chttp.ExitFailedToInitialize,
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

func getAttachment(opts *getAttOpts) error {
	return nil
}

func parseTarget(target string) (root, db, id, filename string, err error) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		url, err := url.Parse(target)
		if err != nil {
			return "", "", "", "", &errors.ExitError{Err: err, ExitCode: chttp.ExitStatusURLMalformed}
		}
		_, db, id, filename, _ := parseTarget(url.Path)
		return fmt.Sprintf("%s://%s/", url.Scheme, url.Host),
			db, id, filename, nil
	}
	if strings.HasPrefix(target, "/") {
		parts := strings.SplitN(target, "/", 4)
		return "", parts[1], parts[2], parts[3], nil
	}
	if strings.Contains(target, "/") {
		parts := strings.SplitN(target, "/", 2)
		return "", "", parts[0], parts[1], nil
	}
	return "", "", "", target, nil
}
