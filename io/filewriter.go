package io

import (
	"io"
	"os"
)

type delayedOpenWriter struct {
	filename string
	clobber  bool
	w        io.Writer
}

var _ io.Writer = &delayedOpenWriter{}

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

func (w *delayedOpenWriter) open() (io.Writer, error) {
	if w.clobber {
		return os.Create(w.filename)
	}
	return os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
}
