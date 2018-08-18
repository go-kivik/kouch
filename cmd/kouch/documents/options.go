package documents

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/go-kivik/kouch"
	"github.com/spf13/pflag"
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

func (o *opts) setBool(f *pflag.FlagSet, flagName string) error {
	v, err := f.GetBool(flagName)
	if err == nil && v {
		o.Values.Add(param(flagName), "true")
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
