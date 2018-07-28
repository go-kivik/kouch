package io

import (
	"io"

	"github.com/spf13/cobra"
)

func init() {
	registerOutputMode("raw", &rawMode{})
}

type rawMode struct {
	defaultMode
}

var _ outputMode = &rawMode{}

func (m *rawMode) config(cmd *cobra.Command) {}

func (m *rawMode) new(cmd *cobra.Command) (OutputProcessor, error) {
	return &rawProcessor{}, nil
}

type rawProcessor struct{}

var _ OutputProcessor = &rawProcessor{}

func (p *rawProcessor) Output(o io.Writer, input io.ReadCloser) error {
	defer input.Close()
	_, err := io.Copy(o, input)
	return err
}
