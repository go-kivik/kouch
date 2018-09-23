package io

import (
	"io"
	"os"
	"path/filepath"
)

type delayedOpenWriter struct {
	filename   string
	createDirs bool
	clobber    bool
	w          io.WriteCloser
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
	return openOutputFile(w.filename, w.clobber, w.createDirs)
}

func openOutputFile(filename string, clobber, createDirs bool) (*os.File, error) {
	if createDirs {
		if path := filepath.Dir(filename); path != "" {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}
	if clobber {
		return os.Create(filename)
	}
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
}
