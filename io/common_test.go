package io

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func testOptions(t *testing.T, expected []string, cmd *cobra.Command) {
	found := make([]string, 0)
	cmd.ParseFlags(nil)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		found = append(found, f.Name)
	})
	if d := diff.Interface(expected, found); d != nil {
		t.Error(d)
	}
}
