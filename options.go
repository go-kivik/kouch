package kouch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/spf13/pflag"
)

// Options represents the accumulated options passed by the user, through config
// files, the commandline, etc.
type Options struct {
	*Target
	*chttp.Options
}

// NewOptions returns a new, empty Options struct.
func NewOptions() *Options {
	return &Options{
		Target:  &Target{},
		Options: &chttp.Options{},
	}
}

// Query returns the url query parameters, initializing it if necessary.
func (o *Options) Query() *url.Values {
	if o.Options.Query == nil {
		o.Options.Query = url.Values{}
	}
	return &o.Options.Query
}

var flagExceptions = map[string]string{
	"shards": "q",
}

func param(flagName string) string {
	if exception, ok := flagExceptions[flagName]; ok {
		return exception
	}
	return strings.Replace(flagName, "-", "_", -1)
}

// SetParamBool sets the query paramater boolean value specified by flagName if
// the provided value differs from the default. This means that values which
// default to true are also supported, but only added to the query when the
// user requests false.
func (o *Options) SetParamBool(f *pflag.FlagSet, flagName string) error {
	v, err := f.GetBool(flagName)
	textV := fmt.Sprintf("%v", v)
	if err == nil && textV != f.Lookup(flagName).DefValue {
		o.Query().Add(param(flagName), textV)
	}
	return err
}

// SetParamStringSlice sets the query param string slice value specified by
// flagName.
func (o *Options) SetParamStringSlice(f *pflag.FlagSet, flagName string) error {
	v, err := f.GetStringSlice(flagName)
	if err == nil && len(v) > 0 {
		enc, e := json.Marshal(v)
		if e != nil {
			return e
		}
		o.Query().Add(param(flagName), string(enc))
	}
	return err
}

// SetParamString sets the query parameter string value specified by flagName,
// if it differs from the default.
func (o *Options) SetParamString(f *pflag.FlagSet, flagName string) error {
	if flag := f.Lookup(flagName); flag == nil {
		return nil
	}
	v, err := f.GetString(flagName)
	if err == nil && v != f.Lookup(flagName).DefValue {
		o.Query().Add(param(flagName), v)
	}
	return err
}

// SetParamInt sets the query parameter string value specified by flagName,
// if it differs from the default.
func (o *Options) SetParamInt(f *pflag.FlagSet, flagName string) error {
	if flag := f.Lookup(flagName); flag == nil {
		return nil
	}
	v, err := f.GetInt(flagName)
	if err == nil && strconv.Itoa(v) != f.Lookup(flagName).DefValue {
		o.Query().Add(param(flagName), strconv.Itoa(v))
	}
	return err
}
