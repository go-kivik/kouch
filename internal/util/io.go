package util

import "io"

type nopWriteCloser struct {
	io.Writer
}

var _ io.WriteCloser = &nopWriteCloser{}

// NopWriteCloser turns an io.Writer into an io.WriteCloser, with a noop
// Close() method.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	if w == nil {
		return nil
	}
	return &nopWriteCloser{w}
}

func (w *nopWriteCloser) Close() error { return nil }
