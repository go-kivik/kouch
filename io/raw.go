package io

import (
	"io"

	"github.com/go-kivik/kouch"
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

func (m *rawMode) new(cmd *cobra.Command) (kouch.OutputProcessor, error) {
	return &rawProcessor{}, nil
}

type rawProcessor struct{}

var _ kouch.OutputProcessor = &rawProcessor{}

func (p *rawProcessor) Output(o io.Writer, input io.ReadCloser) error {
	defer input.Close()
	_, err := io.Copy(o, input)
	return err
}
