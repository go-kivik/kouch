package attachments

import (
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/io"
	"github.com/spf13/cobra"
)

const (
	FlagFilename = "filename"
	FlagDocID    = "id"
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
		Use:     "attachment [filename]",
		Aliases: []string{"att"},
		Short:   "Fetches a file attachment",
		Long:    `Fetches a file attachment, and sends the content to --` + io.FlagOutputFile + `.`,
		RunE:    a.attachmentCmd,
	}
	cmd.Flags().String(FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	cmd.Flags().String(FlagDocID, "", "The document ID. May be provided with the filename argument in the format {id}/{filename}.")
	return cmd
}

/*
   The single argument may be of one of the following formats:
   nil -- in which case the --filename argument is necessary
   {filename}  -- Assumes root url, db, and doc-id are provided with other flags
   {docid}/{filename} -- attachments or docs with slashes must use --filename and --docid respectively
   /{db}/{docid}/{filename}
   http://url/{db}/{docid}/{filename}
*/

type getAttOpts struct {
	filename string
}

func (cx *attCmdCtx) attachmentCmd(cmd *cobra.Command, args []string) error {
	return getAttachment(cx.CmdContext, cmd, args)
}

func (cx *attCmdCtx) getAttachmentOpts(cmd *cobra.Command, args []string) (*getAttOpts, error) {
	if len(args) < 1 || len(args) > 1 {
		return nil, &errors.ExitError{
			Err:      errors.New("Must provide exactly one filename"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	return nil, nil
}

func getAttachment(cx *kouch.CmdContext, cmd *cobra.Command, args []string) error {
	filename, err := cmd.Flags().GetString(FlagFilename)
	if err != nil {
		return err
	}
	if filename != "" && len(args) > 0 {
		return &errors.ExitError{
			Err:      errors.New("Must use --" + FlagFilename + " and pass separate filename"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	if len(args) < 1 || len(args) > 2 {
		return &errors.ExitError{
			Err:      errors.New("Must provide exactly one filename"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}

	return nil
}
