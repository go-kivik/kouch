package diff

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRequest(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual *http.Request
		result           string
	}{
		{
			name:     "two GET requests",
			expected: httptest.NewRequest("GET", "/", nil),
			actual:   httptest.NewRequest("GET", "/", nil),
		},
		{
			name:     "different URLs requests",
			expected: httptest.NewRequest("GET", "/", nil),
			actual:   httptest.NewRequest("GET", "/foo.html", nil),
			result:   "--- expected\n+++ actual\n@@ -1,3 +1,3 @@\n-GET / HTTP/1.1\r\n+GET /foo.html HTTP/1.1\r\n Host: example.com\r\n \r\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HTTPRequest(test.expected, test.actual)
			var resultText string
			if result != nil {
				resultText = result.String()
			}
			if d := Text(test.result, resultText); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestHTTPResponse(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual *http.Response
		result           string
	}{
		{
			name:     "two empty responses",
			expected: &http.Response{},
			actual:   &http.Response{},
		},
		{
			name: "Different headers",
			expected: &http.Response{
				Header: http.Header{"Foo": []string{"bar"}},
			},
			actual: &http.Response{
				Header: http.Header{"Foo": []string{"qux"}},
			},
			result: "--- expected\n+++ actual\n@@ -1,4 +1,4 @@\n HTTP/0.0 000 status code 0\r\n-Foo: bar\r\n+Foo: qux\r\n Content-Length: 0\r\n \r\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HTTPResponse(test.expected, test.actual)
			var resultText string
			if result != nil {
				resultText = result.String()
			}
			if d := Text(test.result, resultText); d != nil {
				t.Error(d)
			}
		})
	}
}
