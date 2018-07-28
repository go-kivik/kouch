package io

import (
	"bytes"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/spf13/cobra"
)

func TestRawModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &rawMode{}
	mode.config(cmd)

	testOptions(t, []string{}, cmd)
}

func TestRawOutput(t *testing.T) {
	input := `{foo bar baz}`
	p := &rawProcessor{}
	t.Run("happy path", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := p.Output(buf, []byte(input))
		testy.Error(t, "", err)
		if d := diff.Text(input, buf.String()); d != nil {
			t.Error(d)
		}
	})
	t.Run("write error", func(t *testing.T) {
		err := p.Output(&errWriter{}, []byte(input))
		testy.Error(t, "errWriter: write error", err)
	})
}
