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
	"github.com/go-kivik/kouch/cmd/kouch/registry"

	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/put"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
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
			opts, err := getDocumentOpts(cmd)
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestGetDocumentCmd(t *testing.T) {
	type gdcTest struct {
		conf   *kouch.Config
		args   []string
		stdout string
		stderr string
		err    string
		status int
	}
	tests := testy.NewTable()
	tests.Add("validation fails", gdcTest{
		args:   []string{},
		err:    "No document ID provided",
		status: chttp.ExitFailedToInitialize,
	})
	tests.Add("success", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return gdcTest{
			args:   []string{s.URL + "/foo/bar"},
			stdout: `{"foo":123}`,
		}
	})
	tests.Add("slashes", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.URL.RawPath != "/foo%2Fba+r/123%2Fb" {
				t.Errorf("Unexpected req path: %s", r.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return gdcTest{
			args: []string{
				"--" + kouch.FlagServerRoot, s.URL,
				"--" + kouch.FlagDatabase, "foo/ba r",
				"--" + kouch.FlagDocument, "123/b",
			},
			stdout: `{"foo":123}`,
		}
	})
	tests.Add("if-none-match", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
			if inm := req.Header.Get("If-None-Match"); inm != "\"oink\"" {
				t.Errorf("Unexpected If-None-Match header: %s", inm)
			}
		})
		tests.Cleanup(s.Close)
		return gdcTest{
			args:   []string{"--" + kouch.FlagIfNoneMatch, "oink", s.URL + "/foo/bar"},
			stdout: `{"foo":123}`,
		}
	})
	tests.Add("rev", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
			if rev := req.URL.Query().Get("rev"); rev != "oink" {
				t.Errorf("Unexpected rev: %s", rev)
			}
		})
		tests.Cleanup(s.Close)
		return gdcTest{
			args:   []string{"--" + kouch.FlagRev, "oink", s.URL + "/foo/bar"},
			stdout: `{"foo":123}`,
		}
	})
	tests.Add("head", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"Date":         []string{"Mon, 20 Aug 2018 08:55:52 GMT"},
			},
			Body: ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.Method != http.MethodHead {
				t.Errorf("Unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/foo/bar/baz.txt" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return gdcTest{
			args: []string{"--" + kouch.FlagHead, s.URL + "/foo/bar/baz.txt"},
			stdout: "Content-Length: 11\r\n" +
				"Content-Type: application/json\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n",
		}
	})

	tests.Run(t, func(t *testing.T, test gdcTest) {
		var err error
		stdout, stderr := testy.RedirIO(nil, func() {
			root := registry.Root()
			root.SetArgs(append([]string{"get", "doc"}, test.args...))
			err = root.Execute()
		})
		if d := diff.Text(test.stdout, stdout); d != nil {
			t.Errorf("STDOUT:\n%s", d)
		}
		if d := diff.Text(test.stderr, stderr); d != nil {
			t.Errorf("STDERR:\n%s", d)
		}
		testy.ExitStatusError(t, test.err, test.status, err)
	})
}
