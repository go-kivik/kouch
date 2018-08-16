package target

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		scope    Scope
		src      string
		expected *Target
		err      string
		status   int
	}{
		{
			name:     "blank input",
			scope:    Root,
			src:      "",
			expected: &Target{},
		},
		{
			name:   "invalid scope",
			src:    "xxx",
			scope:  -1,
			err:    "invalid scope",
			status: 1,
		},
		{
			name:   "invalid scope",
			src:    "xxx",
			scope:  lastScope + 1,
			err:    "invalid scope",
			status: 1,
		},
		{
			name:     "Simple root URL",
			scope:    Root,
			src:      "http://foo.com/",
			expected: &Target{Root: "http://foo.com/"},
		},
		{
			name:     "Simple root URL with path",
			scope:    Root,
			src:      "http://foo.com/db/",
			expected: &Target{Root: "http://foo.com/db/"},
		},
		{
			name:     "implicit scheme",
			scope:    Root,
			src:      "foo.com",
			expected: &Target{Root: "foo.com"},
		},
		{
			name:     "port number",
			scope:    Root,
			src:      "foo.com:5555",
			expected: &Target{Root: "foo.com:5555"},
		},
		{
			name:     "db only",
			scope:    Database,
			src:      "dbname",
			expected: &Target{Database: "dbname"},
		},
		{
			scope:    Database,
			name:     "full url",
			src:      "http://foo.com/dbname",
			expected: &Target{Root: "http://foo.com", Database: "dbname"},
		},
		{
			scope:  Database,
			name:   "invalid url",
			src:    "http://foo.com/%xx",
			err:    `parse http://foo.com/%xx: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			scope: Database,
			name:  "subdir-hosted root, with db",
			src:   "https://foo.com/root/dbname",
			expected: &Target{
				Root:     "https://foo.com/root",
				Database: "dbname",
			},
		},
		{
			scope: Database,
			name:  "No scheme",
			src:   "example.com:5000/foo",
			expected: &Target{
				Root:     "example.com:5000",
				Database: "foo",
			},
		},
		{
			scope: Database,
			name:  "multiple slashes",
			src:   "foo.com/foo/bar/baz",
			expected: &Target{
				Root:     "foo.com/foo/bar",
				Database: "baz",
			},
		},
		{
			scope: Database,
			name:  "encoded slash in dbname",
			src:   "foo.com/foo/bar%2Fbaz",
			expected: &Target{
				Root:     "foo.com/foo",
				Database: "bar%2Fbaz",
			},
		},
		{
			scope:    Document,
			name:     "doc id only",
			src:      "bar",
			expected: &Target{Document: "bar"},
		},
		{
			scope:    Document,
			name:     "db/docid",
			src:      "foo/bar",
			expected: &Target{Database: "foo", Document: "bar"},
		},
		{
			scope:    Document,
			name:     "relative design doc",
			src:      "_design/bar",
			expected: &Target{Document: "_design/bar"},
		},
		{
			scope:    Document,
			name:     "relative local doc",
			src:      "_local/bar",
			expected: &Target{Document: "_local/bar"},
		},
		{
			scope:    Document,
			name:     "relative design doc with db",
			src:      "foo/_design/bar",
			expected: &Target{Database: "foo", Document: "_design/bar"},
		},
		{
			scope:    Document,
			name:     "full url",
			src:      "http://localhost:5984/foo/bar",
			expected: &Target{Root: "http://localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:  Document,
			name:   "incomplete full url",
			src:    "http://localhost:5984/foo",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:  Document,
			name:   "incomplete full url",
			src:    "http://localhost:5984/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
	}
	for _, test := range tests {
		scopeName := ScopeName(test.scope)
		if scopeName == "" {
			scopeName = "Unknown"
		}
		t.Run(scopeName+"_"+test.name, func(t *testing.T) {
			result, err := Parse(test.scope, test.src)
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}
