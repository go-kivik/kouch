package documents

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
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
			expected: &kouch.Options{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					Document: "123",
				},
				Options: &chttp.Options{},
			},
		},
		{
			name: "db included in target",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"/foo/123"},
			expected: &kouch.Options{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "foo",
					Document: "123",
				},
				Options: &chttp.Options{},
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
			expected: &kouch.Options{
				Target: &kouch.Target{
					Root:     "http://foo.com",
					Database: "foo",
					Document: "123",
				},
				Options: &chttp.Options{},
			},
		},
		{
			name: "if-none-match",
			args: []string{"--" + kouch.FlagIfNoneMatch, "foo", "foo.com/bar/baz"},
			expected: &kouch.Options{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					Document: "baz",
				},
				Options: &chttp.Options{IfNoneMatch: "foo"},
			},
		},
		{
			name: "rev",
			args: []string{"--" + kouch.FlagRev, "foo", "foo.com/bar/baz"},
			expected: &kouch.Options{
				Target: &kouch.Target{
					Root:     "foo.com",
					Database: "bar",
					Document: "baz",
				},
				Options: &chttp.Options{
					Query: url.Values{"rev": []string{"foo"}},
				},
			},
		},
		{
			name: "attachments since",
			args: []string{"--" + flagAttsSince, "foo,bar,baz", "docid"},
			expected: &kouch.Options{
				Target: &kouch.Target{Document: "docid"},
				Options: &chttp.Options{
					Query: url.Values{param(flagAttsSince): []string{`["foo","bar","baz"]`}},
				},
			},
		},
		{
			name: "open revs",
			args: []string{"--" + flagOpenRevs, "foo,bar,baz", "docid"},
			expected: &kouch.Options{
				Target: &kouch.Target{Document: "docid"},
				Options: &chttp.Options{
					Query: url.Values{param(flagOpenRevs): []string{`["foo","bar","baz"]`}},
				},
			},
		},
		{
			name: "head",
			args: []string{"--" + kouch.FlagHead, "docid"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Document: "docid"},
				Options: &chttp.Options{},
				Head:    true,
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
			expected: &kouch.Options{
				Target: &kouch.Target{Document: "docid"},
				Options: &chttp.Options{
					Query: url.Values{param(flag): []string{"true"}},
				},
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
		opts     *kouch.Options
		resp     *http.Response
		val      testy.RequestValidator
		expected string
		err      string
		status   int
	}
	tests := []gdTest{
		{
			name:   "validation fails",
			opts:   &kouch.Options{Target: &kouch.Target{}},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "success",
			opts: &kouch.Options{Target: &kouch.Target{Database: "foo", Document: "123"}},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.Path != "/foo/123" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Test\ncontent\n",
		},
		{
			name: "slashes",
			opts: &kouch.Options{Target: &kouch.Target{Database: "foo/ba r", Document: "123/b"}},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb" {
					t.Errorf("Unexpected path: %s", r.URL.RawPath)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Test\ncontent\n",
		},
		{
			name: "if-none-match",
			opts: &kouch.Options{
				Target:  &kouch.Target{Database: "foo", Document: "123"},
				Options: &chttp.Options{IfNoneMatch: "xyz"},
			},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.Path != "/foo/123" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
				if inm := r.Header.Get("If-None-Match"); inm != "\"xyz\"" {
					t.Errorf("Unexpected If-None-Match header: %s", inm)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Test\ncontent\n",
		},
		{
			name: "include query params",
			opts: &kouch.Options{
				Target: &kouch.Target{Database: "foo", Document: "123"},
				Options: &chttp.Options{
					Query: url.Values{"foobar": []string{"baz"}},
				},
			},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.Path != "/foo/123" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
				if val := r.URL.Query().Get("foobar"); val != "baz" {
					t.Errorf("Unexpected query value: %s", val)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Test\ncontent\n",
		},
		{
			name: "head",
			opts: &kouch.Options{
				Target:  &kouch.Target{Database: "foo", Document: "123"},
				Options: &chttp.Options{},
				Head:    true,
			},
			val: func(t *testing.T, r *http.Request) {
				if r.Method != "HEAD" {
					t.Errorf("Unexpected method: %s", r.Method)
				}
				if r.URL.Path != "/foo/123" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"Date": []string{"Mon, 20 Aug 2018 08:28:57 GMT"},
					"ETag": []string{`"2-dcae93de55ac4c27b071654853bca12f"`},
				},
				Body: ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Content-Length: 13\r\n" +
				"Content-Type: text/plain; charset=utf-8\r\n" +
				"Date: Mon, 20 Aug 2018 08:28:57 GMT\r\n" +
				"Etag: \"2-dcae93de55ac4c27b071654853bca12f\"\r\n",
		},
	}
	for _, test := range tests {
		func(test gdTest) {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				if test.resp != nil {
					if test.val != nil {
						s := testy.ServeResponseValidator(t, test.resp, test.val)
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
