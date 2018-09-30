package io

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/pkg/errors"
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

func (m *tmplMode) new(flags *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	tmpl, err := newTmpl(flags)
	if err != nil {
		return nil, err
	}
	return newProcessor(w, func(o io.Writer, i interface{}) error {
		return tmpl.Execute(o, i)
	}), nil
}

func newTmpl(flags *pflag.FlagSet) (*template.Template, error) {
	templateString, err := flags.GetString(optTemplate)
	if err != nil {
		return nil, err
	}
	templateFile, err := flags.GetString(optTemplateFile)
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
		return template.New("").Parse(templateString)
	}
	return template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
}
