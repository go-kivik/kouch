package io

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
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

func (m *jsonMode) config(cmd *cobra.Command) {
	cmd.PersistentFlags().String(optJSONPrefix, "", "Prefix to begin each line of the JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	cmd.PersistentFlags().String(optJSONIndent, "", "Indentation string for JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	cmd.PersistentFlags().Bool(optJSONEscapeHTML, false, "Enable escaping of special HTML characters. See [https://golang.org/pkg/encoding/json/#Encoder.SetEscapeHTML].")
}

func (m *jsonMode) new(cmd *cobra.Command) (processor, error) {
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
	return &jsonProcessor{
		prefix:     prefix,
		indent:     indent,
		escapeHTML: escapeHTML,
	}, nil
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
