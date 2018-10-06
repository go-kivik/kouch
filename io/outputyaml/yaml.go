package outputyaml

import (
	"io"

	"github.com/go-kivik/kouch/io/outputcommon"
	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// YAMLMode outputs as YAML.
type YAMLMode struct{}

var _ kouchio.OutputMode = &YAMLMode{}

// AddFlags does nothing.
func (m *YAMLMode) AddFlags(_ *pflag.FlagSet) {}

// New returns a new YAML outputter.
func (m *YAMLMode) New(_ *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	return outputcommon.NewProcessor(w, func(o io.Writer, i interface{}) error {
		return yaml.NewEncoder(o).Encode(i)
	}), nil
}
