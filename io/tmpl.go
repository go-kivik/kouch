package io

import (
	"encoding/json"
	"html/template"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

func (m *tmplMode) config(cmd *cobra.Command) {
	cmd.PersistentFlags().String(optTemplate, "", "Template string to use with -o=go-template. See [http://golang.org/pkg/text/template/#pkg-overview] for format documetation.")
	cmd.PersistentFlags().String(optTemplateFile, "", "Template file to use with -o=go-template. Alternative to --template.")
}

func (m *tmplMode) new(cmd *cobra.Command) (processor, error) {
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
		return &tmplProcessor{template: tmpl}, e
	}
	tmpl, err := template.New("").ParseFiles(templateFile)
	return &tmplProcessor{template: tmpl}, err
}

type tmplProcessor struct {
	template *template.Template
}

var _ processor = &tmplProcessor{}

func (p *tmplProcessor) Output(o io.Writer, input io.ReadCloser) error {
	defer input.Close()
	var unmarshaled interface{}
	if err := json.NewDecoder(input).Decode(&unmarshaled); err != nil {
		return err
	}
	return p.template.Execute(o, unmarshaled)
}
