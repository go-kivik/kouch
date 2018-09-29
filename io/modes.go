package io

import (
	"io"

	"github.com/spf13/pflag"
)

var outputModes = make(map[string]outputMode)

type outputMode interface {
	// config sets flags for the passed command, at start-up
	config(*pflag.FlagSet)
	// isDefault returns true if this should be the default format. Exactly one
	// output mode must return true.
	isDefault() bool
	// new takes flags, after command line options have been parsed, and returns
	// a new output processor.
	new(*pflag.FlagSet, io.Writer) (io.Writer, error)
}
