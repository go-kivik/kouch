package io

import (
	"encoding/json"
	"html/template"
	"io"
)

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
