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
			scope:  -1,
			name:   "invalid scope",
			src:    "xxx",
			err:    "invalid scope",
			status: 1,
		},
		{
			scope:  lastScope + 1,
			name:   "invalid scope",
			src:    "xxx",
			err:    "invalid scope",
			status: 1,
		},
		{
			scope:    Root,
			name:     "blank input",
			src:      "",
			expected: &Target{},
		},
		{
			name:     "Simple root URL",
			scope:    Root,
			src:      "http://foo.com/",
			expected: &Target{Root: "http://foo.com/"},
		},
		{
			scope:    Root,
			name:     "url with auth",
			src:      "http://xxx:yyy@foo.com/",
			expected: &Target{Root: "http://foo.com/", Username: "xxx", Password: "yyy"},
		},
		{
			scope:    Root,
			name:     "Simple root URL with path",
			src:      "http://foo.com/db/",
			expected: &Target{Root: "http://foo.com/db/"},
		},
		{
			scope:    Root,
			name:     "implicit scheme",
			src:      "foo.com",
			expected: &Target{Root: "foo.com"},
		},
		{
			scope:    Root,
			name:     "port number",
			src:      "foo.com:5555",
			expected: &Target{Root: "foo.com:5555"},
		},
		{
			scope:  Root,
			name:   "invalid url",
			src:    "http://foo.com/%xx/",
			err:    `parse http://foo.com/%xx/: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			scope:    Database,
			name:     "db only",
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
			scope:    Database,
			name:     "url with auth",
			src:      "http://a:b@foo.com/dbname",
			expected: &Target{Root: "http://foo.com", Username: "a", Password: "b", Database: "dbname"},
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
			scope:  Database,
			name:   "missing db",
			src:    "https://foo.com/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
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
			name:     "odd chars",
			src:      "foo/foo:bar@baz",
			expected: &Target{Database: "foo", Document: "foo:bar@baz"},
		},
		{
			scope:    Document,
			name:     "full url",
			src:      "http://localhost:5984/foo/bar",
			expected: &Target{Root: "http://localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:    Document,
			name:     "url with auth",
			src:      "http://foo:bar@localhost:5984/foo/bar",
			expected: &Target{Root: "http://localhost:5984", Username: "foo", Password: "bar", Database: "foo", Document: "bar"},
		},
		{
			scope:    Document,
			name:     "no scheme",
			src:      "localhost:5984/foo/bar",
			expected: &Target{Root: "localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:  Document,
			name:   "url missing doc",
			src:    "http://localhost:5984/foo",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:  Document,
			name:   "url missing db",
			src:    "http://localhost:5984/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    Attachment,
			name:     "filename only",
			src:      "baz.txt",
			expected: &Target{Filename: "baz.txt"},
		},
		{
			scope:    Attachment,
			name:     "doc and filename",
			src:      "bar/baz.jpg",
			expected: &Target{Document: "bar", Filename: "baz.jpg"},
		},
		{
			scope:    Attachment,
			name:     "db, doc, filename",
			src:      "foo/bar/baz.png",
			expected: &Target{Database: "foo", Document: "bar", Filename: "baz.png"},
		},
		{
			scope:    Attachment,
			name:     "db, design doc, filename",
			src:      "foo/_design/bar/baz.html",
			expected: &Target{Database: "foo", Document: "_design/bar", Filename: "baz.html"},
		},
		{
			scope:    Attachment,
			name:     "full url",
			src:      "http://host.com/foo/bar/baz.html",
			expected: &Target{Root: "http://host.com", Database: "foo", Document: "bar", Filename: "baz.html"},
		},
		{
			scope:    Attachment,
			name:     "full url, subdir root",
			src:      "http://host.com/couchdb/foo/bar/baz.html",
			expected: &Target{Root: "http://host.com/couchdb", Database: "foo", Document: "bar", Filename: "baz.html"},
		},
		{
			scope:  Attachment,
			name:   "url missing filename",
			src:    "http://host.com/foo/bar",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    Attachment,
			name:     "full url, no scheme",
			src:      "foo.com:5984/foo/bar/baz.txt",
			expected: &Target{Root: "foo.com:5984", Database: "foo", Document: "bar", Filename: "baz.txt"},
		},
		{
			scope:    Attachment,
			name:     "url with auth",
			src:      "https://admin:abc123@localhost:5984/foo/bar/baz.pdf",
			expected: &Target{Root: "https://localhost:5984", Username: "admin", Password: "abc123", Database: "foo", Document: "bar", Filename: "baz.pdf"},
		},
		{
			scope:    Attachment,
			name:     "odd chars",
			src:      "dbname/foo:bar@baz/@1:2.txt",
			expected: &Target{Database: "dbname", Document: "foo:bar@baz", Filename: "@1:2.txt"},
		},
		{
			scope:    Attachment,
			name:     "odd chars, filename only",
			src:      "@1:2.txt",
			expected: &Target{Filename: "@1:2.txt"},
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
