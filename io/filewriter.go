package io

import (
	"io"
	"os"
)

type delayedOpenWriter struct {
	filename string
	clobber  bool
	w        io.WriteCloser
}

var _ io.WriteCloser = &delayedOpenWriter{}

func (w *delayedOpenWriter) Write(p []byte) (int, error) {
	if w.w == nil {
		var err error
		w.w, err = w.open()
		if err != nil {
			return 0, err
		}
	}
	return w.w.Write(p)
}

func (w *delayedOpenWriter) Close() error {
	return w.w.Close()
}

func (w *delayedOpenWriter) open() (io.WriteCloser, error) {
	return openOutputFile(w.filename, w.clobber)
}

func openOutputFile(filename string, clobber bool) (*os.File, error) {
	if clobber {
		return os.Create(filename)
	}
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
}
