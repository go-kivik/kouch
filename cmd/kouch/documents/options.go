package documents

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/kouch"
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

type opts struct {
	*kouch.Target
	*url.Values
	ifNoneMatch string
	fullCommit  bool
}

func newOpts() *opts {
	return &opts{
		Target: &kouch.Target{},
		Values: &url.Values{},
	}
}

// setBool sets the query paramater boolean value defined by flagName, if the
// provided value differs from the default. This means that values which
// default to true are also supported, but only added to the query when the
// user requests false.
func (o *opts) setBool(f *pflag.FlagSet, flagName string) error {
	v, err := f.GetBool(flagName)
	textV := fmt.Sprintf("%v", v)
	if err == nil && textV != f.Lookup(flagName).DefValue {
		o.Values.Add(param(flagName), textV)
	}
	return err
}

func (o *opts) setStringSlice(f *pflag.FlagSet, flagName string) error {
	v, err := f.GetStringSlice(flagName)
	if err == nil && len(v) > 0 {
		enc, e := json.Marshal(v)
		if e != nil {
			return e
		}
		o.Values.Add(param(flagName), string(enc))
	}
	return err
}

func (o *opts) setRev(f *pflag.FlagSet) error {
	v, err := f.GetString(kouch.FlagRev)
	if err == nil && v != "" {
		o.Values.Add("rev", v)
	}
	return err
}

func (o *opts) setBatch(f *pflag.FlagSet) error {
	v, err := f.GetBool(flagBatch)
	if err == nil && v {
		o.Values.Add(param(flagBatch), "ok")
	}
	return err
}
