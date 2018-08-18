package attachments

import (
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name   string
		target *kouch.Target
		err    string
		status int
	}{
		{
			name:   "no filename",
			target: &kouch.Target{},
			err:    "No filename provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no doc id",
			target: &kouch.Target{Filename: "foo.txt"},
			err:    "No document ID provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no database provided",
			target: &kouch.Target{Document: "123", Filename: "foo.txt"},
			err:    "No database name provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "no root url",
			target: &kouch.Target{Database: "foo", Document: "123", Filename: "foo.txt"},
			err:    "No root URL provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "valid",
			target: &kouch.Target{Root: "xxx", Database: "foo", Document: "123", Filename: "foo.txt"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateTarget(test.target)
			testy.ExitStatusError(t, test.err, test.status, err)
		})
	}
}
