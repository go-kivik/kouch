package documents

import (
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

// Get-doc specific flags
const (
	flagIncludeAttachments      = "attachments"
	flagIncludeAttEncoding      = "att-encoding-info"
	flagAttsSince               = "atts-since"
	flagIncludeConflicts        = "conflicts"
	flagIncludeDeletedConflicts = "deleted-conflicts"
	flagForceLatest             = "latest"
	flagIncludeLocalSeq         = "local-seq"
	flagMeta                    = "meta"
	flagOpenRevs                = "open-revs"
	flagRevs                    = "revs"
	flagRevsInfo                = "revs-info"
)

// Put-doc specific flags
const (
	flagBatch    = "batch"
	flagNewEdits = "new-edits"
)

func param(flagName string) string {
	return strings.Replace(flagName, "-", "_", -1)
}

func setBatch(o *kouch.Options, f *pflag.FlagSet) error {
	v, err := f.GetBool(flagBatch)
	if err == nil && v {
		o.Query().Add(param(flagBatch), "ok")
	}
	return err
}

func validateTarget(t *kouch.Target) error {
	if t.Filename != "" {
		panic("non-nil filename")
	}
	if t.Document == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No document ID provided")
	}
	if t.Database == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No database name provided")
	}
	if t.Root == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No root URL provided")
	}
	return nil
}
