package attachments

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/pkg/errors"
)

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
		{
			name:   "db, missing filename",
			target: "/db/123",
			err:    "invalid target",
			status: chttp.ExitStatusURLMalformed,
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

func TestGetAttachment(t *testing.T) {
	type gaTest struct {
		name     string
		opts     *getAttOpts
		resp     *http.Response
		val      testy.RequestValidator
		expected string
		err      string
		status   int
	}
	tests := []gaTest{
		{
			name:   "validation fails",
			opts:   &getAttOpts{},
			err:    "No filename provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "success",
			opts: &getAttOpts{db: "foo", id: "123", filename: "foo.txt"},
			val: func(r *http.Request) error {
				if r.URL.Path != "/foo/123/foo.txt" {
					return errors.Errorf("Unexpected path: %s", r.URL.Path)
				}
				return nil
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Test\ncontent\n",
		},
	}
	for _, test := range tests {
		func(test gaTest) {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				if test.resp != nil {
					if test.val != nil {
						s := testy.ServeResponseValidator(test.resp, test.val)
						defer s.Close()
						test.opts.root = s.URL
					} else {
						s := testy.ServeResponse(test.resp)
						defer s.Close()
						test.opts.root = s.URL
					}
				}
				result, err := getAttachment(test.opts)
				testy.ExitStatusError(t, test.err, test.status, err)
				defer result.Close()
				content, err := ioutil.ReadAll(result)
				if err != nil {
					t.Fatal(err)
				}
				if d := diff.Text(test.expected, string(content)); d != nil {
					t.Error(d)
				}
			})
		}(test)
	}
}
