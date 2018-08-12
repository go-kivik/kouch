package kouch

import (
	"io"
)

// CmdContext is the command execution context.
type CmdContext struct {
	Verbose  bool
	Conf     *Config
	Output   io.Writer
	Outputer OutputProcessor
}
