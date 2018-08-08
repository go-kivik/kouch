package kouch

import (
	"io"

	"github.com/go-kivik/kouch/log"
)

// CmdContext is the command execution context.
type CmdContext struct {
	Logger   log.Logger
	Conf     *Config
	Output   io.Writer
	Outputer OutputProcessor
}

// OutputProcessor is a copy of kivik.io/OutputProcessor to prevent import cycles.
type OutputProcessor interface {
	Output(io.Writer, io.ReadCloser) error
}
