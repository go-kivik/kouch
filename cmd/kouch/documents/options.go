package documents

import (
	"encoding/json"
	"net/url"

	"github.com/go-kivik/kouch"
	"github.com/spf13/pflag"
)

type opts struct {
	*kouch.Target
	*url.Values
	ifNoneMatch string
}

func newOpts() *opts {
	return &opts{
		Target: &kouch.Target{},
		Values: &url.Values{},
	}
}

func (o *opts) setRev(f *pflag.FlagSet) error {
	v, err := f.GetString(kouch.FlagRev)
	if err == nil && v != "" {
		o.Values.Add("rev", v)
	}
	return err
}

func (o *opts) setIncludeAttachments(f *pflag.FlagSet) error {
	v, err := f.GetBool(flagIncludeAttachments)
	if err == nil && v {
		o.Values.Add(paramIncludeAttachments, "true")
	}
	return err
}

func (o *opts) setIncludeAttEncoding(f *pflag.FlagSet) error {
	v, err := f.GetBool(flagIncludeAttEncoding)
	if err == nil && v {
		o.Values.Add(paramIncludeAttEncoding, "true")
	}
	return err
}

func (o *opts) setAttsSince(f *pflag.FlagSet) error {
	v, err := f.GetStringSlice(flagAttsSince)
	if err == nil && len(v) > 0 {
		enc, e := json.Marshal(v)
		if e != nil {
			return e
		}
		o.Values.Add(paramAttsSince, string(enc))
	}
	return err
}

func (o *opts) setIncludeConflicts(f *pflag.FlagSet) error {
	v, err := f.GetBool(flagIncludeConflicts)
	if err == nil && v {
		o.Values.Add(paramIncludeConflicts, "true")
	}
	return err
}
