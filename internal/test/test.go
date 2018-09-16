// Package test provides test utilities
package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

// CmdTest represents a single test for a complete kouch command.
type CmdTest struct {
	Args   []string
	Stdout string
	Stderr string
	Err    string
	Status int
}

// ValidateCmdTest validates a test. It is meant to be passed to
// testy's tests.Run() method.
func ValidateCmdTest(args []string) func(*testing.T, CmdTest) {
	return func(t *testing.T, test CmdTest) {
		defer testy.RestoreEnv()()
		if e := os.Setenv("HOME", "/dev/null"); e != nil {
			t.Fatal(e)
		}
		var err error
		stdout, stderr := testy.RedirIO(nil, func() {
			root := registry.Root()
			root.SetArgs(append(args, test.Args...))
			err = root.Execute()
		})
		if d := diff.Text(test.Stdout, stdout); d != nil {
			t.Errorf("STDOUT:\n%s", d)
		}
		if d := diff.Text(test.Stderr, stderr); d != nil {
			t.Errorf("STDERR:\n%s", d)
		}
		testy.ExitStatusError(t, test.Err, test.Status, err)
	}
}

// NewRequest wraps httptest.NewRequest, and add reasonable default headers.
func NewRequest(t *testing.T, method, fullPath string, body io.Reader) *http.Request {
	url, err := url.Parse(fullPath)
	if err != nil {
		t.Fatal(err)
	}
	path := url.EscapedPath()
	if q := url.Query(); len(q) > 0 {
		path = path + "?" + q.Encode()
	}
	r := httptest.NewRequest(method, path, body)
	r.Host = url.Host
	r.Header.Set("Content-Type", "application/json")
	if method != http.MethodHead {
		r.Header.Set("Accept-Encoding", "gzip")
	}
	r.Header.Set("Accept", "application/json")
	return r
}

// CheckRequest compares two request, ignoring the User-Agent header.
func CheckRequest(t *testing.T, expected, actual *http.Request) {
	delete(expected.Header, "User-Agent")
	delete(actual.Header, "User-Agent")
	if d := diff.HTTPRequest(expected, actual); d != nil {
		t.Error(d)
	}
}
