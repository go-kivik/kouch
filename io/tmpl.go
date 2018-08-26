package io

import (
	"html/template"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	registerOutputMode("template", &tmplMode{})
}

const (
	optTemplate     = "template"
	optTemplateFile = "template-file"
)

type tmplMode struct {
	defaultMode
}

var _ outputMode = &tmplMode{}

func (m *tmplMode) config(flags *pflag.FlagSet) {
	flags.String(optTemplate, "", "Template string to use with -o=go-template. See [http://golang.org/pkg/text/template/#pkg-overview] for format documetation.")
	flags.String(optTemplateFile, "", "Template file to use with -o=go-template. Alternative to --template.")
}

func (m *tmplMode) new(cmd *cobra.Command, w io.Writer) (io.WriteCloser, error) {
	templateString, err := cmd.Flags().GetString(optTemplate)
	if err != nil {
		return nil, err
	}
	templateFile, err := cmd.Flags().GetString(optTemplateFile)
	if err != nil {
		return nil, err
	}
	if templateString == "" && templateFile == "" {
		return nil, errors.Errorf("Must provide --%s or --%s option", optTemplate, optTemplateFile)
	}
	if templateString != "" && templateFile != "" {
		return nil, errors.Errorf("Both --%s and --%s specified; must provide only one.", optTemplate, optTemplateFile)
	}
	if templateString != "" {
		tmpl, e := template.New("").Parse(templateString)
		return &tmplProcessor{template: tmpl, underlying: w}, e
	}
	tmpl, err := template.New("").ParseFiles(templateFile)
	return &tmplProcessor{template: tmpl, underlying: w}, err
}

type tmplProcessor struct {
	template   *template.Template
	underlying io.Writer
	r          *io.PipeReader
	w          *io.PipeWriter
	done       <-chan struct{}
	err        error
}

var _ io.WriteCloser = &tmplProcessor{}

func (p *tmplProcessor) Write(in []byte) (int, error) {
	if p.w == nil {
		p.init()
	}
	n, e := p.w.Write(in)
	return n, e
}

func (p *tmplProcessor) init() {
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
		p.err = p.template.Execute(p.underlying, unmarshaled)
	}()
}

func (p *tmplProcessor) Close() error {
	if p.w == nil {
		return nil
	}

	<-p.done
	_ = p.w.Close() // always returns nil for PipeWriter
	return p.err
}
