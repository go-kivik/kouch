package kouch

import "io"

// OutputProcessor processes a command's output for display to a user.
type OutputProcessor interface {
	Output(io.Writer, io.ReadCloser) error
}
