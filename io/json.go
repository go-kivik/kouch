package io

import (
	"encoding/json"
	"io"
)

type jsonProcessor struct {
	prefix     string
	indent     string
	escapeHTML bool
}

var _ processor = &jsonProcessor{}

func (p *jsonProcessor) Output(o io.Writer, input []byte) error {
	var unmarshaled interface{}
	if err := json.Unmarshal(input, &unmarshaled); err != nil {
		return err
	}
	enc := json.NewEncoder(o)
	enc.SetIndent(p.prefix, p.indent)
	enc.SetEscapeHTML(p.escapeHTML)
	return enc.Encode(unmarshaled)
}
