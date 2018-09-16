package kouch

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestTargetScopeName(t *testing.T) {
	for scope := TargetScope(0); scope < targetLastScope+1; scope++ {
		result := TargetScopeName(scope)
		if result == "" {
			t.Errorf("No name defined for scope #%d", scope)
		}
	}
}

// borrowed from attachments
func addCommonFlags(flags *pflag.FlagSet) {
	flags.String(FlagFilename, "", "The attachment filename to fetch. Only necessary if the filename contains slashes, to disambiguate from {id}/{filename}.")
	flags.String(FlagDocument, "", "The document ID. May be provided with the target in the format {id}/{filename}.")
	flags.String(FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}/{filename}")
	flags.StringP(FlagRev, FlagShortRev, "", "Retrieves attachment from document of specified revision.")
}

func TestNewTarget(t *testing.T) {
	defaultConfig := &Config{
		DefaultContext: "foo",
		Contexts:       []NamedContext{{Name: "foo", Context: &Context{Root: "foo.com"}}},
	}

	type newTargetTest struct {
		scope    TargetScope
		addFlags func(*pflag.FlagSet)
		conf     *Config
		args     []string
		expected *Target
		err      string
		status   int
	}
	tests := testy.NewTable()
	tests.Add("no flags", newTargetTest{
		scope: TargetRoot,
		conf:  defaultConfig,
		expected: &Target{
			Root: "foo.com",
		},
	})
	tests.Add("target only", newTargetTest{
		scope:    TargetRoot,
		args:     []string{"http://localhost/"},
		conf:     defaultConfig,
		expected: &Target{Root: "http://localhost/"},
	})
	tests.Add("duplicate filenames", newTargetTest{
		scope:    TargetAttachment,
		addFlags: addCommonFlags,
		args:     []string{"--" + FlagFilename, "foo.txt", "foo.txt"},
		err:      "Must not use --" + FlagFilename + " and pass separate filename",
		status:   chttp.ExitFailedToInitialize,
	})
	tests.Add("id from target", newTargetTest{
		scope:    TargetAttachment,
		addFlags: addCommonFlags,
		conf:     defaultConfig,
		args:     []string{"123/foo.txt", "--database", "bar"},
		expected: &Target{
			Root:     "foo.com",
			Database: "bar",
			Document: "123",
			Filename: "foo.txt",
		},
	})
	tests.Add("doc ID provided twice", newTargetTest{
		scope:    TargetAttachment,
		addFlags: addCommonFlags,
		args:     []string{"123/foo.txt", "--" + FlagDocument, "321"},
		err:      "Must not use --id and pass document ID as part of the target",
		status:   chttp.ExitFailedToInitialize,
	})
	tests.Add("db included in target", newTargetTest{
		scope:    TargetAttachment,
		addFlags: addCommonFlags,
		conf:     defaultConfig,
		args:     []string{"/foo/123/foo.txt"},
		expected: &Target{
			Root:     "foo.com",
			Database: "foo",
			Document: "123",
			Filename: "foo.txt",
		},
	})
	tests.Add("db provided twice", newTargetTest{
		scope:    TargetAttachment,
		addFlags: addCommonFlags,
		args:     []string{"/foo/123/foo.txt", "--" + FlagDatabase, "foo"},
		err:      "Must not use --" + FlagDatabase + " and pass database as part of the target",
		status:   chttp.ExitFailedToInitialize,
	})
	tests.Add("full url target", newTargetTest{
		scope:    TargetAttachment,
		addFlags: addCommonFlags,
		conf:     defaultConfig,
		args:     []string{"http://xyz.com/qrs/123/tuv.txt"},
		expected: &Target{
			Root:     "http://xyz.com",
			Database: "qrs",
			Document: "123",
			Filename: "tuv.txt",
		},
	})

	tests.Run(t, func(t *testing.T, test newTargetTest) {
		cmd := &cobra.Command{}
		if af := test.addFlags; af != nil {
			af(cmd.PersistentFlags())
		}
		if e := cmd.ParseFlags(test.args); e != nil {
			t.Fatal(e)
		}
		ctx := GetContext(cmd)
		ctx = SetConf(ctx, test.conf)
		if flags := cmd.Flags().Args(); len(flags) > 0 {
			ctx = SetTarget(ctx, flags[0])
		}
		target, err := NewTarget(ctx, test.scope, cmd.Flags())
		testy.ExitStatusError(t, test.err, test.status, err)
		if d := diff.Interface(test.expected, target); d != nil {
			t.Error(d)
		}
	})
}

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name     string
		scope    TargetScope
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
			scope:  targetLastScope + 1,
			name:   "invalid scope",
			src:    "xxx",
			err:    "invalid scope",
			status: 1,
		},
		{
			scope:    TargetRoot,
			name:     "blank input",
			src:      "",
			expected: &Target{},
		},
		{
			name:     "Simple root URL",
			scope:    TargetRoot,
			src:      "http://foo.com/",
			expected: &Target{Root: "http://foo.com/"},
		},
		{
			scope:    TargetRoot,
			name:     "url with auth",
			src:      "http://xxx:yyy@foo.com/",
			expected: &Target{Root: "http://foo.com/", Username: "xxx", Password: "yyy"},
		},
		{
			scope:    TargetRoot,
			name:     "Simple root URL with path",
			src:      "http://foo.com/db/",
			expected: &Target{Root: "http://foo.com/db/"},
		},
		{
			scope:    TargetRoot,
			name:     "implicit scheme",
			src:      "foo.com",
			expected: &Target{Root: "foo.com"},
		},
		{
			scope:    TargetRoot,
			name:     "port number",
			src:      "foo.com:5555",
			expected: &Target{Root: "foo.com:5555"},
		},
		{
			scope:  TargetRoot,
			name:   "invalid url",
			src:    "http://foo.com/%xx/",
			err:    `parse http://foo.com/%xx/: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			scope:    TargetDatabase,
			name:     "db only",
			src:      "dbname",
			expected: &Target{Database: "dbname"},
		},
		{
			scope:    TargetDatabase,
			name:     "full url",
			src:      "http://foo.com/dbname",
			expected: &Target{Root: "http://foo.com", Database: "dbname"},
		},
		{
			scope:    TargetDatabase,
			name:     "url with auth",
			src:      "http://a:b@foo.com/dbname",
			expected: &Target{Root: "http://foo.com", Username: "a", Password: "b", Database: "dbname"},
		},
		{
			scope:  TargetDatabase,
			name:   "invalid url",
			src:    "http://foo.com/%xx",
			err:    `parse http://foo.com/%xx: invalid URL escape "%xx"`,
			status: chttp.ExitStatusURLMalformed,
		},
		{
			scope: TargetDatabase,
			name:  "subdir-hosted TargetRoot, with db",
			src:   "https://foo.com/root/dbname",
			expected: &Target{
				Root:     "https://foo.com/root",
				Database: "dbname",
			},
		},
		{
			scope: TargetDatabase,
			name:  "No scheme",
			src:   "example.com:5000/foo",
			expected: &Target{
				Root:     "example.com:5000",
				Database: "foo",
			},
		},
		{
			scope: TargetDatabase,
			name:  "multiple slashes",
			src:   "foo.com/foo/bar/baz",
			expected: &Target{
				Root:     "foo.com/foo/bar",
				Database: "baz",
			},
		},
		{
			scope: TargetDatabase,
			name:  "encoded slash in dbname",
			src:   "foo.com/foo/bar%2Fbaz",
			expected: &Target{
				Root:     "foo.com/foo",
				Database: "bar%2Fbaz",
			},
		},
		{
			scope:  TargetDatabase,
			name:   "missing db",
			src:    "https://foo.com/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    TargetDocument,
			name:     "doc id only",
			src:      "bar",
			expected: &Target{Document: "bar"},
		},
		{
			scope:    TargetDocument,
			name:     "db/docid",
			src:      "foo/bar",
			expected: &Target{Database: "foo", Document: "bar"},
		},
		{
			scope:    TargetDocument,
			name:     "relative design doc",
			src:      "_design/bar",
			expected: &Target{Document: "_design/bar"},
		},
		{
			scope:    TargetDocument,
			name:     "relative local doc",
			src:      "_local/bar",
			expected: &Target{Document: "_local/bar"},
		},
		{
			scope:    TargetDocument,
			name:     "relative design doc with db",
			src:      "foo/_design/bar",
			expected: &Target{Database: "foo", Document: "_design/bar"},
		},
		{
			scope:    TargetDocument,
			name:     "odd chars",
			src:      "foo/foo:bar@baz",
			expected: &Target{Database: "foo", Document: "foo:bar@baz"},
		},
		{
			scope:    TargetDocument,
			name:     "full url",
			src:      "http://localhost:5984/foo/bar",
			expected: &Target{Root: "http://localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:    TargetDocument,
			name:     "url with auth",
			src:      "http://foo:bar@localhost:5984/foo/bar",
			expected: &Target{Root: "http://localhost:5984", Username: "foo", Password: "bar", Database: "foo", Document: "bar"},
		},
		{
			scope:    TargetDocument,
			name:     "no scheme",
			src:      "localhost:5984/foo/bar",
			expected: &Target{Root: "localhost:5984", Database: "foo", Document: "bar"},
		},
		{
			scope:  TargetDocument,
			name:   "url missing doc",
			src:    "http://localhost:5984/foo",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:  TargetDocument,
			name:   "url missing db",
			src:    "http://localhost:5984/",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    TargetAttachment,
			name:     "filename only",
			src:      "baz.txt",
			expected: &Target{Filename: "baz.txt"},
		},
		{
			scope:    TargetAttachment,
			name:     "doc and filename",
			src:      "bar/baz.jpg",
			expected: &Target{Document: "bar", Filename: "baz.jpg"},
		},
		{
			scope:    TargetAttachment,
			name:     "db, doc, filename",
			src:      "foo/bar/baz.png",
			expected: &Target{Database: "foo", Document: "bar", Filename: "baz.png"},
		},
		{
			scope:    TargetAttachment,
			name:     "db, design doc, filename",
			src:      "foo/_design/bar/baz.html",
			expected: &Target{Database: "foo", Document: "_design/bar", Filename: "baz.html"},
		},
		{
			scope:    TargetAttachment,
			name:     "full url",
			src:      "http://host.com/foo/bar/baz.html",
			expected: &Target{Root: "http://host.com", Database: "foo", Document: "bar", Filename: "baz.html"},
		},
		{
			scope:    TargetAttachment,
			name:     "full url, subdir root",
			src:      "http://host.com/couchdb/foo/bar/baz.html",
			expected: &Target{Root: "http://host.com/couchdb", Database: "foo", Document: "bar", Filename: "baz.html"},
		},
		{
			scope:  TargetAttachment,
			name:   "url missing filename",
			src:    "http://host.com/foo/bar",
			err:    "incomplete target URL",
			status: chttp.ExitFailedToInitialize,
		},
		{
			scope:    TargetAttachment,
			name:     "full url, no scheme",
			src:      "foo.com:5984/foo/bar/baz.txt",
			expected: &Target{Root: "foo.com:5984", Database: "foo", Document: "bar", Filename: "baz.txt"},
		},
		{
			scope:    TargetAttachment,
			name:     "url with auth",
			src:      "https://admin:abc123@localhost:5984/foo/bar/baz.pdf",
			expected: &Target{Root: "https://localhost:5984", Username: "admin", Password: "abc123", Database: "foo", Document: "bar", Filename: "baz.pdf"},
		},
		{
			scope:    TargetAttachment,
			name:     "odd chars",
			src:      "dbname/foo:bar@baz/@1:2.txt",
			expected: &Target{Database: "dbname", Document: "foo:bar@baz", Filename: "@1:2.txt"},
		},
		{
			scope:    TargetAttachment,
			name:     "odd chars, filename only",
			src:      "@1:2.txt",
			expected: &Target{Filename: "@1:2.txt"},
		},
	}
	for _, test := range tests {
		scopeName := TargetScopeName(test.scope)
		if scopeName == "" {
			scopeName = "Unknown"
		}
		t.Run(scopeName+"_"+test.name, func(t *testing.T) {
			result, err := ParseTarget(test.scope, test.src)
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, result); d != nil {
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

func TestDocumentFromFlags(t *testing.T) {
	idFlagSet := func() *pflag.FlagSet {
		return flagSet(func(pf *pflag.FlagSet) {
			pf.String(FlagDocument, "", "id")
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
			flags:    idFlagSet(),
			expected: &Target{},
		},
		{
			name:   "id already set",
			target: &Target{Document: "321"},
			flags: func() *pflag.FlagSet {
				fs := idFlagSet()
				if err := fs.Set("id", "123"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			err: "Must not use --" + FlagDocument + " and pass document ID as part of the target",
		},
		{
			name:   "id set anew",
			target: &Target{},
			flags: func() *pflag.FlagSet {
				fs := idFlagSet()
				if err := fs.Set("id", "123"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			expected: &Target{Document: "123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.target.DocumentFromFlags(test.flags)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, test.target); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestDatabaseFromFlags(t *testing.T) {
	dbFlagSet := func() *pflag.FlagSet {
		return flagSet(func(pf *pflag.FlagSet) {
			pf.String(FlagDatabase, "", "db")
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
			flags:    dbFlagSet(),
			expected: &Target{},
		},
		{
			name:   "id already set",
			target: &Target{Database: "foo"},
			flags: func() *pflag.FlagSet {
				fs := dbFlagSet()
				if err := fs.Set("database", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			err: "Must not use --" + FlagDatabase + " and pass database as part of the target",
		},
		{
			name:   "id set anew",
			target: &Target{},
			flags: func() *pflag.FlagSet {
				fs := dbFlagSet()
				if err := fs.Set("database", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			expected: &Target{Database: "bar"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.target.DatabaseFromFlags(test.flags)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, test.target); d != nil {
				t.Error(d)
			}
		})
	}
}
