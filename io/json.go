package io

import (
	"encoding/json"
	"io"

	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/pflag"
)

const (
	optJSONPrefix     = "json-prefix"
	optJSONIndent     = "json-indent"
	optJSONEscapeHTML = "json-escape-html"
)

type jsonMode struct{}

var _ kouchio.OutputMode = &jsonMode{}

func (m *jsonMode) AddFlags(flags *pflag.FlagSet) {
	flags.String(optJSONPrefix, "", "Prefix to begin each line of the JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	flags.String(optJSONIndent, "", "Indentation string for JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	flags.Bool(optJSONEscapeHTML, false, "Enable escaping of special HTML characters. See [https://golang.org/pkg/encoding/json/#Encoder.SetEscapeHTML].")
}

func (m *jsonMode) New(flags *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	prefix, err := flags.GetString(optJSONPrefix)
	if err != nil {
		return nil, err
	}
	indent, err := flags.GetString(optJSONIndent)
	if err != nil {
		return nil, err
	}
	escapeHTML, err := flags.GetBool(optJSONEscapeHTML)
	if err != nil {
		return nil, err
	}
	return newProcessor(w, func(o io.Writer, i interface{}) error {
		enc := json.NewEncoder(o)
		enc.SetIndent(prefix, indent)
		enc.SetEscapeHTML(escapeHTML)
		return enc.Encode(i)
	}), nil
}
