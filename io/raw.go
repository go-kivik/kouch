package io

import (
	"io"

	"github.com/spf13/pflag"
)

type rawMode struct{}

var _ outputMode = &rawMode{}

func (m *rawMode) AddFlags(_ *pflag.FlagSet) {}

func (m *rawMode) new(_ *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	return w, nil
}
