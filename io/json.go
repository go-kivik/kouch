package io

import (
	"encoding/json"
	"io"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
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

func (m *jsonMode) new(cmd *cobra.Command) (OutputProcessor, error) {
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

var _ OutputProcessor = &jsonProcessor{}

func (p *jsonProcessor) Output(o io.Writer, input io.ReadCloser) error {
	defer input.Close()
	unmarshaled, err := unmarshal(input)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(o)
	enc.SetIndent(p.prefix, p.indent)
	enc.SetEscapeHTML(p.escapeHTML)
	return enc.Encode(unmarshaled)
}

func unmarshal(r io.Reader) (interface{}, error) {
	var unmarshaled interface{}
	if err := json.NewDecoder(r).Decode(&unmarshaled); err != nil {
		return nil, &errors.ExitError{Err: err, ExitCode: chttp.ExitWeirdReply}
	}
	return unmarshaled, nil
}
