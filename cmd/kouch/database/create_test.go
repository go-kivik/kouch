package database

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/create"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestCreateDatabaseCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "no server root specified",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("create success", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 201,
			Body:       ioutil.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/oink" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/oink"},
			Stdout: `{"ok":true}`,
		}
	})
	tests.Add("shards", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 201,
			Body:       ioutil.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/oink" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
			if q := r.URL.Query().Get("q"); q != "5" {
				t.Errorf("Unexpected q value: %v", q)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/oink", "--" + kouch.FlagShards, "5"},
			Stdout: `{"ok":true}`,
		}
	})
	tests.Add("auth in target", func(t *testing.T) interface{} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method != kivik.MethodPut {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/oink" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
			if q := r.URL.Query().Get("q"); q != "5" {
				t.Errorf("Unexpected q value: %v", q)
			}
			if auth := r.Header.Get("Authorization"); auth != "Basic YWRtaW46YWJjMTIz" {
				t.Errorf("Unexpected Authorization header: %s", auth)
			}
			w.WriteHeader(kivik.StatusCreated)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		tests.Cleanup(s.Close)
		addr, _ := url.Parse(s.URL)
		addr.User = url.UserPassword("admin", "abc123")
		addr.Path = "/oink"
		return test.CmdTest{
			Args:   []string{addr.String(), "--" + kouch.FlagShards, "5"},
			Stdout: `{"ok":true}`,
		}
	})
	tests.Add("auth in cli opts", func(t *testing.T) interface{} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method != kivik.MethodPut {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/oink" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
			if q := r.URL.Query().Get("q"); q != "5" {
				t.Errorf("Unexpected q value: %v", q)
			}
			if auth := r.Header.Get("Authorization"); auth != "Basic YWRtaW46YWJjMTIz" {
				t.Errorf("Unexpected Authorization header: %s", auth)
			}
			w.WriteHeader(kivik.StatusCreated)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/oink", "--" + kouch.FlagShards, "5", "--" + kouch.FlagUser, "admin", "--" + kouch.FlagPassword, "abc123"},
			Stdout: `{"ok":true}`,
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"create", "database"}))
}
