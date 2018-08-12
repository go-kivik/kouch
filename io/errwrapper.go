package io

import (
	"io"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/errors"
)

type errWrapper struct {
	kouch.OutputProcessor
}

var _ kouch.OutputProcessor = &errWrapper{}

func (p *errWrapper) Output(w io.Writer, r io.ReadCloser) error {
	return p.OutputProcessor.Output(&exitStatusWriter{w}, &exitStatusReadCloser{r})
}

type exitStatusWriter struct {
	io.Writer
}

var _ io.Writer = &exitStatusWriter{}

func (w *exitStatusWriter) Write(p []byte) (int, error) {
	c, err := w.Writer.Write(p)
	if err == nil {
		return c, nil
	}
	return c, &errors.ExitError{Err: err, ExitCode: chttp.ExitWriteError}
}

type exitStatusReadCloser struct {
	io.ReadCloser
}

var _ io.ReadCloser = &exitStatusReadCloser{}

func (r *exitStatusReadCloser) Read(p []byte) (int, error) {
	c, err := r.ReadCloser.Read(p)
	if err == nil || err == io.EOF {
		return c, err
	}
	return c, &errors.ExitError{Err: err, ExitCode: chttp.ExitReadError}
}
