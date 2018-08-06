package kouch

import (
	"errors"
	"testing"

	"github.com/go-kivik/couchdb/chttp"
)

func TestExit(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
		status   int
	}{
		{
			name:     "standard error",
			err:      errors.New("foo"),
			expected: "foo",
			status:   chttp.ExitUnknownFailure,
		},
		{
			name:     "HTTP status error",
			err:      &httpErr{error: errors.New("foo"), httpStatus: 404, exitStatus: chttp.ExitNotRetrieved},
			expected: "The requested URL returned error: 404 foo",
			status:   chttp.ExitNotRetrieved,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, status := exit(test.err)
			if msg != test.expected {
				t.Errorf("Unexpected result:\nExpected: %s\n  Actual: %s\n", test.expected, msg)
			}
			if status != test.status {
				t.Errorf("Unexpected exit status:\nExpected: %d\n  Actual: %d\n", test.status, status)
			}
		})
	}
}

type httpErr struct {
	httpStatus, exitStatus int
	error
}

func (e *httpErr) StatusCode() int {
	return e.httpStatus
}

func (e *httpErr) ExitStatus() int {
	return e.exitStatus
}
