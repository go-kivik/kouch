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
		{
			name: "head",
			args: []string{"--" + kouch.FlagHead, "foo.txt"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Filename: "foo.txt"},
				Options: &chttp.Options{},
				Head:    true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := getAttCmd()
			cmd.ParseFlags(test.args)
			ctx := kouch.GetContext(cmd)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(kouch.SetConf(ctx, test.conf), cmd)
			opts, err := getAttachmentOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestGetAttachment(t *testing.T) {
	type gaTest struct {
		name     string
		opts     *kouch.Options
		resp     *http.Response
		val      testy.RequestValidator
		expected string
		err      string
		status   int
	}
	tests := []gaTest{
		{
			name:   "validation fails",
			opts:   &kouch.Options{Target: &kouch.Target{}},
			err:    "No filename provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "success",
			opts: &kouch.Options{Target: &kouch.Target{Database: "foo", Document: "123", Filename: "foo.txt"}},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.Path != "/foo/123/foo.txt" {
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
			opts: &kouch.Options{Target: &kouch.Target{Database: "foo/ba r", Document: "123/b", Filename: "foo/bar.txt"}},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb/foo%2Fbar.txt" {
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
				Target:  &kouch.Target{Database: "foo/ba r", Document: "123/b", Filename: "foo/bar.txt"},
				Options: &chttp.Options{IfNoneMatch: "xyz"},
			},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb/foo%2Fbar.txt" {
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
			name: "rev",
			opts: &kouch.Options{
				Target: &kouch.Target{Database: "foo/ba r", Document: "123/b", Filename: "foo/bar.txt"},
				Options: &chttp.Options{
					Query: url.Values{"rev": []string{"xyz"}},
				},
			},
			val: func(t *testing.T, r *http.Request) {
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb/foo%2Fbar.txt" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
				if rev := r.URL.Query().Get("rev"); rev != "xyz" {
					t.Errorf("Unexpected revision: %s", rev)
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
				Target:  &kouch.Target{Database: "foo/ba r", Document: "123/b", Filename: "foo/bar.txt"},
				Options: &chttp.Options{},
				Head:    true,
			},
			val: func(t *testing.T, r *http.Request) {
				if r.Method != http.MethodHead {
					t.Errorf("Unepxected method: %s", r.Method)
				}
				if r.URL.RawPath != "/foo%2Fba+r/123%2Fb/foo%2Fbar.txt" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"Date": []string{"Mon, 20 Aug 2018 08:55:52 GMT"},
				},
				Body: ioutil.NopCloser(strings.NewReader("Test\ncontent\n")),
			},
			expected: "Content-Length: 13\r\n" +
				"Content-Type: text/plain; charset=utf-8\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n",
		},
	}
	for _, test := range tests {
		func(test gaTest) {
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
