package io

import (
	"io"

	"github.com/spf13/cobra"
)

type format int

// The availale output formats
const (
	FormatRaw format = iota
	FormatJSON
	FormatYAML
	FormatGoTmpl
)

// Output outputs result and/or error according to the configuration found in
// cmd.
func Output(cmd *cobra.Command, result []byte, err error) {

}

type outputMode interface {
	config(*cobra.Command)
	new(*cobra.Command) (processor, error)
}

type processor interface {
	Output(io.Writer, []byte) error
}
