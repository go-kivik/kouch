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

// SetParam sets the named parameter
func (o *Options) SetParam(f *pflag.FlagSet, flag string) error {
	parser, ok := flagParsers[flag]
	if !ok {
		panic(fmt.Sprintf("No setter for %s flag", flag))
	}
	values, err := parser(f, flag)
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	if validator, ok := flagValidators[flag]; ok {
		if e := validator(flag, values); e != nil {
			return e
		}
	}
	paramName := param(flag)
	for _, value := range values {
		o.Query().Add(paramName, value)
	}
	return nil
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

func parseParamStringArray(f *pflag.FlagSet, flag string) ([]string, error) {
	return f.GetStringArray(flag)
}

func parseParamString(f *pflag.FlagSet, flag string) ([]string, error) {
	if flag := f.Lookup(flag); flag == nil {
		return nil, nil
	}
	v, err := f.GetString(flag)
	if err == nil && v != f.Lookup(flag).DefValue {
		return []string{v}, nil
	}
	return nil, err
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
