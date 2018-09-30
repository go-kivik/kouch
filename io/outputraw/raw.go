package outputraw

import (
	"io"

	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/pflag"
)

// RawMode passes through the output, unaltered.
type RawMode struct{}

var _ kouchio.OutputMode = &RawMode{}

// AddFlags does nothing.
func (m *RawMode) AddFlags(_ *pflag.FlagSet) {}

// New returns w, unaltered.
func (m *RawMode) New(_ *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	return w, nil
}
