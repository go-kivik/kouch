package documents

import (
	"net/http"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
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
			target.HelpText(target.Document),
		RunE: putDocumentCmd,
	}
	f := cmd.Flags()
	f.String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}.")
	f.String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}.")
	f.StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves document of specified revision.")
	f.Bool(kouch.FlagFullCommit, false, "Overrides server’s commit policy.")
	f.BoolP(kouch.FlagAutoRev, kouch.FlagShortAutoRev, false, "Fetch the current rev before update. Use with caution!")

	f.Bool(flagBatch, false, "Store document in batch mode.")
	f.Bool(flagNewEdits, true, "When disabled, prevents insertion of conflicting documents.")
	return cmd
}

func putDocumentOpts(cmd *cobra.Command, _ []string) (*kouch.Options, error) {
	ctx := kouch.GetContext(cmd)
	o := kouch.NewOptions()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Document, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := o.Target.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	var err error
	o.Options.FullCommit, err = cmd.Flags().GetBool(kouch.FlagFullCommit)
	if err != nil {
		return nil, err
	}

	if defCtx, e := kouch.Conf(ctx).DefaultCtx(); e == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}
	autoRev, err := cmd.Flags().GetBool(kouch.FlagAutoRev)
	if err != nil {
		return nil, err
	}
	if autoRev {
		rev, err := util.FetchRev(ctx, o)
		if err != nil {
			return nil, err
		}
		o.Query().Set("rev", rev)
	} else {
		if e := o.SetParamString(cmd.Flags(), kouch.FlagRev); e != nil {
			return nil, e
		}
	}
	if e := setBatch(o, cmd.Flags()); e != nil {
		return nil, e
	}
	if e := o.SetParamBool(cmd.Flags(), flagNewEdits); e != nil {
		return nil, e
	}

	return o, nil
}

func putDocumentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := putDocumentOpts(cmd, args)
	if err != nil {
		return err
	}
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodPut, util.DocPath(o), o)
}
