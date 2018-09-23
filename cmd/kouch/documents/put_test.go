package documents

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/put"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestPutDocumentOpts(t *testing.T) {
	tests := testy.NewTable()

	tests.Add("duplicate id", test.OptionsTest{
		Args:   []string{"--" + kouch.FlagDocument, "foo", "bar"},
		Err:    "Must not use --" + kouch.FlagDocument + " and pass document ID as part of the target",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("id from target", test.OptionsTest{
		Conf: &kouch.Config{
			DefaultContext: "foo",
			Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
		},
		Args: []string{"--database", "bar", "123"},
		Expected: &kouch.Options{
			Target: &kouch.Target{
				Root:     "foo.com",
				Database: "bar",
				Document: "123",
			},
			Options: &chttp.Options{},
		},
	})
	tests.Add("db included in target", test.OptionsTest{
		Conf: &kouch.Config{
			DefaultContext: "foo",
			Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
		},
		Args: []string{"/foo/123"},
		Expected: &kouch.Options{
			Target: &kouch.Target{
				Root:     "foo.com",
				Database: "foo",
				Document: "123",
			},
			Options: &chttp.Options{},
		},
	})
	tests.Add("db provided twice", test.OptionsTest{
		Args:   []string{"/foo/123/foo.txt", "--" + kouch.FlagDatabase, "foo"},
		Err:    "Must not use --" + kouch.FlagDatabase + " and pass database as part of the target",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("full url target", test.OptionsTest{
		Args: []string{"http://foo.com/foo/123"},
		Expected: &kouch.Options{
			Target: &kouch.Target{
				Root:     "http://foo.com",
				Database: "foo",
				Document: "123",
			},
			Options: &chttp.Options{},
		},
	})
	tests.Add("full commit", test.OptionsTest{
		Args: []string{"http://foo.com/foo/123", "--" + kouch.FlagFullCommit},
		Expected: &kouch.Options{
			Target: &kouch.Target{
				Root:     "http://foo.com",
				Database: "foo",
				Document: "123",
			},
			Options: &chttp.Options{
				FullCommit: true,
			},
		},
	})
	tests.Add("batch", test.OptionsTest{
		Args: []string{"--" + kouch.FlagBatch, "docid"},
		Expected: &kouch.Options{
			Target: &kouch.Target{Document: "docid"},
			Options: &chttp.Options{
				Query: url.Values{param(kouch.FlagBatch): []string{"ok"}},
			},
		},
	})
	tests.Add("new edits", test.OptionsTest{
		Args: []string{"--" + kouch.FlagNewEdits + "=false", "docid"},
		Expected: &kouch.Options{
			Target: &kouch.Target{Document: "docid"},
			Options: &chttp.Options{
				Query: url.Values{param(kouch.FlagNewEdits): []string{"false"}},
			},
		},
	})

	tests.Run(t, test.Options(putDocCmd, putDocumentOpts))
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
			_, _ = w.Write([]byte(`{"ok":true,"id":"bar","rev":"2-967a00dff5e02add41819138abb3284d"}`))
		}))
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar", "-d", `{"oink":foo}`, "-F", "yaml", "--auto-rev"},
			Stdout: "id: bar\nok: true\nrev: 2-967a00dff5e02add41819138abb3284d",
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"put", "doc"}))
}
