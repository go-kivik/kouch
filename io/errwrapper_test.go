package io

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

func TestErrWrapper(t *testing.T) {
	tests := []struct {
		name   string
		w      io.Writer
		r      io.ReadCloser
		p      kouch.OutputProcessor
		err    string
		status int
	}{
		{
			name: "no errors",
			w:    &bytes.Buffer{},
			r:    ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
			p:    &rawProcessor{},
		},
		{
			name:   "read error",
			w:      &bytes.Buffer{},
			r:      ioutil.NopCloser(&errReader{}),
			p:      &rawProcessor{},
			err:    "errReader: read error",
			status: chttp.ExitReadError,
		},
		{
			name:   "write error",
			w:      &errWriter{},
			r:      ioutil.NopCloser(strings.NewReader(`{"foo":123}`)),
			p:      &rawProcessor{},
			err:    "errWriter: write error",
			status: chttp.ExitWriteError,
		},
		{
			name:   "invalid JSON error",
			w:      &bytes.Buffer{},
			r:      ioutil.NopCloser(strings.NewReader(`oink`)),
			p:      &jsonProcessor{},
			err:    "invalid character 'o' looking for beginning of value",
			status: chttp.ExitWeirdReply,
		},
		{
			name:   "invalid JSON error",
			w:      &bytes.Buffer{},
			r:      ioutil.NopCloser(strings.NewReader(`oink`)),
			p:      &jsonProcessor{},
			err:    "invalid character 'o' looking for beginning of value",
			status: chttp.ExitWeirdReply,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &errWrapper{test.p}
			err := p.Output(test.w, test.r)
			testy.ExitStatusError(t, test.err, test.status, err)
		})
	}
}
