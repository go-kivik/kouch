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
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/delete"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestDeleteDatabaseCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "no server root specified",
		Status: chttp.ExitFailedToInitialize,
	})
	tests.Add("delete success", func(t *testing.T) interface{} {
		var s *httptest.Server
		s = testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, func(t *testing.T, r *http.Request) {
			expected := test.NewRequest(t, "DELETE", s.URL+"/oink", nil)
			test.CheckRequest(t, expected, r)
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{s.URL + "/oink"},
			Stdout: `{"ok":true}`,
		}
	})
	tests.Add("auth in target", func(t *testing.T) interface{} {
		var s *httptest.Server
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expected := test.NewRequest(t, "DELETE", s.URL+"/oink", nil)
			expected.Header.Set("Authorization", "Basic YWRtaW46YWJjMTIz")
			test.CheckRequest(t, expected, r)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(kivik.StatusCreated)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		tests.Cleanup(s.Close)
		addr, _ := url.Parse(s.URL)
		addr.User = url.UserPassword("admin", "abc123")
		addr.Path = "/oink"
		return test.CmdTest{
			Args:   []string{addr.String()},
			Stdout: `{"ok":true}`,
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"delete", "database"}))
}
