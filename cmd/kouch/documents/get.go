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
	registry.Register([]string{"get"}, getDocCmd)
}

func getDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "document [target]",
		Aliases: []string{"doc"},
		Short:   "Fetches a single document.",
		Long: "Fetches a single document.\n\n" +
			kouch.TargetHelpText(kouch.TargetDocument),
		RunE: getDocumentCmd,
	}
	f := cmd.Flags()
	f.String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}.")
	f.String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}.")
	f.StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves document of specified revision.")
	f.String(kouch.FlagIfNoneMatch, "", "Optionally fetch the document, only if the current rev does not match the one provided")

	cmd.PersistentFlags().BoolP(kouch.FlagHead, kouch.FlagShortHead, false, "Fetch the headers only.")
	f.Bool(kouch.FlagIncludeAttachments, false, "Include Base64-encoded content of attachments in response.")
	f.Bool(kouch.FlagIncludeAttEncoding, false, "Include encoding information in attachment stubs for compressed attachments.")
	f.StringSlice(kouch.FlagAttsSince, nil, "Include attachments only since, but not including, the specified revisions.")
	f.Bool(kouch.FlagConflicts, false, "Include document conflicts information.")
	f.Bool(kouch.FlagIncludeDeletedConflicts, false, "Include information about deleted conflicted revisions.")
	f.Bool(kouch.FlagForceLatest, false, `Force retrieving latest “leaf” revision, no matter what rev was requested.`)
	f.Bool(kouch.FlagIncludeLocalSeq, false, "Include last update sequence for the document.")
	f.Bool(kouch.FlagMeta, false, "Same as: --"+kouch.FlagConflicts+" --"+kouch.FlagIncludeDeletedConflicts+" --"+kouch.FlagRevsInfo)
	f.StringSlice(kouch.FlagOpenRevs, nil, "Retrieve documents of specified leaf revisions. May use the value 'all' to return all leaf revisions.")
	f.Bool(kouch.FlagRevs, false, "Include list of all known document revisions.")
	f.Bool(kouch.FlagRevsInfo, false, "Include detailed information for all known document revisions")
	return cmd
}

func getDocumentCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := getDocumentOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	return getDocument(ctx, o)
}

func getDocumentOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, kouch.TargetDocument, flags)
	if err != nil {
		return nil, err
	}

	o.Options.IfNoneMatch, err = flags.GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}
	if e := o.SetParams(flags,
		kouch.FlagAttsSince, kouch.FlagOpenRevs,
		kouch.FlagIncludeAttachments, kouch.FlagIncludeAttEncoding,
		kouch.FlagConflicts, kouch.FlagIncludeDeletedConflicts,
		kouch.FlagForceLatest, kouch.FlagIncludeLocalSeq, kouch.FlagMeta,
		kouch.FlagRevs, kouch.FlagRevsInfo,
	); e != nil {
		return nil, e
	}

	return o, nil
}

func getDocument(ctx context.Context, o *kouch.Options) error {
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodGet, util.DocPath(o), o)
}
