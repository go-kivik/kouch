package documents

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/pkg/errors"
)

func TestGetDocumentOpts(t *testing.T) {
	type gdoTest struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected interface{}
		err      string
		status   int
	}
	tests := []gdoTest{
		{
			name:   "duplicate id",
			args:   []string{"--" + kouch.FlagDocument, "foo", "bar"},
			err:    "Must not use --" + kouch.FlagDocument + " and pass document ID as part of the target",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "id from target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"--database", "bar", "123"},
			expected: &opts{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					Document: "123",
				},
				Values: &url.Values{},
			},
		},
		{
			name: "db included in target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"/foo/123"},
			expected: &opts{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "foo",
					Document: "123",
				},
				Values: &url.Values{},
			},
		},
		{
			name:   "db provided twice",
			args:   []string{"/foo/123/foo.txt", "--" + kouch.FlagDatabase, "foo"},
			err:    "Must not use --" + kouch.FlagDatabase + " and pass database as part of the target",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "full url target",
			args: []string{"http://foo.com/foo/123"},
			expected: &opts{
				Target: &kouch.Target{
					Root:     "http://foo.com",
					Database: "foo",
					Document: "123",
				},
				Values: &url.Values{},
			},
		},
		{
			name: "if-none-match",
			args: []string{"--" + kouch.FlagIfNoneMatch, "foo", "foo.com/bar/baz"},
			expected: &opts{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					Document: "baz",
				},
				Values:      &url.Values{},
				ifNoneMatch: "foo"},
		},
		{
			name: "rev",
			args: []string{"--" + kouch.FlagRev, "foo", "foo.com/bar/baz"},
			expected: &opts{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					Document: "baz",
				},
				Values: &url.Values{"rev": []string{"foo"}},
			},
		},
		{
			name: "attachments since",
			args: []string{"--" + flagAttsSince, "foo,bar,baz", "docid"},
			expected: &opts{
				Target: &kouch.Target{Document: "docid"},
				Values: &url.Values{param(flagAttsSince): []string{`["foo","bar","baz"]`}},
			},
		},
		{
			name: "open revs",
			args: []string{"--" + flagOpenRevs, "foo,bar,baz", "docid"},
			expected: &opts{
				Target: &kouch.Target{Document: "docid"},
				Values: &url.Values{param(flagOpenRevs): []string{`["foo","bar","baz"]`}},
			},
		},
	}
	for _, flag := range []string{
		flagIncludeAttachments, flagIncludeAttEncoding, flagIncludeConflicts,
		flagIncludeDeletedConflicts, flagForceLatest, flagIncludeLocalSeq,
		flagMeta, flagRevs, flagRevsInfo,
	} {
		tests = append(tests, gdoTest{
			name: flag,
			args: []string{"--" + flag, "docid"},
			expected: &opts{
				Target: &kouch.Target{Document: "docid"},
				Values: &url.Values{param(flag): []string{"true"}},
			},
		})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := getDocCmd()
			if err := cmd.ParseFlags(test.args); err != nil {
				t.Fatal(err)
			}
			ctx := kouch.GetContext(cmd)
			ctx = kouch.SetConf(ctx, test.conf)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(ctx, cmd)
			opts, err := getDocumentOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name   string
		target *kouch.Target
		err    string
		status int
	}{
		{
			name:   "no doc id",
			target: &kouch.Target{},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no database provided",
			target: &kouch.Target{Document: "123"},
			err:    "No database name provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no root url",
			target: &kouch.Target{Database: "foo", Document: "123"},
			err:    "No root URL provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "valid",
			target: &kouch.Target{Root: "xxx", Database: "foo", Document: "123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateTarget(test.target)
			testy.ExitStatusError(t, test.err, test.status, err)
		})
	}
}

func TestGetDocument(t *testing.T) {
	type gdTest struct {
		name     string
		opts     *opts
		resp     *http.Response
		val      testy.RequestValidator
		expected string
		err      string
		status   int
	}
	tests := []gdTest{
		{
			name:   "validation fails",
			opts:   &opts{Target: &kouch.Target{}, Values: &url.Values{}},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "success",
			opts: &opts{Target: &kouch.Target{Database: "foo", Document: "123"}, Values: &url.Values{}},
			val: func(r *http.Request) error {
				if r.URL.Path != "/foo/123" {
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
			opts: &opts{Target: &kouch.Target{Database: "foo/ba r", Document: "123/b"}, Values: &url.Values{}},
			val: func(r *http.Request) error {
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb" {
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
		{
			name: "if-none-match",
			opts: &opts{Target: &kouch.Target{Database: "foo", Document: "123"}, Values: &url.Values{}, ifNoneMatch: "xyz"},
			val: func(r *http.Request) error {
				if r.URL.Path != "/foo/123" {
					err := errors.Errorf("Unexpected path: %s", r.URL.Path)
					fmt.Println(err)
					return err
				}
				if inm := r.Header.Get("If-None-Match"); inm != "\"xyz\"" {
					err := errors.Errorf("Unexpected If-None-Match header: %s", inm)
					fmt.Println(err)
					return err
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
			name: "include query params",
			opts: &opts{Target: &kouch.Target{Database: "foo", Document: "123"},
				Values: &url.Values{"foobar": []string{"baz"}},
			},
			val: func(r *http.Request) error {
				if r.URL.Path != "/foo/123" {
					err := errors.Errorf("Unexpected path: %s", r.URL.Path)
					fmt.Println(err)
					return err
				}
				if val := r.URL.Query().Get("foobar"); val != "baz" {
					err := errors.Errorf("Unexpected query value: %s", val)
					fmt.Println(err)
					return err
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
		func(test gdTest) {
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
				result, err := getDocument(test.opts)
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
