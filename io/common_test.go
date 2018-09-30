package io

import (
	"errors"
	"io"
)

type errWriter struct{}

var _ io.Writer = &errWriter{}

func (w *errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("errWriter: write error")
}

type errReader struct{}

var _ io.Reader = &errReader{}

func (r *errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("errReader: read error")
}
