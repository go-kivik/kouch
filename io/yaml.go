package io

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type yamlMode struct{}

var _ outputMode = &yamlMode{}

func (m *yamlMode) config(cmd *cobra.Command) {}

func (m *yamlMode) new(cmd *cobra.Command) (processor, error) {
	return &yamlProcessor{}, nil
}

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
