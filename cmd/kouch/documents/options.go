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

func (o *opts) setBool(f *pflag.FlagSet, flagName, paramName string) error {
	v, err := f.GetBool(flagName)
	if err == nil && v {
		o.Values.Add(paramName, "true")
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

func (o *opts) setIncludeAttachments(f *pflag.FlagSet) error {
	return o.setBool(f, flagIncludeAttachments, paramIncludeAttachments)
}

func (o *opts) setIncludeAttEncoding(f *pflag.FlagSet) error {
	return o.setBool(f, flagIncludeAttEncoding, paramIncludeAttEncoding)
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
	return o.setBool(f, flagIncludeConflicts, paramIncludeConflicts)
}

func (o *opts) setIncludeDeletedConflicts(f *pflag.FlagSet) error {
	return o.setBool(f, flagIncludeDeletedConflicts, paramIncludeDeletedConflicts)
}

func (o *opts) setForceLatest(f *pflag.FlagSet) error {
	return o.setBool(f, flagForceLatest, paramForceLatest)
}
