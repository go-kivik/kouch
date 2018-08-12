package io

import (
	"io"

	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

func init() {
	registerOutputMode("yaml", &yamlMode{})
}

type yamlMode struct {
	defaultMode
}

var _ outputMode = &yamlMode{}

func (m *yamlMode) config(_ *pflag.FlagSet) {}

func (m *yamlMode) new(cmd *cobra.Command) (kouch.OutputProcessor, error) {
	return &yamlProcessor{}, nil
}

type yamlProcessor struct {
}

var _ kouch.OutputProcessor = &yamlProcessor{}

func (p *yamlProcessor) Output(o io.Writer, input io.ReadCloser) error {
	defer input.Close()
	unmarshaled, err := unmarshal(input)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(o)
	return enc.Encode(unmarshaled)
}
