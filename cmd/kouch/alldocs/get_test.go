package alldocs

import (
	"net/url"
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"
)

/*
f.Bool(kouch.FlagUpdateSeq, false, "Whether to include in the response an `update_seq` value indicating the sequence id of the database the view reflects.")
*/

func TestGetAllDocsOpts(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("defaults", test.OptionsTest{
		Expected: &kouch.Options{
			Target:  &kouch.Target{},
			Options: &chttp.Options{},
		},
	})
	tests.Add("conflicts", test.OptionsTest{
		Args: []string{"--" + kouch.FlagConflicts},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"conflicts": []string{"true"},
				},
			},
		},
	})
	tests.Add("descending", test.OptionsTest{
		Args: []string{"--" + kouch.FlagDescending},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"descending": []string{"true"},
				},
			},
		},
	})
	tests.Add("endkey", test.OptionsTest{
		Args: []string{"--" + kouch.FlagEndKey, "oink", "--" + kouch.FlagEndKeyDocID, "moo", "--" + kouch.FlagInclusiveEnd + "=false"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"endkey":        []string{"oink"},
					"endkey_docid":  []string{"moo"},
					"inclusive_end": []string{"false"},
				},
			},
		},
	})
	tests.Add("startkey", test.OptionsTest{
		Args: []string{"--" + kouch.FlagStartKey, "oink",
			"--" + kouch.FlagStartKeyDocID, "moo"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"startkey":       []string{"oink"},
					"startkey_docid": []string{"moo"},
				},
			},
		},
	})
	tests.Add("group", test.OptionsTest{
		Args: []string{"--" + kouch.FlagGroup, "--" + kouch.FlagGroupLevel, "5"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"group":       []string{"true"},
					"group_level": []string{"5"},
				},
			},
		},
	})
	tests.Add("no reduce", test.OptionsTest{
		Args: []string{"--" + kouch.FlagReduce + "=false"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"reduce": []string{"false"},
				},
			},
		},
	})
	tests.Add("include docs", test.OptionsTest{
		Args: []string{"--" + kouch.FlagIncludeDocs},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"include_docs": []string{"true"},
				},
			},
		},
	})
	tests.Add("attachments", test.OptionsTest{
		Args: []string{"--" + kouch.FlagIncludeAttachments, "--" + kouch.FlagIncludeAttEncoding},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"attachments":       []string{"true"},
					"att_encoding_info": []string{"true"},
				},
			},
		},
	})
	tests.Add("key", test.OptionsTest{
		Args: []string{"--" + kouch.FlagKey, "oink"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"key": []string{"oink"},
				},
			},
		},
	})
	tests.Add("keys", test.OptionsTest{
		Args: []string{"--" + kouch.FlagKeys, "oink",
			"--" + kouch.FlagKeys, "moo"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"keys": []string{"oink", "moo"},
				},
			},
		},
	})
	tests.Add("limit & skip", test.OptionsTest{
		Args: []string{"--" + kouch.FlagLimit, "10",
			"--" + kouch.FlagSkip, "50"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"limit": []string{"10"},
					"skip":  []string{"50"},
				},
			},
		},
	})
	tests.Add("sorted", test.OptionsTest{
		Args: []string{"--" + kouch.FlagSorted + "=false"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"sorted": []string{"false"},
				},
			},
		},
	})
	tests.Add("stable", test.OptionsTest{
		Args: []string{"--" + kouch.FlagStable},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"stable": []string{"true"},
				},
			},
		},
	})
	tests.Add("stale", test.OptionsTest{
		Args: []string{"--" + kouch.FlagStale, "ok"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"stale": []string{"ok"},
				},
			},
		},
	})
	tests.Add("invalid stale value", test.OptionsTest{
		Args:   []string{"--" + kouch.FlagStale, "yes"},
		Err:    "Invalid value for --" + kouch.FlagStale + ". Supported options: `ok`, `update_after`, `false`",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("update", test.OptionsTest{
		Args: []string{"--" + kouch.FlagUpdate, "lazy"},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"update": []string{"lazy"},
				},
			},
		},
	})
	tests.Add("invalid update value", test.OptionsTest{
		Args:   []string{"--" + kouch.FlagUpdate, "yes"},
		Err:    "Invalid value for --" + kouch.FlagUpdate + ". Supported options: `true`, `false`, `lazy`",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("update sequence", test.OptionsTest{
		Args: []string{"--" + kouch.FlagUpdateSeq},
		Expected: &kouch.Options{
			Target: &kouch.Target{},
			Options: &chttp.Options{
				Query: url.Values{
					"update_seq": []string{"true"},
				},
			},
		},
	})

	tests.Run(t, test.Options(getAllDocsCmd, getAllDocsOpts))
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
