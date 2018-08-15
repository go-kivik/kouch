package root

import (
	"context"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected *cobra.Command
	}{
	// {
	// 	name:    "normal",
	// 	log:     discardLogger,
	// 	version: "1.2.3",
	// 	expected: &cobra.Command{
	// 		Version: "1.2.3",
	// 		Use:     "kouch",
	// 		Short:   "kouch is a command-line tool for interacting with CouchDB",
	// 	},
	// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := rootCmd(test.version)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestSetTarget(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		err      string
		status   int
	}{
		{
			name:     "no arguments",
			args:     nil,
			expected: "",
		},
		{
			name:   "too many arguments",
			args:   []string{"foo", "bar"},
			err:    "Too many targets provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:     "success",
			args:     []string{"foo"},
			expected: "foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var err error
			ctx, err = setTarget(ctx, test.args)
			testy.ExitStatusError(t, test.err, test.status, err)
			val := kouch.GetTarget(ctx)
			if val != test.expected {
				t.Errorf("Unexpected value: %s", val)
			}
		})
	}
}
