package io

import (
	"bytes"
	"io"

	"github.com/spf13/cobra"
)

type rawMode struct{}

var _ outputMode = &rawMode{}

func (m *rawMode) config(cmd *cobra.Command) {}

func (m *rawMode) new(cmd *cobra.Command) processor {
	return &rawProcessor{}
}

type rawProcessor struct{}

var _ processor = &rawProcessor{}

func (p *rawProcessor) Output(o io.Writer, input []byte) error {
	_, err := io.Copy(o, bytes.NewReader(input))
	return err
}
