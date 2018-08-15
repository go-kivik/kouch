package kouch

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/spf13/pflag"
)

func TestParseAttachmentTarget(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected *Target
		err      string
		status   int
	}{
		{
			name:     "simple filename only",
			target:   "foo.txt",
			expected: &Target{Filename: "foo.txt"},
		},
		{
			name:     "simple id/filename",
			target:   "123/foo.txt",
			expected: &Target{DocID: "123", Filename: "foo.txt"},
		},
		{
			name:     "simple /db/id/filename",
			target:   "/foo/123/foo.txt",
			expected: &Target{Database: "foo", DocID: "123", Filename: "foo.txt"},
		},
		{
			name:     "id + filename with slash",
			target:   "123/foo/bar.txt",
			expected: &Target{DocID: "123", Filename: "foo/bar.txt"},
		},
		{
			name:   "invalid url",
			target: "http://foo.com/%xx",
			err:    `parse http://foo.com/%xx: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			name:     "full url",
			target:   "http://foo.com/foo/123/foo.txt",
			expected: &Target{Root: "http://foo.com/", Database: "foo", DocID: "123", Filename: "foo.txt"},
		},
		{
			name:   "db, missing filename",
			target: "/db/123",
			err:    "invalid target",
			status: chttp.ExitStatusURLMalformed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target, err := ParseAttachmentTarget(test.target)
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, target); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestFilenameFromFlags(t *testing.T) {
	filenameFlagSet := func() *pflag.FlagSet {
		return flagSet(func(pf *pflag.FlagSet) {
			pf.String(FlagFilename, "", "filename")
		})
	}
	tests := []struct {
		name     string
		target   *Target
		flags    *pflag.FlagSet
		expected *Target
		err      string
	}{
		{
			name:     "no flags",
			target:   &Target{},
			flags:    filenameFlagSet(),
			expected: &Target{},
		},
		{
			name:   "no flag defined",
			target: &Target{},
			flags:  flagSet(),
			err:    "flag accessed but not defined: filename",
		},
		{
			name:   "filename already set",
			target: &Target{Filename: "foo"},
			flags: func() *pflag.FlagSet {
				fs := filenameFlagSet()
				if err := fs.Set("filename", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			err: "Must not use --" + FlagFilename + " and pass separate filename",
		},
		{
			name:   "filename set anew",
			target: &Target{},
			flags: func() *pflag.FlagSet {
				fs := filenameFlagSet()
				if err := fs.Set("filename", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			expected: &Target{Filename: "bar"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.target.FilenameFromFlags(test.flags)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, test.target); d != nil {
				t.Error(d)
			}
		})
	}
}
