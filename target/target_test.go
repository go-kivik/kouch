package target

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		scope    kouch.TargetScope
		src      string
		expected *kouch.Target
		err      string
		status   int
	}{
		{
			scope:  -1,
			name:   "invalid scope",
			src:    "xxx",
			err:    "invalid scope",
			status: 1,
		},
		{
			scope:  kouch.TargetLastScope + 1,
			name:   "invalid scope",
			src:    "xxx",
			err:    "invalid scope",
			status: 1,
		},
		{
			scope:    kouch.TargetRoot,
			name:     "blank input",
			src:      "",
			expected: &kouch.Target{},
		},
		{
			name:     "Simple root URL",
			scope:    kouch.TargetRoot,
			src:      "http://foo.com/",
			expected: &kouch.Target{Root: "http://foo.com/"},
		},
		{
			scope:    kouch.TargetRoot,
			name:     "url with auth",
			src:      "http://xxx:yyy@foo.com/",
			expected: &kouch.Target{Root: "http://foo.com/", Username: "xxx", Password: "yyy"},
		},
		{
			scope:    kouch.TargetRoot,
			name:     "Simple root URL with path",
			src:      "http://foo.com/db/",
			expected: &kouch.Target{Root: "http://foo.com/db/"},
		},
		{
			scope:    kouch.TargetRoot,
			name:     "implicit scheme",
			src:      "foo.com",
			expected: &kouch.Target{Root: "foo.com"},
		},
		{
			scope:    kouch.TargetRoot,
			name:     "port number",
			src:      "foo.com:5555",
			expected: &kouch.Target{Root: "foo.com:5555"},
		},
		{
			scope:  kouch.TargetRoot,
			name:   "invalid url",
			src:    "http://foo.com/%xx/",
			err:    `parse http://foo.com/%xx/: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			scope:    kouch.TargetDatabase,
			name:     "db only",
			src:      "dbname",
			expected: &kouch.Target{Database: "dbname"},
		},
		{
			scope:    kouch.TargetDatabase,
			name:     "full url",
			src:      "http://foo.com/dbname",
			expected: &kouch.Target{Root: "http://foo.com", Database: "dbname"},
		},
		{
			scope:    kouch.TargetDatabase,
			name:     "url with auth",
			src:      "http://a:b@foo.com/dbname",
			expected: &kouch.Target{Root: "http://foo.com", Username: "a", Password: "b", Database: "dbname"},
		},
		{
			scope:  kouch.TargetDatabase,
			name:   "invalid url",
			src:    "http://foo.com/%xx",
			err:    `parse http://foo.com/%xx: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			scope: kouch.TargetDatabase,
			name:  "subdir-hosted kouch.TargetRoot, with db",
			src:   "https://foo.com/root/dbname",
			expected: &kouch.Target{
				Root:     "https://foo.com/root",
				Database: "dbname",
			},
		},
		{
			scope: kouch.TargetDatabase,
			name:  "No scheme",
			src:   "example.com:5000/foo",
			expected: &kouch.Target{
				Root:     "example.com:5000",
				Database: "foo",
			},
		},
		{
			scope: kouch.TargetDatabase,
			name:  "multiple slashes",
			src:   "foo.com/foo/bar/baz",
			expected: &kouch.Target{
				Root:     "foo.com/foo/bar",
				Database: "baz",
			},
		},
		{
			scope: kouch.TargetDatabase,
			name:  "encoded slash in dbname",
			src:   "foo.com/foo/bar%2Fbaz",
			expected: &kouch.Target{
				Root:     "foo.com/foo",
				Database: "bar%2Fbaz",
			},
		},
		{
			scope:  kouch.TargetDatabase,
			name:   "missing db",
			src:    "https://foo.com/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    kouch.TargetDocument,
			name:     "doc id only",
			src:      "bar",
			expected: &kouch.Target{Document: "bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "db/docid",
			src:      "foo/bar",
			expected: &kouch.Target{Database: "foo", Document: "bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "relative design doc",
			src:      "_design/bar",
			expected: &kouch.Target{Document: "_design/bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "relative local doc",
			src:      "_local/bar",
			expected: &kouch.Target{Document: "_local/bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "relative design doc with db",
			src:      "foo/_design/bar",
			expected: &kouch.Target{Database: "foo", Document: "_design/bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "odd chars",
			src:      "foo/foo:bar@baz",
			expected: &kouch.Target{Database: "foo", Document: "foo:bar@baz"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "full url",
			src:      "http://localhost:5984/foo/bar",
			expected: &kouch.Target{Root: "http://localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "url with auth",
			src:      "http://foo:bar@localhost:5984/foo/bar",
			expected: &kouch.Target{Root: "http://localhost:5984", Username: "foo", Password: "bar", Database: "foo", Document: "bar"},
		},
		{
			scope:    kouch.TargetDocument,
			name:     "no scheme",
			src:      "localhost:5984/foo/bar",
			expected: &kouch.Target{Root: "localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:  kouch.TargetDocument,
			name:   "url missing doc",
			src:    "http://localhost:5984/foo",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:  kouch.TargetDocument,
			name:   "url missing db",
			src:    "http://localhost:5984/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "filename only",
			src:      "baz.txt",
			expected: &kouch.Target{Filename: "baz.txt"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "doc and filename",
			src:      "bar/baz.jpg",
			expected: &kouch.Target{Document: "bar", Filename: "baz.jpg"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "db, doc, filename",
			src:      "foo/bar/baz.png",
			expected: &kouch.Target{Database: "foo", Document: "bar", Filename: "baz.png"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "db, design doc, filename",
			src:      "foo/_design/bar/baz.html",
			expected: &kouch.Target{Database: "foo", Document: "_design/bar", Filename: "baz.html"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "full url",
			src:      "http://host.com/foo/bar/baz.html",
			expected: &kouch.Target{Root: "http://host.com", Database: "foo", Document: "bar", Filename: "baz.html"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "full url, subdir root",
			src:      "http://host.com/couchdb/foo/bar/baz.html",
			expected: &kouch.Target{Root: "http://host.com/couchdb", Database: "foo", Document: "bar", Filename: "baz.html"},
		},
		{
			scope:  kouch.TargetAttachment,
			name:   "url missing filename",
			src:    "http://host.com/foo/bar",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "full url, no scheme",
			src:      "foo.com:5984/foo/bar/baz.txt",
			expected: &kouch.Target{Root: "foo.com:5984", Database: "foo", Document: "bar", Filename: "baz.txt"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "url with auth",
			src:      "https://admin:abc123@localhost:5984/foo/bar/baz.pdf",
			expected: &kouch.Target{Root: "https://localhost:5984", Username: "admin", Password: "abc123", Database: "foo", Document: "bar", Filename: "baz.pdf"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "odd chars",
			src:      "dbname/foo:bar@baz/@1:2.txt",
			expected: &kouch.Target{Database: "dbname", Document: "foo:bar@baz", Filename: "@1:2.txt"},
		},
		{
			scope:    kouch.TargetAttachment,
			name:     "odd chars, filename only",
			src:      "@1:2.txt",
			expected: &kouch.Target{Filename: "@1:2.txt"},
		},
	}
	for _, test := range tests {
		scopeName := kouch.TargetScopeName(test.scope)
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
