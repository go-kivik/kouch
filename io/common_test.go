package io

import (
	"errors"
	"io"
	"testing"

	"github.com/flimzy/diff"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type errWriter struct{}

var _ io.Writer = &errWriter{}

func (w *errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("errWriter: write error")
}

func testOptions(t *testing.T, expected []string, cmd *cobra.Command) {
	found := make([]string, 0)
	if e := cmd.ParseFlags(nil); e != nil {
		t.Fatal(e)
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		found = append(found, f.Name)
	})
	if d := diff.Interface(expected, found); d != nil {
		t.Error(d)
	}
}

type errReader struct{}

var _ io.Reader = &errReader{}

func (r *errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("errReader: read error")
}
