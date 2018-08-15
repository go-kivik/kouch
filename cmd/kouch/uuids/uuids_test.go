package uuids

import (
	"context"
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
		ctx      context.Context
		opts     *opts
		expected string
		err      string
		cleanup  func()
	}
	tests := []guTest{
		func() guTest {
			url, close := uuidServer(&uuidResponse{count: 1})
			return guTest{
				name:     "defaults",
				opts:     &opts{root: url, count: 1},
				expected: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6"]}`,
				cleanup:  close,
			}
		}(),
		func() guTest {
			url, close := uuidServer(&uuidResponse{count: 3})
			return guTest{
				name:     "3 uuids",
				opts:     &opts{root: url, count: 3},
				expected: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6","3cd2f787fc320c6654befd3a4a005c10","3cd2f787fc320c6654befd3a4a00624e"]}`,
				cleanup:  close,
			}
		}(),
		{
			name: "invalid url",
			opts: &opts{root: "http://%xxfoo.com/"},
			err:  `parse http://%xxfoo.com/: invalid URL escape "%xx"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanup != nil {
				defer test.cleanup()
			}
			ctx := test.ctx
			if ctx == nil {
				ctx = context.Background()
			}
			result, err := getUUIDs(ctx, test.opts)
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

func TestGetUUIDsOpts(t *testing.T) {
	tests := []struct {
		name     string
		conf     *kouch.Config
		args     []string
		expected *opts
		err      string
		status   int
	}{
		{
			name:     "count specified",
			args:     []string{"--count", "123"},
			expected: &opts{count: 123},
		},
		{
			name: "root from context",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			expected: &opts{
				root:  "foo.com",
				count: 1,
			},
		},
		{
			name: "root from command line",
			conf: &kouch.Config{
				DefaultContext: "foo",
				Contexts:       []kouch.NamedContext{{Name: "foo", Context: &kouch.Context{Root: "foo.com"}}},
			},
			args: []string{"--count", "4", "example.com:555"},
			expected: &opts{
				root:  "example.com:555",
				count: 4,
			},
		},
		{
			name:   "too many arguments",
			args:   []string{"foo", "bar"},
			err:    "Too many targets provided",
			status: chttp.ExitFailedToInitialize,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.conf == nil {
				test.conf = &kouch.Config{}
			}
			cmd := uuidsCmd()
			kouch.SetContext(kouch.SetConf(kouch.GetContext(cmd), test.conf), cmd)
			cmd.ParseFlags(test.args)
			opts, err := getUUIDsOpts(cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}
