package attachments

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

/*
func TestGetAttachments(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		resp    *http.Response
		content string
		err     string
		status  int
	}{
		{
			name:   "no filename",
			err:    "Must provide exactly one filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "too many filenames",
			args:   []string{"foo.txt", "bar.jpg"},
			err:    "Must provide exactly one filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "duplicate filenames",
			args:   []string{"--" + FlagFilename, "foo.txt", "foo.txt"},
			err:    "Must use --" + FlagFilename + " and pass separate filename",
			status: chttp.ExitFailedToInitialize,
		},
		// {
		// 	name: "no doc id provided",
		// 	args: []string{"foo.txt"},
		// 	err:  "document id must be provided as part of the filename argument or with the --" + FlagDocID + " flag",
		// },
		// {
		//     name: "Filename with slash",
		//     args: []string{"--filename", "foo/bar.txt"},
		//     err: "xxx",
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := &testy.ResponseHandler{Response: test.resp}
				h.ServeHTTP(w, r)
			}))
			conf := &kouch.Config{
				Contexts: []kouch.NamedContext{
					{
						Name:    "default",
						Context: &kouch.Context{Root: s.URL},
					},
				},
				DefaultContext: "default",
			}
			buf := &bytes.Buffer{}
			cx := &kouch.CmdContext{
				Conf:   conf,
				Output: buf,
			}
			root := &cobra.Command{}
			get := &cobra.Command{Use: "get"}
			get.AddCommand(attCmd(cx))
			root.AddCommand(get)
			root.ParseFlags(test.args)
			attCmd := attachmentCmd(cx)
			err := attCmd(root, nil)
			testy.ExitStatusErrorRE(t, test.err, test.status, err)
			if d := diff.Text(test.content, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}
*/

func TestGetAttachmentOpts(t *testing.T) {
	tests := []struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected interface{}
		err      string
		status   int
	}{
		{
			name:   "no filename",
			args:   nil,
			err:    "No filename provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "too many filenames",
			args:   []string{"foo.txt", "bar.jpg"},
			err:    "Too many targets provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "duplicate filenames",
			args:   []string{"--" + FlagFilename, "foo.txt", "foo.txt"},
			err:    "Must not use --" + FlagFilename + " and pass separate filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "id from target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"123/foo.txt", "--database", "bar"},
			expected: &getAttOpts{
				root:     "foo.com",
				db:       "bar",
				id:       "123",
				filename: "foo.txt",
			},
		},
		{
			name:   "doc ID provided twice",
			args:   []string{"123/foo.txt", "--" + FlagDocID, "321"},
			err:    "Must not use --id and pass doc ID as part of the target",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no doc ID provided",
			args:   []string{"foo.txt"},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "db included in target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"/foo/123/foo.txt"},
			expected: &getAttOpts{
				root:     "foo.com",
				db:       "foo",
				id:       "123",
				filename: "foo.txt",
			},
		},
		{
			name:   "db provided twice",
			args:   []string{"/foo/123/foo.txt", "--" + FlagDatabase, "foo"},
			err:    "Must not use --" + FlagDatabase + " and pass database as part of the target",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "full url target",
			args: []string{"http://foo.com/foo/123/foo.txt"},
			expected: &getAttOpts{
				root:     "http://foo.com/",
				db:       "foo",
				id:       "123",
				filename: "foo.txt",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cx := &attCmdCtx{&kouch.CmdContext{
				Conf: test.conf,
			}}
			cmd := attCmd(cx.CmdContext)
			cmd.ParseFlags(test.args)
			opts, err := cx.getAttachmentOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected *getAttOpts
		err      string
		status   int
	}{
		{
			name:     "simple filename only",
			target:   "foo.txt",
			expected: &getAttOpts{filename: "foo.txt"},
		},
		{
			name:     "simple id/filename",
			target:   "123/foo.txt",
			expected: &getAttOpts{id: "123", filename: "foo.txt"},
		},
		{
			name:     "simple /db/id/filename",
			target:   "/foo/123/foo.txt",
			expected: &getAttOpts{db: "foo", id: "123", filename: "foo.txt"},
		},
		{
			name:     "id + filename with slash",
			target:   "123/foo/bar.txt",
			expected: &getAttOpts{id: "123", filename: "foo/bar.txt"},
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
			expected: &getAttOpts{root: "http://foo.com/", db: "foo", id: "123", filename: "foo.txt"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts, err := parseTarget(test.target)
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestGetAttOpts_Validate(t *testing.T) {
	tests := []struct {
		name   string
		opts   *getAttOpts
		err    string
		status int
	}{
		{
			name:   "no filename",
			opts:   &getAttOpts{},
			err:    "No filename provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no doc id",
			opts:   &getAttOpts{filename: "foo.txt"},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no database provided",
			opts:   &getAttOpts{id: "123", filename: "foo.txt"},
			err:    "No database name provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no root url",
			opts:   &getAttOpts{db: "foo", id: "123", filename: "foo.txt"},
			err:    "No root URL provided",
			status: chttp.ExitFailedToInitialize,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.opts.validate()
			testy.ExitStatusError(t, test.err, test.status, err)
		})
	}
}
