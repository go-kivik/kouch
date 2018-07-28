package io

import (
	"bytes"
	"io"

	"github.com/spf13/cobra"
)

func init() {
	registerOutputMode("raw", &rawMode{})
}

type rawMode struct{}

var _ outputMode = &rawMode{}

func (m *rawMode) config(cmd *cobra.Command) {}

func (m *rawMode) new(cmd *cobra.Command) (processor, error) {
	return &rawProcessor{}, nil
}

type rawProcessor struct{}

var _ processor = &rawProcessor{}

func (p *rawProcessor) Output(o io.Writer, input []byte) error {
	_, err := io.Copy(o, bytes.NewReader(input))
	return err
}
