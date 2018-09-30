package io

import (
	"io"

	"github.com/spf13/pflag"
)

const defaultOutputMode = "json"

var outputModes = map[string]outputMode{
	defaultOutputMode: &jsonMode{},
	"yaml":            &yamlMode{},
	"raw":             &rawMode{},
	"template":        &tmplMode{},
}

type outputMode interface {
	// AddFlags adds flags for the passed command, at start-up
	AddFlags(*pflag.FlagSet)
	// new takes flags, after command line options have been parsed, and returns
	// a new output processor.
	new(*pflag.FlagSet, io.Writer) (io.Writer, error)
}
