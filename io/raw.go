package io

import (
	"io"

	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/pflag"
)

type rawMode struct{}

var _ kouchio.OutputMode = &rawMode{}

func (m *rawMode) AddFlags(_ *pflag.FlagSet) {}

func (m *rawMode) New(_ *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	return w, nil
}
