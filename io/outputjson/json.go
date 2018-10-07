package outputjson

import (
	"encoding/json"
	"io"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/io/outputcommon"
	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/pflag"
)

// JSONMode pretty-prints the JSON output.
type JSONMode struct{}

var _ kouchio.OutputMode = &JSONMode{}

// AddFlags adds JSON-specific flags.
func (m *JSONMode) AddFlags(flags *pflag.FlagSet) {
	flags.String(kouch.FlagJSONPrefix, "", "Prefix to begin each line of the JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	flags.String(kouch.FlagJSONIndent, "", "Indentation string for JSON output. See [https://golang.org/pkg/encoding/json/#Indent] for more information.")
	flags.Bool(kouch.FlagJSONEscapeHTML, false, "Enable escaping of special HTML characters. See [https://golang.org/pkg/encoding/json/#Encoder.SetEscapeHTML].")
}

// New returns a new output processor.
func (m *JSONMode) New(flags *pflag.FlagSet, w io.Writer) (io.Writer, error) {
	prefix, err := flags.GetString(kouch.FlagJSONPrefix)
	if err != nil {
		return nil, err
	}
	indent, err := flags.GetString(kouch.FlagJSONIndent)
	if err != nil {
		return nil, err
	}
	escapeHTML, err := flags.GetBool(kouch.FlagJSONEscapeHTML)
	if err != nil {
		return nil, err
	}
	return outputcommon.NewProcessor(w, func(o io.Writer, i interface{}) error {
		enc := json.NewEncoder(o)
		enc.SetIndent(prefix, indent)
		enc.SetEscapeHTML(escapeHTML)
		return enc.Encode(i)
	}), nil
}
