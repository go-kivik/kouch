package io

import (
	"encoding/json"
	"io"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	registerOutputMode("json", &jsonMode{defaultMode: true})
}

const (
	optJSONPrefix     = "json-prefix"
	optJSONIndent     = "json-indent"
	optJSONEscapeHTML = "json-escape-html"
)

type jsonMode struct {
	defaultMode
}

var _ outputMode = &jsonMode{}

func (m *jsonMode) config(flags *pflag.FlagSet) {
	flags.String(optJSONPrefix, "", "Prefix to begin each line of the JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	flags.String(optJSONIndent, "", "Indentation string for JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	flags.Bool(optJSONEscapeHTML, false, "Enable escaping of special HTML characters. See [https://golang.org/pkg/encoding/json/#Encoder.SetEscapeHTML].")
}

func (m *jsonMode) new(cmd *cobra.Command, w io.Writer) (io.WriteCloser, error) {
	prefix, err := cmd.Flags().GetString(optJSONPrefix)
	if err != nil {
		return nil, err
	}
	indent, err := cmd.Flags().GetString(optJSONIndent)
	if err != nil {
		return nil, err
	}
	escapeHTML, err := cmd.Flags().GetBool(optJSONEscapeHTML)
	if err != nil {
		return nil, err
	}
	return newJSONProcessor(prefix, indent, escapeHTML, w), nil
}

type jsonProcessor struct {
	prefix     string
	indent     string
	escapeHTML bool
	underlying io.Writer
	r          *io.PipeReader
	w          *io.PipeWriter
	done       <-chan struct{}
	err        error
}

var _ io.WriteCloser = &jsonProcessor{}

func newJSONProcessor(prefix, indent string, escapeHTML bool, w io.Writer) io.WriteCloser {
	return &jsonProcessor{
		prefix:     prefix,
		indent:     indent,
		escapeHTML: escapeHTML,
		underlying: w,
	}
}

func (p *jsonProcessor) Write(in []byte) (int, error) {
	if p.w == nil {
		p.init()
	}
	n, e := p.w.Write(in)
	return n, e
}

func (p *jsonProcessor) init() {
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
		enc := json.NewEncoder(p.underlying)
		enc.SetIndent(p.prefix, p.indent)
		enc.SetEscapeHTML(p.escapeHTML)
		p.err = enc.Encode(unmarshaled)
	}()
}

func (p *jsonProcessor) Close() error {
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
