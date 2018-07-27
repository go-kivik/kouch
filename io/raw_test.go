package io

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
)

type errWriter struct{}

var _ io.Writer = &errWriter{}

func (w *errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("errWriter: write error")
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
