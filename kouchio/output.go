package kouchio

import (
	"context"
	"io"

	"github.com/spf13/pflag"
)

// OutputMode is the common interface for all output modes.
type OutputMode interface {
	// AddFlags adds flags for the passed command, at start-up
	AddFlags(*pflag.FlagSet)
	// New takes flags, after command line options have been parsed, and returns
	// a new output processor.
	New(context.Context, io.Writer) (io.Writer, error)
}
