package kouch

import (
	"io"

	kio "github.com/go-kivik/kouch/io"
	"github.com/go-kivik/kouch/log"
)

// CmdContext is the command execution context.
type CmdContext struct {
	Logger   log.Logger
	Conf     *Config
	Output   io.Writer
	Outputer kio.OutputProcessor
}
