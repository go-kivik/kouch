package documents

import (
	"context"
	"net/http"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	registry.Register([]string{"put"}, putDocCmd)
}

func putDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "document [target]",
		Aliases: []string{"doc"},
		Short:   "Create or update a single document.",
		Long: "Fetches a single document.\n\n" +
			kouch.TargetHelpText(kouch.TargetDocument),
		RunE: putDocumentCmd,
	}
	f := cmd.Flags()
	f.String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}.")
	f.String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}.")
	f.StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves document of specified revision.")
	f.Bool(kouch.FlagFullCommit, false, "Overrides server’s commit policy.")
	f.BoolP(kouch.FlagAutoRev, kouch.FlagShortAutoRev, false, "Fetch the current rev before update. Use with caution!")

	f.Bool(kouch.FlagBatch, false, "Store document in batch mode.")
	f.Bool(kouch.FlagNewEdits, true, "When disabled, prevents insertion of conflicting documents.")
	return cmd
}

func putDocumentOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, kouch.TargetDocument, flags)
	if err != nil {
		return nil, err
	}

	o.Options.FullCommit, err = flags.GetBool(kouch.FlagFullCommit)
	if err != nil {
		return nil, err
	}

	if e := setBatch(o, flags); e != nil {
		return nil, e
	}
	if e := o.SetParam(flags, kouch.FlagNewEdits); e != nil {
		return nil, e
	}

	return o, nil
}

func putDocumentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := putDocumentOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodPut, util.DocPath(o), o)
}
