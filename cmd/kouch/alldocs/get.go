package alldocs

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
	registry.Register([]string{"get"}, getAllDocsCmd)
}

func getAllDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alldocs [target]",
		Aliases: []string{"doc"},
		Short:   "Fetches all documents in a database.",
		Long: "Fetches all documents in a database, subject to possible restrictions.\n\n" +
			kouch.TargetHelpText(kouch.TargetDatabase),
		RunE: getAllDocumentsCmd,
	}
	f := cmd.Flags()
	f.Bool(kouch.FlagConflicts, false, "Include conflicts information in response. Ignored if --"+kouch.FlagIncludeDocs+" isn’t true.")
	f.Bool(kouch.FlagDescending, false, "Return the documents in descending order by key.")
	f.String(kouch.FlagEndKey, "", "Stop returning records when the specified key, in JSON format, is reached.")
	f.String(kouch.FlagEndKeyDocID, "", "Stop returning records when the specified document ID is reached. Ignored if --"+kouch.FlagEndKey+" is not set.")
	f.Bool(kouch.FlagGroup, false, "Group the results using the reduce function to a group or single row. Implies --"+kouch.FlagReduce+" is true and the maximum --"+kouch.FlagGroupLevel+" value.")
	f.Int(kouch.FlagGroupLevel, 0, "Specify the group level to be used. Implies --"+kouch.FlagGroup+" is true.")
	f.Bool(kouch.FlagIncludeDocs, false, "Include the associated document with each row.")
	f.Bool(kouch.FlagIncludeAttachments, false, "Include Base64-encoded content of attachments in the response if --"+kouch.FlagIncludeDocs+" is true. Ignored if --"+kouch.FlagIncludeDocs+" is not true.")
	f.Bool(kouch.FlagIncludeAttEncoding, false, "Include encoding information in attachment stubs for compressed attachments if --"+kouch.FlagIncludeDocs+" is true.")
	f.Bool(kouch.FlagInclusiveEnd, true, "Specifies whether the specified end key should be included in the result.")
	f.String(kouch.FlagKey, "", "Return only documents that match the specified key in JSON format.")
	f.StringArray(kouch.FlagKeys, []string{}, "Return only documents matching one the keys specified in the array.")
	f.Int(kouch.FlagLimit, 0, "The maximum number of documents to be returned.")
	f.Bool(kouch.FlagReduce, true, "Use the reduction function. Default is true when a reduce function is defined, false otherwise.")
	f.Int(kouch.FlagSkip, 0, "Skip this number of records before starting to return the results.")
	f.Bool(kouch.FlagSorted, true, "Sort returned rows. Setting this to false offers a performance boost. The `total_rows` and `offset` fields are not available in the result when this is disabled.")
	f.Bool(kouch.FlagStable, false, "Whether or not the view results should be returned from a stable set of shards. Supported values: `ok`, `update_after` and `false`.")
	f.String(kouch.FlagStale, "false", "Allow the results from a stale view to be used.")
	f.String(kouch.FlagStartKey, "", "Return records starting with the specified key.")
	f.String(kouch.FlagStartKeyDocID, "", "Return records starting with the specified document ID. Ignored if --"+kouch.FlagStartKey+" is not set.")
	f.String(kouch.FlagUpdate, "true", "Whether or not the view in question should be updated prior to responding to the user. Supported values: `true`, `false`, `lazy`.")
	f.Bool(kouch.FlagUpdateSeq, false, "Whether to include in the response an `update_seq` value indicating the sequence id of the database the view reflects.")

	return cmd
}

func getAllDocumentsCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := getAllDocsOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	return getAllDocs(ctx, o)
}

func getAllDocsOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, kouch.TargetDocument, flags)
	if err != nil {
		return nil, err
	}

	if e := o.SetParams(flags,
		kouch.FlagEndKey, kouch.FlagEndKeyDocID, kouch.FlagKey, kouch.FlagStale,
		kouch.FlagStartKey, kouch.FlagStartKeyDocID, kouch.FlagUpdate,
		kouch.FlagKeys, kouch.FlagGroupLevel, kouch.FlagLimit, kouch.FlagSkip,
	); e != nil {
		return nil, e
	}

	for _, flag := range []string{
		kouch.FlagConflicts, kouch.FlagDescending, kouch.FlagGroup,
		kouch.FlagIncludeDocs, kouch.FlagIncludeAttachments,
		kouch.FlagIncludeAttEncoding, kouch.FlagInclusiveEnd, kouch.FlagReduce,
		kouch.FlagSorted, kouch.FlagStable, kouch.FlagUpdateSeq,
	} {
		if e := o.SetParamBool(flags, flag); e != nil {
			return nil, e
		}
	}

	return o, nil
}

func getAllDocs(ctx context.Context, o *kouch.Options) error {
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodGet, util.DocPath(o), o)
}
