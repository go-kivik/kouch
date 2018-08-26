package io

import (
	"io"

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

func (m *yamlMode) new(cmd *cobra.Command, w io.Writer) (io.WriteCloser, error) {
	return newYAMLProcessor(w), nil
}

type yamlProcessor struct {
	underlying io.Writer
	r          *io.PipeReader
	w          *io.PipeWriter
	done       <-chan struct{}
	err        error
}

var _ io.WriteCloser = &yamlProcessor{}

func newYAMLProcessor(w io.Writer) io.WriteCloser {
	return &yamlProcessor{
		underlying: w,
	}
}

func (p *yamlProcessor) Write(in []byte) (int, error) {
	if p.w == nil {
		p.init()
	}
	n, e := p.w.Write(in)
	return n, e
}

func (p *yamlProcessor) init() {
	p.r, p.w = io.Pipe()
	done := make(chan struct{})
	p.done = done
	go func() {
		defer func() { close(done) }()
		defer p.r.Close()
		unmarshaled, err := unmarshal(p.r)
		if err != nil {
			p.err = err
			return
		}
		p.err = yaml.NewEncoder(p.underlying).Encode(unmarshaled)
	}()
}

func (p *yamlProcessor) Close() error {
	if p.w == nil {
		return nil
	}

	<-p.done
	_ = p.w.Close() // always returns nil for PipeWriter
	return p.err
}
