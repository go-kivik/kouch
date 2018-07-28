package io

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var outputModes = make(map[string]outputMode)

func registerOutputMode(name string, m outputMode) {
	if _, ok := outputModes[name]; ok {
		panic(fmt.Sprintf("Output mode '%s' already registered", name))
	}
	outputModes[name] = m
}

type outputMode interface {
	config(*cobra.Command)
	new(*cobra.Command) (processor, error)
}

type processor interface {
	Output(io.Writer, []byte) error
}
