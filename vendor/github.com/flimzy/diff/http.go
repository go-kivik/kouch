package diff

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// HTTPRequest compares the metadata and bodies of the two HTTP requests, and
// returns the difference.
func HTTPRequest(expected, actual *http.Request) *Result {
	expDump, err := httputil.DumpRequest(expected, true)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump expected request: %s", err)}
	}
	actDump, err := httputil.DumpRequest(actual, true)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump actual request: %s", err)}
	}
	return Text(string(expDump), string(actDump))
}

// HTTPResponse compares the metadata and bodies of the two HTTP responses, and
// returns the difference.
func HTTPResponse(expected, actual *http.Response) *Result {
	expDump, err := httputil.DumpResponse(expected, true)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump expected response: %s", err)}
	}
	actDump, err := httputil.DumpResponse(actual, true)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump actual response: %s", err)}
	}
	return Text(string(expDump), string(actDump))
}
