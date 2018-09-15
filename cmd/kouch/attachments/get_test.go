package attachments

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
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
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
			name: "if none match",
			args: []string{"--" + kouch.FlagIfNoneMatch, "xyz", "foo.txt"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Filename: "foo.txt"},
				Options: &chttp.Options{IfNoneMatch: "xyz"},
			},
		},
		{
			name: "rev",
			args: []string{"--" + kouch.FlagRev, "xyz", "foo.txt"},
			expected: &kouch.Options{
				Target: &kouch.Target{Filename: "foo.txt"},
				Options: &chttp.Options{
					Query: url.Values{"rev": []string{"xyz"}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := getAttCmd()
			if e := cmd.ParseFlags(test.args); e != nil {
				t.Fatal(e)
			}
			ctx := kouch.GetContext(cmd)
			ctx = kouch.SetConf(ctx, test.conf)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(kouch.SetConf(ctx, test.conf), cmd)
			opts, err := getAttachmentOpts(ctx, cmd.Flags())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestGetAttachmentCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "No filename provided",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("slashes", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("attachment content")),
		}, func(t *testing.T, r *http.Request) {
			if r.URL.RawPath != "/foo%2Fba+r/123%2Fb/foo%2Fbar.txt" {
				t.Errorf("Unexpected req path: %s", r.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{
				"--" + kouch.FlagServerRoot, s.URL,
				"--" + kouch.FlagDatabase, "foo/ba r",
				"--" + kouch.FlagDocument, "123/b",
				"--" + kouch.FlagFilename, "foo/bar.txt",
			},
			Stdout: "attachment content\n",
		}
	})
	tests.Add("success", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("attachment content")),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar/foo.txt" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar/foo.txt"},
			Stdout: "attachment content",
		}
	})
	tests.Add("if-none-match", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("attachment content")),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar/baz.txt" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
			if inm := req.Header.Get("If-None-Match"); inm != "\"oink\"" {
				t.Errorf("Unexpected If-None-Match header: %s", inm)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{"--" + kouch.FlagIfNoneMatch, "oink", s.URL + "/foo/bar/baz.txt"},
			Stdout: "attachment content",
		}
	})
	tests.Add("rev", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("attachment content")),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar/baz.txt" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
			if rev := req.URL.Query().Get("rev"); rev != "oink" {
				t.Errorf("Unexpected rev: %s", rev)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{"--" + kouch.FlagRev, "oink", s.URL + "/foo/bar/baz.txt"},
			Stdout: "attachment content",
		}
	})
	tests.Add("head", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Date": []string{"Mon, 20 Aug 2018 08:55:52 GMT"}},
			Body:       ioutil.NopCloser(strings.NewReader("attachment content")),
		}, func(t *testing.T, req *http.Request) {
			if req.Method != http.MethodHead {
				t.Errorf("Unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/foo/bar/baz.txt" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{"--" + kouch.FlagHead, s.URL + "/foo/bar/baz.txt"},
			Stdout: "Content-Length: 18\r\n" +
				"Content-Type: text/plain; charset=utf-8\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n",
		}
	})
	tests.Add("dump header", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Date": []string{"Mon, 20 Aug 2018 08:55:52 GMT"}},
			Body:       ioutil.NopCloser(strings.NewReader("attachment content")),
		}, func(t *testing.T, req *http.Request) {
			if req.Method != http.MethodGet {
				t.Errorf("Unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/foo/bar/baz.txt" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{"--" + kouch.FlagDumpHeader, "%", s.URL + "/foo/bar/baz.txt"},
			Stdout: "attachment content",
			Stderr: "Content-Length: 18\r\n" +
				"Content-Type: text/plain; charset=utf-8\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n",
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"get", "att"}))
}
