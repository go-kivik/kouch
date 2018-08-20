package util

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
)

func TestChttpDo(t *testing.T) {
	type cdTest struct {
		path         string
		options      *kouch.Options
		head, body   bool // Indicate whether to pass an io.Writer for the respective part
		expectedHead string
		expectedBody string
		err          string
		status       int
		cleanup      func()
	}
	tests := map[string]func(*testing.T) cdTest{
		"body": func(t *testing.T) cdTest {
			s := testy.ServeResponseValidator(t,
				&http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("body content\n")),
				},
				func(t *testing.T, r *http.Request) {
					if r.Method != http.MethodGet {
						t.Errorf("Unexpected method: %s", r.Method)
					}
					if r.URL.Path != "/foo/bar" {
						t.Errorf("Unexpected path: %s", r.URL.Path)
					}
				})
			return cdTest{
				options:      &kouch.Options{Target: &kouch.Target{Root: s.URL + "/foo"}},
				path:         "bar",
				body:         true,
				expectedBody: "body content\n",
				cleanup:      s.Close,
			}
		},
		"head": func(t *testing.T) cdTest {
			s := testy.ServeResponseValidator(t,
				&http.Response{
					StatusCode:    200,
					ContentLength: 13,
					Header: http.Header{
						"Content-Length": []string{"13"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
						"Date":           []string{"Mon, 20 Aug 2018 10:23:52 GMT"},
					},
				},
				func(t *testing.T, r *http.Request) {
					if r.Method != http.MethodHead {
						t.Errorf("Unexpected method: %s", r.Method)
					}
					if r.URL.Path != "/foo/bar" {
						t.Errorf("Unexpected path: %s", r.URL.Path)
					}
				})
			return cdTest{
				options: &kouch.Options{Target: &kouch.Target{Root: s.URL + "/foo"}},
				path:    "bar",
				head:    true,
				expectedHead: "Content-Length: 13\r\n" +
					"Content-Type: text/plain; charset=utf-8\r\n" +
					"Date: Mon, 20 Aug 2018 10:23:52 GMT\r\n",
				cleanup: s.Close,
			}
		},
		"both": func(t *testing.T) cdTest {
			s := testy.ServeResponseValidator(t,
				&http.Response{
					StatusCode: 200,
					Header: http.Header{
						"Date": []string{"Mon, 20 Aug 2018 10:23:52 GMT"},
					},
					Body: ioutil.NopCloser(strings.NewReader("body content\n")),
				},
				func(t *testing.T, r *http.Request) {
					if r.Method != http.MethodGet {
						t.Errorf("Unexpected method: %s", r.Method)
					}
					if r.URL.Path != "/foo/bar" {
						t.Errorf("Unexpected path: %s", r.URL.Path)
					}
				})
			return cdTest{
				options: &kouch.Options{Target: &kouch.Target{Root: s.URL + "/foo"}},
				path:    "bar",
				head:    true,
				body:    true,
				expectedHead: "Content-Length: 13\r\n" +
					"Content-Type: text/plain; charset=utf-8\r\n" +
					"Date: Mon, 20 Aug 2018 10:23:52 GMT\r\n",
				expectedBody: "body content\n",
				cleanup:      s.Close,
			}
		},
	}
	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			test := fn(t)
			var head, body *bytes.Buffer
			if test.head {
				head = &bytes.Buffer{}
			}
			if test.body {
				fmt.Printf("body\n")
				body = &bytes.Buffer{}
			}
			err := ChttpGet(context.Background(), test.path, test.options, head, body)
			testy.ExitStatusError(t, test.err, test.status, err)
			var resultHead, resultBody string
			if head != nil {
				resultHead = head.String()
			}
			if body != nil {
				resultBody = body.String()
			}
			if d := diff.Text(test.expectedHead, resultHead); d != nil {
				t.Errorf("Head differs:\n%s", d)
			}
			if d := diff.Text(test.expectedBody, resultBody); d != nil {
				t.Errorf("Body differs:\n%s", d)
			}
		})
	}
}
