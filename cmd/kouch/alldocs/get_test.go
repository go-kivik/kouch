package alldocs

import (
	"testing"

	"github.com/flimzy/testy"
)

func TestGetAllDocsOpts(t *testing.T) {
	type gadoTest struct {
	}

	tests := testy.NewTable()

	tests.Run(t, func(t *testing.T, test gadoTest) {

	})
}

/*
TODO: no-flag functionality
Each flag below:

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
f.Bool(kouch.FlagStable, false, "Whether or not the view results should be returned from a stable set of shards.")
f.String(kouch.FlagStale, "false", "Allow the results from a stale view to be used.")
f.String(kouch.FlagStartKey, "", "Return records starting with the specified key.")
f.String(kouch.FlagStartKeyDocID, "", "Return records starting with the specified document ID. Ignored if --"+kouch.FlagStartKey+" is not set.")
f.String(kouch.FlagUpdate, "true", "Whether or not the view in question should be updated prior to responding to the user. Supported values: `true`, ``false`, `lazy`.")
f.Bool(kouch.FlagUpdateSeq, false, "Whether to include in the response an `update_seq` value indicating the sequence id of the database the view reflects.")
*/
