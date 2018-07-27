package io

import (
	"bytes"
	"io"
)

type rawProcessor struct{}

var _ processor = &rawProcessor{}

func (p *rawProcessor) Output(o io.Writer, input []byte) error {
	_, err := io.Copy(o, bytes.NewReader(input))
	return err
}
