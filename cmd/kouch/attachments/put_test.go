package attachments

import (
	"context"
	"io"
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

func TestPutAttachmentOpts(t *testing.T) {
	input := ioutil.NopCloser(strings.NewReader(""))
	tests := []struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected interface{}
		err      string
		status   int
	}{
		{
			name: "rev",
			args: []string{"--" + kouch.FlagRev, "xyz", "foo.txt"},
			expected: &kouch.Options{
				Target: &kouch.Target{Filename: "foo.txt"},
				Options: &chttp.Options{
					Query: url.Values{"rev": []string{"xyz"}},
					Body:  input,
				},
			},
		},
		{
			name: "content type",
			args: []string{"--" + flagContentType, "image/oink", "foo.txt"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Filename: "foo.txt"},
				Options: &chttp.Options{ContentType: "image/oink", Body: input},
			},
		},
		{
			name: "guess content type",
			args: []string{"--" + flagGuessContentType, "foo.txt"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Filename: "foo.txt"},
				Options: &chttp.Options{ContentType: "text/plain; charset=utf-8", Body: input},
			},
		},
		{
			name: "guess content type failure",
			args: []string{"--" + flagGuessContentType, "foo.xxxxxxx"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Filename: "foo.xxxxxxx"},
				Options: &chttp.Options{ContentType: "application/octet-stream", Body: input},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := putAttCmd()
			cmd.ParseFlags(test.args)
			ctx := kouch.GetContext(cmd)
			ctx = kouch.SetInput(ctx, input)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(kouch.SetConf(ctx, test.conf), cmd)
			opts, err := putAttachmentOpts(cmd)
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestPutAttachment(t *testing.T) {
	type gaTest struct {
		name     string
		content  io.ReadCloser
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
			opts: &kouch.Options{
				Target: &kouch.Target{Database: "foo", Document: "oink", Filename: "foo.txt"},
				Options: &chttp.Options{
					ContentType: "text/plain",
					Body:        ioutil.NopCloser(strings.NewReader("test data")),
				},
			},
			val: func(t *testing.T, r *http.Request) {
				defer r.Body.Close()
				if r.URL.Path != "/foo/oink/foo.txt" {
					t.Errorf("Unexpected path: %s", r.URL.Path)
				}
				if ct := r.Header.Get("Content-Type"); ct != "text/plain" {
					t.Errorf("Unexpected Content-Type: %s", ct)
				}
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if d := diff.Text("test data", body); d != nil {
					t.Errorf("Unexpected body: %s", d)
				}
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(`{"ok":true,"id":"oink","rev":"3-13438fbeeac7271383a42b57511f03ea"}`)),
			},
			expected: `{"ok":true,"id":"oink","rev":"3-13438fbeeac7271383a42b57511f03ea"}`,
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
				ctx := context.Background()
				if test.content != nil {
					ctx = kouch.SetInput(ctx, test.content)
				}
				result, err := putAttachment(ctx, test.opts)
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
