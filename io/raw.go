package io

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	registerOutputMode("raw", &rawMode{})
}

type rawMode struct {
	defaultMode
}

var _ outputMode = &rawMode{}

func (m *rawMode) config(_ *pflag.FlagSet) {}

func (m *rawMode) new(cmd *cobra.Command, w io.Writer) (io.WriteCloser, error) {
	if t, ok := w.(io.WriteCloser); ok {
		return t, nil
	}
	return &nopCloser{w}, nil
}
