package attachments

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name   string
		target *kouch.Target
		err    string
		status int
	}{
		{
			name:   "no filename",
			target: &kouch.Target{},
			err:    "No filename provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no doc id",
			target: &kouch.Target{Filename: "foo.txt"},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no database provided",
			target: &kouch.Target{Document: "123", Filename: "foo.txt"},
			err:    "No database name provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no root url",
			target: &kouch.Target{Database: "foo", Document: "123", Filename: "foo.txt"},
			err:    "No root URL provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "valid",
			target: &kouch.Target{Root: "xxx", Database: "foo", Document: "123", Filename: "foo.txt"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateTarget(test.target)
			testy.ExitStatusError(t, test.err, test.status, err)
		})
	}
}

func TestCommonOpts(t *testing.T) {
	tests := []struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected interface{}
		err      string
		status   int
	}{
		{
			name:   "duplicate filenames",
			args:   []string{"--" + kouch.FlagFilename, "foo.txt", "foo.txt"},
			err:    "Must not use --" + kouch.FlagFilename + " and pass separate filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "id from target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"123/foo.txt", "--database", "bar"},
			expected: &opts{Target: &kouch.Target{
				Root:     "foo.com",
				Database: "bar",
				Document: "123",
				Filename: "foo.txt",
			}},
		},
		{
			name:   "doc ID provided twice",
			args:   []string{"123/foo.txt", "--" + kouch.FlagDocument, "321"},
			err:    "Must not use --id and pass document ID as part of the target",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "db included in target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"/foo/123/foo.txt"},
			expected: &opts{Target: &kouch.Target{
				Root:     "foo.com",
				Database: "foo",
				Document: "123",
				Filename: "foo.txt",
			}},
		},
		{
			name:   "db provided twice",
			args:   []string{"/foo/123/foo.txt", "--" + kouch.FlagDatabase, "foo"},
			err:    "Must not use --" + kouch.FlagDatabase + " and pass database as part of the target",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "full url target",
			args: []string{"http://foo.com/foo/123/foo.txt"},
			expected: &opts{Target: &kouch.Target{
				Root:     "http://foo.com",
				Database: "foo",
				Document: "123",
				Filename: "foo.txt",
			}},
		},
		{
			name: "rev",
			args: []string{"--" + kouch.FlagRev, "xyz", "foo.txt"},
			expected: &opts{
				Target: &kouch.Target{
					Filename: "foo.txt",
				},
				rev: "xyz",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := getAttCmd()
			cmd.ParseFlags(test.args)
			ctx := kouch.GetContext(cmd)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(kouch.SetConf(ctx, test.conf), cmd)
			opts, err := commonOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}
