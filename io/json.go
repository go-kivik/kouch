package io

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
)

type jsonMode struct{}

var _ outputMode = &jsonMode{}

func (m *jsonMode) config(cmd *cobra.Command) {
	cmd.PersistentFlags().String("json-prefix", "", "Prefix to begin each line of the JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	cmd.PersistentFlags().String("json-indent", "", "Indentation string for JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	cmd.PersistentFlags().Bool("json-escape-html", false, "Enable escaping of special HTML characters. See [https://golang.org/pkg/encoding/json/#Encoder.SetEscapeHTML].")
}

func (m *jsonMode) new(cmd *cobra.Command) processor {
	return &jsonProcessor{}
}

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
