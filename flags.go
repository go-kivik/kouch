package kouch

import (
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

// Common command line flags
const (
	// Curl-equivalent flags
	FlagVerbose    = "verbose"
	FlagOutputFile = "output"
	FlagData       = "data"
	FlagHead       = "head"
	FlagDumpHeader = "dump-header"
	FlagUser       = "user"
	FlagCreateDirs = "create-dirs"

	// Custom flags
	FlagClobber                 = "force"
	FlagConfigFile              = "kouchconfig"
	FlagServerRoot              = "root"
	FlagDataJSON                = "data-json"
	FlagDataYAML                = "data-yaml"
	FlagOutputFormat            = "output-format"
	FlagFilename                = "filename"
	FlagDocument                = "id"
	FlagDatabase                = "database"
	FlagFullCommit              = "full-commit"
	FlagIfNoneMatch             = "if-none-match"
	FlagRev                     = "rev"
	FlagAutoRev                 = "auto-rev"
	FlagShards                  = "shards"
	FlagPassword                = "password"
	FlagContext                 = "context"
	FlagConflicts               = "conflicts"
	FlagDescending              = "descending"
	FlagEndKey                  = "endkey"
	FlagEndKeyDocID             = "endkey-docid"
	FlagGroup                   = "group"
	FlagGroupLevel              = "group-level"
	FlagIncludeDocs             = "include-docs"
	FlagIncludeAttachments      = "attachments"
	FlagIncludeAttEncoding      = "att-encoding-info"
	FlagInclusiveEnd            = "inclusive-end"
	FlagKey                     = "key"
	FlagKeys                    = "keys"
	FlagLimit                   = "limit"
	FlagReduce                  = "reduce"
	FlagSkip                    = "skip"
	FlagSorted                  = "sorted"
	FlagStable                  = "stable"
	FlagStale                   = "stale"
	FlagStartKey                = "startkey"
	FlagStartKeyDocID           = "startkey-docid"
	FlagUpdate                  = "update"
	FlagUpdateSeq               = "update-seq"
	FlagAttsSince               = "atts-since"
	FlagIncludeDeletedConflicts = "deleted-conflicts"
	FlagForceLatest             = "latest"
	FlagIncludeLocalSeq         = "local-seq"
	FlagMeta                    = "meta"
	FlagOpenRevs                = "open-revs"
	FlagRevs                    = "revs"
	FlagRevsInfo                = "revs-info"
	FlagBatch                   = "batch"
	FlagNewEdits                = "new-edits"

	// Curl-equivalent short flags
	FlagShortVerbose    = "v"
	FlagShortOutputFile = "o"
	FlagShortData       = "d"
	FlagShortHead       = "I"
	FlagShortDumpHeader = "D"
	FlagShortUser       = "u"

	// Short versions, custom
	FlagShortServerRoot   = "S"
	FlagShortOutputFormat = "F"
	FlagShortRev          = "r"
	FlagShortAutoRev      = "R"
	FlagShortShards       = "q"
	FlagShortPassword     = "p"
)

type paramParser func(flags *pflag.FlagSet, flagName string) ([]string, error)

var flagParsers = map[string]paramParser{
	FlagEndKey:                  parseParamString,
	FlagEndKeyDocID:             parseParamString,
	FlagKey:                     parseParamString,
	FlagStale:                   parseParamString,
	FlagStartKey:                parseParamString,
	FlagStartKeyDocID:           parseParamString,
	FlagUpdate:                  parseParamString,
	FlagRev:                     parseParamString,
	FlagKeys:                    parseParamStringArray,
	FlagGroupLevel:              parseParamInt,
	FlagLimit:                   parseParamInt,
	FlagSkip:                    parseParamInt,
	FlagShards:                  parseParamInt,
	FlagConflicts:               parseParamBool,
	FlagDescending:              parseParamBool,
	FlagGroup:                   parseParamBool,
	FlagIncludeDocs:             parseParamBool,
	FlagIncludeAttachments:      parseParamBool,
	FlagIncludeAttEncoding:      parseParamBool,
	FlagInclusiveEnd:            parseParamBool,
	FlagReduce:                  parseParamBool,
	FlagSorted:                  parseParamBool,
	FlagStable:                  parseParamBool,
	FlagUpdateSeq:               parseParamBool,
	FlagIncludeDeletedConflicts: parseParamBool,
	FlagAttsSince:               parseParamStringSlice,
	FlagOpenRevs:                parseParamStringSlice,
	FlagForceLatest:             parseParamBool,
	FlagIncludeLocalSeq:         parseParamBool,
	FlagMeta:                    parseParamBool,
	FlagRevs:                    parseParamBool,
	FlagRevsInfo:                parseParamBool,
	FlagNewEdits:                parseParamBool,
}

type paramValidator func(flag string, value []string) error

var flagValidators = map[string]paramValidator{
	FlagStale: func(flag string, v []string) error {
		switch v[0] {
		case "ok", "update_after", "false":
			return nil
		}
		return errors.NewExitError(chttp.ExitFailedToInitialize, "Invalid value for --%s. Supported options: `ok`, `update_after`, `false`", flag)
	},
	FlagUpdate: func(flag string, v []string) error {
		switch v[0] {
		case "true", "false", "lazy":
			return nil
		}
		return errors.NewExitError(chttp.ExitFailedToInitialize, "Invalid value for --%s. Supported options: `true`, `false`, `lazy`", flag)
	},
}
