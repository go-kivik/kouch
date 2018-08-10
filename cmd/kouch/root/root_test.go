package root

import (
	"testing"

	"github.com/flimzy/diff"
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
