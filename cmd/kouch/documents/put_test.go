package documents

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/pkg/errors"
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
			name: "full commit",
			args: []string{"http://foo.com/foo/123", "--" + kouch.FlagFullCommit},
			expected: &opts{
				Target: &kouch.Target{
					Root:     "http://foo.com",
					Database: "foo",
					Document: "123",
				},
				fullCommit: true,
				Values:     &url.Values{},
			},
		},
		{
			name: "batch",
			args: []string{"--" + flagBatch, "docid"},
			expected: &opts{
				Target: &kouch.Target{Document: "docid"},
				Values: &url.Values{param(flagBatch): []string{"ok"}},
			},
		},
		{
			name: "new edits",
			args: []string{"--" + flagNewEdits + "=false", "docid"},
			expected: &opts{
				Target: &kouch.Target{Document: "docid"},
				Values: &url.Values{param(flagNewEdits): []string{"false"}},
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

func TestPutDocument(t *testing.T) {
	type pdTest struct {
		name     string
		content  io.ReadCloser
		opts     *opts
		resp     *http.Response
		val      testy.RequestValidator
		expected string
		err      string
		status   int
	}
	tests := []pdTest{
		{
			name:   "validation fails",
			opts:   &opts{Target: &kouch.Target{}, Values: &url.Values{}},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:    "success",
			content: ioutil.NopCloser(strings.NewReader(`{"_id":"oink"}`)),
			opts:    &opts{Target: &kouch.Target{Database: "foo", Document: "123"}, Values: &url.Values{}},
			val: func(r *http.Request) error {
				if r.Method != "PUT" {
					return errors.Errorf("Unexpected method: %s", r.Method)
				}
				if r.URL.Path != "/foo/123" {
					return errors.Errorf("Unexpected path: %s", r.URL.Path)
				}
				if ct := r.Header.Get("Content-Type"); ct != "application/json" {
					return errors.Errorf("Unexpected Content-Type: %s", ct)
				}
				var doc interface{}
				if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
					fmt.Printf("json err: %s\n", err)
					return err
				}
				return nil
			},
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(`{"ok":true,"id":"bar","rev":"1-967a00dff5e02add41819138abb3284d"}`)),
			},
			expected: `{"ok":true,"id":"bar","rev":"1-967a00dff5e02add41819138abb3284d"}`,
		},
	}
	for _, test := range tests {
		func(test pdTest) {
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
				ctx := context.Background()
				if test.content != nil {
					ctx = kouch.SetInput(ctx, test.content)
				}
				result, err := putDocument(ctx, test.opts)
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
