package database

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/create"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestCreateDatabaseCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("validation fails", test.CmdTest{
		Args:   []string{},
		Err:    "no URL specified",
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
		return test.CmdTest{
			Args:   []string{s.URL + "/oink", "--" + kouch.FlagShards, "5"},
			Stdout: `{"ok":true}`,
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"create", "database"}))
}
