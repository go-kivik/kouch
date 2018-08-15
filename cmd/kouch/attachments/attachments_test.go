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
			expected: &getAttOpts{
				kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					DocID:    "123",
					Filename: "foo.txt",
				},
			},
		},
		{
			name:   "doc ID provided twice",
			args:   []string{"123/foo.txt", "--" + kouch.FlagDocID, "321"},
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
			expected: &getAttOpts{
				kouch.Target{
					Root:     "foo.com",
					Database: "foo",
					DocID:    "123",
					Filename: "foo.txt",
				},
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
				kouch.Target{
					Root:     "http://foo.com/",
					Database: "foo",
					DocID:    "123",
					Filename: "foo.txt",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := attCmd()
			kouch.SetContext(kouch.SetConf(kouch.GetContext(cmd), test.conf), cmd)
			cmd.ParseFlags(test.args)
			opts, err := getAttachmentOpts(cmd, cmd.Flags().Args())
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
			opts:   &getAttOpts{kouch.Target{Filename: "foo.txt"}},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no database provided",
			opts:   &getAttOpts{kouch.Target{DocID: "123", Filename: "foo.txt"}},
			err:    "No database name provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no root url",
			opts:   &getAttOpts{kouch.Target{Database: "foo", DocID: "123", Filename: "foo.txt"}},
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
			opts: &getAttOpts{kouch.Target{Database: "foo", DocID: "123", Filename: "foo.txt"}},
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
		{
			name: "slashes",
			opts: &getAttOpts{kouch.Target{Database: "foo/ba r", DocID: "123/b", Filename: "foo/bar.txt"}},
			val: func(r *http.Request) error {
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb/foo%2Fbar.txt" {
					return errors.Errorf("Unexpected path: %s", r.URL.RawPath)
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
						test.opts.Root = s.URL
					} else {
						s := testy.ServeResponse(test.resp)
						defer s.Close()
						test.opts.Root = s.URL
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
