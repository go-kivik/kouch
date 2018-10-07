package outputraw

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"
	"github.com/spf13/cobra"
)

func TestRawModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &RawMode{}
	mode.AddFlags(cmd.PersistentFlags())

	test.Flags(t, []string{}, cmd)
}

func TestRawNew(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		parseErr string
		expected io.Writer
		err      string
	}{
		{
			name:     "happy path",
			args:     nil,
			expected: &bytes.Buffer{},
		},
		{
			name:     "invalid args",
			args:     []string{"--foo"},
			parseErr: "unknown flag: --foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			cmd := &cobra.Command{}
			mode := &RawMode{}
			mode.AddFlags(cmd.PersistentFlags())

			err := cmd.ParseFlags(test.args)
			testy.Error(t, test.parseErr, err)

			ctx = kouch.SetFlags(ctx, cmd.Flags())
			result, err := mode.New(ctx, &bytes.Buffer{})
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}
