package io

import (
	"io"

	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type yamlMode struct{}

var _ kouchio.OutputMode = &yamlMode{}

func (m *yamlMode) AddFlags(_ *pflag.FlagSet) {}

func (m *yamlMode) New(_ *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	return newProcessor(w, func(o io.Writer, i interface{}) error {
		return yaml.NewEncoder(o).Encode(i)
	}), nil
}
