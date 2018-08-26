package documents

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/put"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestPutDocumentOpts(t *testing.T) {
	type pdoTest struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected interface{}
		err      string
		status   int
	}
	tests := []pdoTest{
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
			name: "full commit",
			args: []string{"http://foo.com/foo/123", "--" + kouch.FlagFullCommit},
			expected: &kouch.Options{
				Target: &kouch.Target{
					Root:     "http://foo.com",
					Database: "foo",
					Document: "123",
				},
				Options: &chttp.Options{
					FullCommit: true,
				},
			},
		},
		{
			name: "batch",
			args: []string{"--" + flagBatch, "docid"},
			expected: &kouch.Options{
				Target: &kouch.Target{Document: "docid"},
				Options: &chttp.Options{
					Query: url.Values{param(flagBatch): []string{"ok"}},
				},
			},
		},
		{
			name: "new edits",
			args: []string{"--" + flagNewEdits + "=false", "docid"},
			expected: &kouch.Options{
				Target: &kouch.Target{Document: "docid"},
				Options: &chttp.Options{
					Query: url.Values{param(flagNewEdits): []string{"false"}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := putDocCmd()
			if err := cmd.ParseFlags(test.args); err != nil {
				t.Fatal(err)
			}
			ctx := kouch.GetContext(cmd)
			ctx = kouch.SetConf(ctx, test.conf)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(ctx, cmd)
			opts, err := putDocumentOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestPutDocCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "No document ID provided",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("create success", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"ok":true,"id":"bar","rev":"1-967a00dff5e02add41819138abb3284d"}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Unexpected Content-Type: %s", ct)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar", "-d", `{"oink":foo}`, "-F", "yaml"},
			Stdout: "id: bar\nok: true\nrev: 1-967a00dff5e02add41819138abb3284d",
		}
	})
	tests.Add("manual rev, dump header", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Date":         []string{"Sun, 26 Aug 2018 17:30:01 GMT"},
				"Content-Type": []string{"application/json"},
			},
			Body: ioutil.NopCloser(strings.NewReader(`{"ok":true,"id":"bar","rev":"2-967a00dff5e02add41819138abb3284d"}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Unexpected Content-Type: %s", ct)
			}
			if rev := r.URL.Query().Get("rev"); rev != "1-967a00dff5e02add41819138abb3284d" {
				t.Errorf("Unexpected rev: %s", rev)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{s.URL + "/foo/bar", "-d", `{"oink":foo}`, "-F", "yaml", "--rev", "1-967a00dff5e02add41819138abb3284d", "--dump-header", "-"},
			Stdout: "Content-Length: 65\r\n" +
				"Content-Type: application/json\r\n" +
				"Date: Sun, 26 Aug 2018 17:30:01 GMT\r\n" +
				"\r\n" +
				"id: bar\nok: true\nrev: 2-967a00dff5e02add41819138abb3284d",
		}
	})
	tests.Add("auto rev", func(t *testing.T) interface{} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if r.Method == http.MethodHead {
				w.Header().Add("ETag", `"1-xyz"`)
				w.WriteHeader(200)
				return
			}
			if rev := r.URL.Query().Get("rev"); rev != "1-xyz" {
				t.Errorf("Unexpected rev: %s", rev)
			}
			w.WriteHeader(200)
			w.Write([]byte(`{"ok":true,"id":"bar","rev":"2-967a00dff5e02add41819138abb3284d"}`))
		}))
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar", "-d", `{"oink":foo}`, "-F", "yaml", "--auto-rev"},
			Stdout: "id: bar\nok: true\nrev: 2-967a00dff5e02add41819138abb3284d",
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"put", "doc"}))
}
