package io

import (
	"io"

	"github.com/spf13/pflag"
)

type rawMode struct {
	defaultMode
}

var _ outputMode = &rawMode{}

func (m *rawMode) config(_ *pflag.FlagSet) {}

func (m *rawMode) new(_ *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	return w, nil
}
