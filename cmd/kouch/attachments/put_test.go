package attachments

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

func TestPutAttachmentCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "No filename provided",
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
			if r.URL.Path != "/foo/bar/baz.txt" {
				t.Errorf("Unexpected req path: %s", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "text/plain" {
				t.Errorf("Unexpected Content-Type: %s", ct)
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != `{"oink":"foo"}` {
				t.Errorf("Unexpected body: %s", string(body))
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar/baz.txt", "-d", `{"oink":"foo"}`, "-F", "yaml", "--content-type", "text/plain"},
			Stdout: "id: bar\nok: true\nrev: 1-967a00dff5e02add41819138abb3284d",
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
			Args:   []string{s.URL + "/foo/bar/baz.txt", "-d", `{"oink":"foo"}`, "-F", "yaml", "--content-type", "text/plain", "--auto-rev"},
			Stdout: "id: bar\nok: true\nrev: 2-967a00dff5e02add41819138abb3284d",
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"put", "att"}))
}
