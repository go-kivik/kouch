package kouch

import (
	"net/url"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/spf13/pflag"
)

func TestSetParam(t *testing.T) {
	type spTest struct {
		flags    *pflag.FlagSet
		flag     string
		expected *Options
		err      string
		status   int
	}
	tests := testy.NewTable()
	tests.Add("no flags set", spTest{
		flags: pflag.NewFlagSet("foo", 1),
		flag:  FlagEndKey,
		expected: &Options{
			Target:  &Target{},
			Options: &chttp.Options{},
		},
	})
	tests.Add("endkey", func() interface{} {
		f := pflag.NewFlagSet("foo", 1)
		f.String(FlagEndKey, "", "x")
		_ = f.Set(FlagEndKey, "oink")
		return spTest{
			flags: f,
			flag:  FlagEndKey,
			expected: &Options{
				Target: &Target{},
				Options: &chttp.Options{
					Query: url.Values{
						"endkey": []string{"oink"},
					},
				},
			},
		}
	})

	tests.Run(t, func(t *testing.T, test spTest) {
		o := NewOptions()
		err := o.SetParam(test.flags, test.flag)
		testy.ExitStatusError(t, test.err, test.status, err)
		if d := diff.Interface(test.expected, o); d != nil {
			t.Error(d)
		}
	})
}

type parseTest struct {
	flags    *pflag.FlagSet
	flag     string
	expected []string
	err      string
	status   int
}

func TestParseParamString(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("flag not set", parseTest{
		flags: pflag.NewFlagSet("foo", 1),
		flag:  "oink",
	})
	tests.Add("wrong flag type", func() interface{} {
		f := pflag.NewFlagSet("foo", 1)
		f.Bool("oink", false, "x")
		return parseTest{
			flags:  f,
			flag:   "oink",
			err:    "trying to get string value of flag of type bool",
			status: 1,
		}
	})
	tests.Add("success", func() interface{} {
		f := pflag.NewFlagSet("foo", 1)
		f.String("oink", "", "x")
		_ = f.Set("oink", "moo")
		return parseTest{
			flags:    f,
			flag:     "oink",
			expected: []string{"moo"},
		}
	})

	tests.Run(t, func(t *testing.T, test parseTest) {
		result, err := parseParamString(test.flags, test.flag)
		testy.ExitStatusError(t, test.err, test.status, err)
		if d := diff.Interface(test.expected, result); d != nil {
			t.Error(d)
		}
	})
}
