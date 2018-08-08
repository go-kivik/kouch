package uuids

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

var uuids = []interface{}{
	"3cd2f787fc320c6654befd3a4a004df6", "3cd2f787fc320c6654befd3a4a005c10",
	"3cd2f787fc320c6654befd3a4a00624e", "3cd2f787fc320c6654befd3a4a007099",
	"3cd2f787fc320c6654befd3a4a007898", "3cd2f787fc320c6654befd3a4a007c60",
	"3cd2f787fc320c6654befd3a4a008b53", "3cd2f787fc320c6654befd3a4a009675",
	"3cd2f787fc320c6654befd3a4a009ad0", "3cd2f787fc320c6654befd3a4a00a9fb",
}

type uuidResponse struct {
	count int
}

var _ http.Handler = &uuidResponse{}

func (ur *uuidResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count, err := strconv.Atoi("0" + r.URL.Query().Get("count"))
	if err != nil {
		w.WriteHeader(kivik.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"reason":"Unparseable count param: %s"}`, err)))
		return
	}
	if count != ur.count {
		w.WriteHeader(kivik.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"reason":"Unexpected count param: Got %d, expected %d"}`, count, ur.count)))
		return
	}
	result := map[string]interface{}{
		"uuids": uuids[0:count],
	}
	body, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	w.Write(body)
}

func uuidServer(r *uuidResponse) (url string, close func()) {
	s := httptest.NewServer(r)
	return s.URL, func() { s.Close() }
}

func TestGetUUIDs(t *testing.T) {
	type guTest struct {
		name     string
		opts     *getUUIDsOpts
		expected string
		err      string
		cleanup  func()
	}
	tests := []guTest{
		func() guTest {
			url, close := uuidServer(&uuidResponse{count: 1})
			return guTest{
				name:     "defaults",
				opts:     &getUUIDsOpts{Count: 1, Root: url},
				expected: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6"]}`,
				cleanup:  close,
			}
		}(),
		func() guTest {
			url, close := uuidServer(&uuidResponse{count: 3})
			return guTest{
				name:     "3 uuids",
				opts:     &getUUIDsOpts{Count: 3, Root: url},
				expected: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6","3cd2f787fc320c6654befd3a4a005c10","3cd2f787fc320c6654befd3a4a00624e"]}`,
				cleanup:  close,
			}
		}(),
		{
			name: "invalid url",
			opts: &getUUIDsOpts{Count: 1, Root: "http://%xxfoo.com/"},
			err:  `parse http://%xxfoo.com/: invalid URL escape "%xx"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanup != nil {
				defer test.cleanup()
			}
			result, err := getUUIDs(test.opts)
			testy.Error(t, test.err, err)
			defer result.Close()
			resultJSON, err := ioutil.ReadAll(result)
			if err != nil {
				t.Fatal(err)
			}
			if d := diff.JSON([]byte(test.expected), resultJSON); d != nil {
				t.Error(d)
			}
		})
	}
}

var fooConf = &kouch.Config{
	DefaultContext: "foo",
	Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
}

func TestGetUUIDsOpts(t *testing.T) {
	tests := []struct {
		name     string
		conf     *kouch.Config
		cmd      *cobra.Command
		args     []string
		expected interface{}
		err      string
		status   int
	}{
		{
			name:   "no context",
			conf:   &kouch.Config{},
			err:    "No default context",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name: "default context",
			conf: fooConf,
			expected: &getUUIDsOpts{
				Count: 1,
				Root:  "foo.com",
			},
		},
		{
			name: "count from args",
			conf: fooConf,
			args: []string{"--count", "3"},
			expected: &getUUIDsOpts{
				Count: 3,
				Root:  "foo.com",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cx := &getUUIDsCtx{&kouch.CmdContext{
				Conf: test.conf,
			}}
			cmd := getUUIDsCmd(cx.CmdContext)
			cmd.ParseFlags(test.args)
			opts, err := cx.getUUIDsOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}
