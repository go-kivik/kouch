package io

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/cobra"
)

const (
	// FlagOutputFile specifies where to write output.
	FlagOutputFile   = "output"
	flagOutputFormat = "output-format"
	// flagClobber indicates whether output files should be overwritten
	flagClobber = "force"
)

type defaultMode bool

func (m defaultMode) isDefault() bool {
	return bool(m)
}

var outputModes = make(map[string]outputMode)

func registerOutputMode(name string, m outputMode) {
	if _, ok := outputModes[name]; ok {
		panic(fmt.Sprintf("Output mode '%s' already registered", name))
	}
	outputModes[name] = m
}

// AddFlags adds command line flags for all configured output modes.
func AddFlags(cmd *cobra.Command) {
	defaults := make([]string, 0)
	formats := make([]string, 0, len(outputModes))
	for name, mode := range outputModes {
		if mode.isDefault() {
			defaults = append(defaults, name)
		}
		mode.config(cmd)
		formats = append(formats, name)
	}
	if len(defaults) == 0 {
		panic("No default output mode configured")
	}
	if len(defaults) > 1 {
		panic(fmt.Sprintf("Multiple default output modes configured: %s", strings.Join(defaults, ", ")))
	}
	sort.Strings(formats)
	pf := cmd.PersistentFlags()
	pf.StringP(flagOutputFormat, "F", defaults[0], fmt.Sprintf("Specify output format. Available options: %s", strings.Join(formats, ", ")))
	pf.StringP(FlagOutputFile, "o", "-", "Output destination. Use '-' for stdout")
	pf.BoolP(flagClobber, "", false, "Overwrite destination files")
}

// SelectOutput returns an io.Writer for the output.
func SelectOutput(cmd *cobra.Command) (io.Writer, error) {
	output, err := cmd.Flags().GetString(FlagOutputFile)
	if err != nil {
		return nil, err
	}
	if output == "" || output == "-" {
		// Default to stdout
		return os.Stdout, nil
	}
	clobber, err := cmd.Flags().GetBool(flagClobber)
	if err != nil {
		return nil, err
	}

	return &delayedOpenWriter{
		filename: output,
		clobber:  clobber,
	}, nil
}

// SelectOutputProcessor selects and configures the desired output processor
// based on the flags provided in cmd.
func SelectOutputProcessor(cmd *cobra.Command) (OutputProcessor, error) {
	name, err := cmd.Flags().GetString(flagOutputFormat)
	if err != nil {
		return nil, err
	}
	processor, ok := outputModes[name]
	if !ok {
		return nil, errors.Errorf("Unrecognized output format '%s'", name)
	}
	p, err := processor.new(cmd)
	return &errWrapper{p}, err
}

type outputMode interface {
	// config sets flags for the passed command, at start-up
	config(*cobra.Command)
	// isDefault returns true if this should be the default format. Exactly one
	// output mode must return true.
	isDefault() bool
	// new takes cmd, after command line options have been parsed, and returns
	// a new output processor.
	new(*cobra.Command) (OutputProcessor, error)
}

// OutputProcessor processes a command's output for display to a user.
type OutputProcessor interface {
	Output(io.Writer, io.ReadCloser) error
}
