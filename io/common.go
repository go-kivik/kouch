package io

import (
	"encoding/json"
	"io"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/go-kivik/kouch/kouchio"
)

type processorFunc func(io.Writer, interface{}) error

// processor implements a basic processor
type processor struct {
	underlying io.Writer
	r          *io.PipeReader
	w          *io.PipeWriter
	done       <-chan struct{}
	err        error
	fn         processorFunc
}

var _ io.WriteCloser = &processor{}
var _ kouchio.WrappedWriter = &processor{}

func newProcessor(w io.Writer, fn processorFunc) io.WriteCloser {
	return &processor{
		underlying: w,
		fn:         fn,
	}
}

func (p *processor) Underlying() io.Writer {
	return p.underlying
}

func (p *processor) Write(in []byte) (int, error) {
	if p.w == nil {
		p.init()
	}
	n, e := p.w.Write(in)
	return n, e
}

func (p *processor) init() {
	p.r, p.w = io.Pipe()
	done := make(chan struct{})
	p.done = done
	go func() {
		defer func() { close(done) }()
		defer p.r.Close() // nolint: errcheck
		unmarshaled, err := unmarshal(p.r)
		if err != nil {
			p.err = err
			return
		}
		p.err = p.fn(p.underlying, unmarshaled)
	}()
}

func (p *processor) Close() error {
	if p.w == nil {
		return nil
	}

	<-p.done
	_ = p.w.Close() // always returns nil for PipeWriter
	return p.err
}

func unmarshal(r io.Reader) (interface{}, error) {
	var unmarshaled interface{}
	err := json.NewDecoder(r).Decode(&unmarshaled)
	return unmarshaled, errors.WrapExitError(chttp.ExitWeirdReply, err)
}
