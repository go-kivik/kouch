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

// SetParams sets parameters based on the provided flags
func (o *Options) SetParams(f *pflag.FlagSet, flags ...string) error {
	for _, flag := range flags {
		if err := o.SetParam(f, flag); err != nil {
			return err
		}
	}
	return nil
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

func parseParamBool(f *pflag.FlagSet, flag string) ([]string, error) {
	if flag := f.Lookup(flag); flag == nil {
		return nil, nil
	}
	v, err := f.GetBool(flag)
	textV := fmt.Sprintf("%v", v)
	if err == nil && textV != f.Lookup(flag).DefValue {
		return []string{textV}, nil
	}
	return nil, err
}

func parseParamStringSlice(f *pflag.FlagSet, flagName string) ([]string, error) {
	v, err := f.GetStringSlice(flagName)
	if err == nil && len(v) > 0 {
		enc, e := json.Marshal(v)
		if e != nil {
			return nil, e
		}
		return []string{string(enc)}, nil
	}
	return nil, err

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

func parseParamInt(f *pflag.FlagSet, flag string) ([]string, error) {
	if flag := f.Lookup(flag); flag == nil {
		return nil, nil
	}
	v, err := f.GetInt(flag)
	if err == nil && strconv.Itoa(v) != f.Lookup(flag).DefValue {
		return []string{strconv.Itoa(v)}, nil
	}
	return nil, err
}
