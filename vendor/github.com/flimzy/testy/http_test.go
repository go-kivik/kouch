package testy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flimzy/diff"
)

func TestServeResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		expected *http.Response
	}{
		{
			name: "Simple GET response",
			response: &http.Response{
				Header: http.Header{
					"X-Foo": []string{"foo"},
					"Date":  []string{"Tue, 07 Aug 2018 20:18:51 GMT"},
				},
			},
			expected: &http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"X-Foo": []string{"foo"},
					"Date":  []string{"Tue, 07 Aug 2018 20:18:51 GMT"},
				},
			},
		},
		{
			name: "Simple response with body",
			response: &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"Date": []string{"Tue, 07 Aug 2018 20:18:51 GMT"},
				},
				Body: ioutil.NopCloser(strings.NewReader("the body\nof the response\n")),
			},
			expected: &http.Response{
				StatusCode:    200,
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: 25,
				Header: http.Header{
					"Date":         []string{"Tue, 07 Aug 2018 20:18:51 GMT"},
					"Content-Type": []string{"text/plain; charset=utf-8"},
				},
				Body: ioutil.NopCloser(strings.NewReader("the body\nof the response\n")),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := httptest.NewServer(&ResponseHandler{test.response})
			defer s.Close()
			res, err := http.Get(s.URL)
			if err != nil {
				t.Fatal(err)
			}
			if d := diff.HTTPResponse(test.expected, res); d != nil {
				t.Error(d)
			}
		})
	}
}
