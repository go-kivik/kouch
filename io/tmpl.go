package io

import (
	"encoding/json"
	"html/template"
	"io"

	"github.com/spf13/cobra"
)

type tmplMode struct{}

var _ outputMode = &tmplMode{}

func (m *tmplMode) config(cmd *cobra.Command) {
	cmd.PersistentFlags().String("template", "", "Template string to use with -o=go-template. See [http://golang.org/pkg/text/template/#pkg-overview] for format documetation.")
	cmd.PersistentFlags().String("template-file", "", "Template file to use with -o=go-template. Alternative to --template.")
}

func (m *tmplMode) new(cmd *cobra.Command) processor {
	return &tmplProcessor{}
}

type tmplProcessor struct {
	template string
}

var _ processor = &tmplProcessor{}

func (p *tmplProcessor) Output(o io.Writer, input []byte) error {
	var unmarshaled interface{}
	if err := json.Unmarshal(input, &unmarshaled); err != nil {
		return err
	}
	tmpl, err := template.New("").Parse(p.template)
	if err != nil {
		return err
	}
	return tmpl.Execute(o, unmarshaled)
}
