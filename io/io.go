package io

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/icza/dyno"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

const (
	// FlagOutputFile specifies where to write output.
	FlagOutputFile   = "output"
	flagOutputFormat = "output-format"
	// flagClobber indicates whether output files should be overwritten
	flagClobber = "force"
	flagStderr  = "stderr"
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
func AddFlags(flags *pflag.FlagSet) {
	defaults := make([]string, 0)
	formats := make([]string, 0, len(outputModes))
	for name, mode := range outputModes {
		if mode.isDefault() {
			defaults = append(defaults, name)
		}
		mode.config(flags)
		formats = append(formats, name)
	}
	if len(defaults) == 0 {
		panic("No default output mode configured")
	}
	if len(defaults) > 1 {
		panic(fmt.Sprintf("Multiple default output modes configured: %s", strings.Join(defaults, ", ")))
	}
	sort.Strings(formats)
	flags.StringP(flagOutputFormat, "F", defaults[0], fmt.Sprintf("Specify output format. Available options: %s", strings.Join(formats, ", ")))
	flags.StringP(FlagOutputFile, "o", "-", "Output destination. Use '-' for stdout")
	flags.BoolP(flagClobber, "", false, "Overwrite destination files")
	flags.String(flagStderr, "", `Where to redirect stderr (use "-" for stdout)`)

	flags.StringP(kouch.FlagData, kouch.FlagShortData, "", "HTTP request body data. Prefix with '@' to specify a filename.")
	flags.String(kouch.FlagDataJSON, "", "HTTP request body data, in JSON format.")
	flags.String(kouch.FlagDataYAML, "", "HTTP request body data, in YAML format.")
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
func SelectOutputProcessor(cmd *cobra.Command) (kouch.OutputProcessor, error) {
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
	config(*pflag.FlagSet)
	// isDefault returns true if this should be the default format. Exactly one
	// output mode must return true.
	isDefault() bool
	// new takes cmd, after command line options have been parsed, and returns
	// a new output processor.
	new(*cobra.Command) (kouch.OutputProcessor, error)
}

// RedirStderr redirects stderr based on configuration.
func RedirStderr(flags *pflag.FlagSet) error {
	filename, err := flags.GetString(flagStderr)
	if err != nil {
		return err
	}
	if filename == "" {
		return nil
	}
	if filename == "-" {
		os.Stderr = os.Stdout
		return nil
	}
	clobber, err := flags.GetBool(flagClobber)
	if err != nil {
		return err
	}
	f, err := openOutputFile(filename, clobber)
	if err != nil {
		return &errors.ExitError{
			Err:      err,
			ExitCode: chttp.ExitWriteError,
		}
	}
	os.Stderr = f
	return nil
}

// whichInput returns the input flag which was set, and the flag value
func whichInput(cmd *cobra.Command) (flag, value string, err error) {
	var found int
	for _, f := range []string{kouch.FlagData, kouch.FlagDataJSON, kouch.FlagDataYAML} {
		v, err := cmd.Flags().GetString(f)
		if err != nil {
			return "", "", err
		}
		if v != "" {
			found++
			flag = f
			value = v
		}
	}
	if found > 1 {
		return "", "", errors.NewExitError(chttp.ExitFailedToInitialize, "Only one data option may be provided")
	}
	return flag, value, nil
}

// SelectInput returns an io.ReadCloser for the input.
func SelectInput(cmd *cobra.Command) (io.ReadCloser, error) {
	flag, data, err := whichInput(cmd)
	if err != nil {
		return nil, err
	}
	if data == "" {
		// Default to stdin
		return os.Stdin, nil
	}
	var in io.ReadCloser
	if data[0] == '@' {
		var err error
		in, err = os.Open(data[1:])
		if err != nil {
			return nil, errors.WrapExitError(chttp.ExitReadError, err)
		}
	} else {
		in = ioutil.NopCloser(strings.NewReader(data))
	}
	if flag == kouch.FlagData {
		return in, nil
	}
	var i interface{}
	defer in.Close()
	switch flag {
	case kouch.FlagDataJSON:
		if err := json.NewDecoder(in).Decode(&i); err != nil {
			return nil, errors.WrapExitError(chttp.ExitPostError, err)
		}
	case kouch.FlagDataYAML:
		var j interface{}
		if err := yaml.NewDecoder(in).Decode(&j); err != nil {
			return nil, errors.WrapExitError(chttp.ExitPostError, err)
		}
		i = dyno.ConvertMapI2MapS(j)
	default:
		panic("Unknown flag: " + flag)
	}
	r, w := io.Pipe()
	go func() {
		err := json.NewEncoder(w).Encode(i)
		w.CloseWithError(err)
	}()
	return r, nil
}
