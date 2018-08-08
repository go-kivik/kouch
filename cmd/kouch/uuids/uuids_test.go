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
	"github.com/go-kivik/kivik"
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
		count    int
		url      string
		expected string
		err      string
		cleanup  func()
	}
	tests := []guTest{
		func() guTest {
			url, close := uuidServer(&uuidResponse{count: 1})
			return guTest{
				name:     "defaults",
				count:    1,
				url:      url + "/_uuids",
				expected: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6"]}`,
				cleanup:  close,
			}
		}(),
		func() guTest {
			url, close := uuidServer(&uuidResponse{count: 3})
			return guTest{
				name:     "3 uuids",
				count:    3,
				url:      url + "/_uuids",
				expected: `{"uuids":["3cd2f787fc320c6654befd3a4a004df6","3cd2f787fc320c6654befd3a4a005c10","3cd2f787fc320c6654befd3a4a00624e"]}`,
				cleanup:  close,
			}
		}(),
		{
			name: "invalid url",
			url:  "http://%xxfoo.com/",
			err:  `parse http://%xxfoo.com/: invalid URL escape "%xx"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if test.cleanup != nil {
				defer test.cleanup()
			}
			result, err := getUUIDs(test.url, test.count)
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
