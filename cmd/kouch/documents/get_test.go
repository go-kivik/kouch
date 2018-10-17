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

	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestGetDocumentOpts(t *testing.T) {
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
	tests.Add("if-none-match", test.OptionsTest{
		Args: []string{"--" + kouch.FlagIfNoneMatch, "foo", "foo.com/bar/baz"},
		Expected: &kouch.Options{
			Target: &kouch.Target{
				Root:     "foo.com",
				Database: "bar",
				Document: "baz",
			},
			Options: &chttp.Options{IfNoneMatch: "foo"},
		},
	})
	tests.Add("rev", test.OptionsTest{
		Args: []string{"--" + kouch.FlagRev, "foo", "foo.com/bar/baz"},
		Expected: &kouch.Options{
			Target: &kouch.Target{
				Root:     "foo.com",
				Database: "bar",
				Document: "baz",
			},
			Options: &chttp.Options{
				Query: url.Values{"rev": []string{"foo"}},
			},
		},
	})
	tests.Add("attachments since", test.OptionsTest{
		Args: []string{"--" + kouch.FlagAttsSince, "foo,bar,baz", "docid"},
		Expected: &kouch.Options{
			Target: &kouch.Target{Document: "docid"},
			Options: &chttp.Options{
				Query: url.Values{param(kouch.FlagAttsSince): []string{`["foo","bar","baz"]`}},
			},
		},
	})
	tests.Add("open revs", test.OptionsTest{
		Args: []string{"--" + kouch.FlagOpenRevs, "foo,bar,baz", "docid"},
		Expected: &kouch.Options{
			Target: &kouch.Target{Document: "docid"},
			Options: &chttp.Options{
				Query: url.Values{param(kouch.FlagOpenRevs): []string{`["foo","bar","baz"]`}},
			},
		},
	})
	for _, flag := range []string{
		kouch.FlagIncludeAttachments, kouch.FlagIncludeAttEncoding,
		kouch.FlagConflicts, kouch.FlagIncludeDeletedConflicts, kouch.FlagForceLatest,
		kouch.FlagIncludeLocalSeq, kouch.FlagMeta, kouch.FlagRevs, kouch.FlagRevsInfo,
	} {
		tests.Add(flag, test.OptionsTest{
			Args: []string{"--" + flag, "docid"},
			Expected: &kouch.Options{
				Target: &kouch.Target{Document: "docid"},
				Options: &chttp.Options{
					Query: url.Values{param(flag): []string{"true"}},
				},
			},
		})
	}

	tests.Run(t, test.Options(getDocCmd, getDocumentOpts))
}

func TestGetDocumentCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "No document ID provided",
		Status: chttp.ExitFailedToInitialize,
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
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar"},
			Stdout: `{"foo":123}`,
		}
	})
	tests.Add("slashes", func(t *testing.T) interface{} {
		var s *httptest.Server
		s = testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			expected := test.NewRequest(t, "GET", s.URL+"/foo%2Fba+r/123%2Fb", nil)
			test.CheckRequest(t, expected, req)
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{
				"--" + kouch.FlagServerRoot, s.URL,
				"--" + kouch.FlagDatabase, "foo/ba r",
				"--" + kouch.FlagDocument, "123/b",
			},
			Stdout: `{"foo":123}`,
		}
	})
	tests.Add("if-none-match", func(t *testing.T) interface{} {
		var s *httptest.Server
		s = testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			expected := test.NewRequest(t, "GET", s.URL+"/foo/bar", nil)
			expected.Header.Set("If-None-Match", `"oink"`)
			test.CheckRequest(t, expected, req)
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{"--" + kouch.FlagIfNoneMatch, "oink", s.URL + "/foo/bar"},
			Stdout: `{"foo":123}`,
		}
	})
	tests.Add("rev", func(t *testing.T) interface{} {
		var s *httptest.Server
		s = testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			expected := test.NewRequest(t, "GET", s.URL+"/foo/bar?rev=oink", nil)
			test.CheckRequest(t, expected, req)
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{"--" + kouch.FlagRev, "oink", s.URL + "/foo/bar"},
			Stdout: `{"foo":123}`,
		}
	})
	tests.Add("head", func(t *testing.T) interface{} {
		var s *httptest.Server
		s = testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"Date":         []string{"Mon, 20 Aug 2018 08:55:52 GMT"},
			},
			Body: ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			expected := test.NewRequest(t, "HEAD", s.URL+"/foo/bar/baz.txt", nil)
			test.CheckRequest(t, expected, req)
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{"--" + kouch.FlagHead, s.URL + "/foo/bar/baz.txt"},
			Stdout: "Content-Length: 11\r\n" +
				"Content-Type: application/json\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n",
		}
	})
	tests.Add("yaml", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/foo/bar", "--" + kouch.FlagOutputFormat, "yaml"},
			Stdout: `foo: 123`,
		}
	})
	tests.Add("dump header to stdout", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"Date":         []string{"Mon, 20 Aug 2018 08:55:52 GMT"},
			},
			Body: ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{s.URL + "/foo/bar", "--" + kouch.FlagOutputFormat, "yaml", "--dump-header", "-"},
			Stdout: "Content-Length: 11\r\n" +
				"Content-Type: application/json\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n" +
				"\r\n" +
				"foo: 123",
		}
	})
	tests.Add("dump header to stderr", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"Date":         []string{"Mon, 20 Aug 2018 08:55:52 GMT"},
			},
			Body: ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
		}, func(t *testing.T, req *http.Request) {
			if req.URL.Path != "/foo/bar" {
				t.Errorf("Unexpected req path: %s", req.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{s.URL + "/foo/bar", "--" + kouch.FlagOutputFormat, "yaml", "--dump-header", "%"},
			Stderr: "Content-Length: 11\r\n" +
				"Content-Type: application/json\r\n" +
				"Date: Mon, 20 Aug 2018 08:55:52 GMT\r\n",
			Stdout: "foo: 123",
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"get", "doc"}))
}
