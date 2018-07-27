package io

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v2"
)

type yamlProcessor struct {
}

var _ processor = &yamlProcessor{}

func (p *yamlProcessor) Output(o io.Writer, input []byte) error {
	var unmarshaled interface{}
	if err := json.Unmarshal(input, &unmarshaled); err != nil {
		return err
	}
	enc := yaml.NewEncoder(o)
	return enc.Encode(unmarshaled)
}
