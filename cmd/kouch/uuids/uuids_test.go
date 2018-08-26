package uuids

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/get"
	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestGetUUIDsOpts(t *testing.T) {
	tests := []struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected *kouch.Options
		err      string
		status   int
	}{
		{
			name: "count specified",
			args: []string{"--count", "123"},
			expected: &kouch.Options{
				Target: &kouch.Target{},
				Options: &chttp.Options{
					Query: url.Values{"count": []string{"123"}},
				},
			},
		},
		{
			name: "root from context",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			expected: &kouch.Options{
				Target:  &kouch.Target{Root: "foo.com"},
				Options: &chttp.Options{},
			},
		},
		{
			name: "root from command line",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"--count", "4", "example.com:555"},
			expected: &kouch.Options{
				Target:  &kouch.Target{Root: "example.com:555"},
				Options: &chttp.Options{Query: url.Values{"count": []string{"4"}}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := uuidsCmd()
			cmd.ParseFlags(test.args)
			ctx := kouch.GetContext(cmd)
			if flags := cmd.Flags().Args(); len(flags) > 0 {
				ctx = kouch.SetTarget(ctx, flags[0])
			}
			kouch.SetContext(kouch.SetConf(ctx, test.conf), cmd)
			opts, err := getUUIDsOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestGetUUIDsCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("invalid url", test.CmdTest{
		Args:   []string{"http://%xxfoo.com"},
		Err:    `parse http://%xxfoo.com: invalid URL escape "%xx"`,
		Status: chttp.ExitStatusURLMalformed,
	})
	tests.Add("defaults", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"uuids":["3cd2f787fc320c6654befd3a4a004df6"]}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/_uuids" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
			if len(r.URL.Query()) != 0 {
				t.Errorf("Expected no query params, got %s", r.URL.RawQuery)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args:   []string{"--root", s.URL},
			Stdout: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6"]}`,
		}
	})
	tests.Add("3 uuids", func(t *testing.T) interface{} {
		s := testy.ServeResponseValidator(t, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"uuids":["3cd2f787fc320c6654befd3a4a004df6","3cd2f787fc320c6654befd3a4a005c10","3cd2f787fc320c6654befd3a4a00624e"]}`)),
		}, func(t *testing.T, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/_uuids" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
			if q := r.URL.RawQuery; q != "count=3" {
				t.Errorf("Unexpected query: %s", r.URL.RawQuery)
			}
		})
		tests.Cleanup(s.Close)
		return test.CmdTest{
			Args: []string{"--root", s.URL, "--count", "3", "--output-format", "yaml"},
			Stdout: "uuids:\n" +
				"- 3cd2f787fc320c6654befd3a4a004df6\n" +
				"- 3cd2f787fc320c6654befd3a4a005c10\n" +
				"- 3cd2f787fc320c6654befd3a4a00624e\n",
		}
	})

	tests.Run(t, test.ValidateCmdTest([]string{"get", "uuids"}))
}
