package kouch

import (
	"io"

	kio "github.com/go-kivik/kouch/io"
)

// CmdContext is the command execution context.
type CmdContext struct {
	Verbose  bool
	Conf     *Config
	Output   io.Writer
	Outputer kio.OutputProcessor
}
