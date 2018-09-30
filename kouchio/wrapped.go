package kouchio

import "io"

// WrappedWriter represents an io.Writerr wrapped by some logic.
type WrappedWriter interface {
	// Underlying returns the original, unwrapped, io.WriteCloser.
	Underlying() io.Writer
}

// Underlying returns the unwrapped io.Writer, or the original if it was not
// wrapped.
func Underlying(w io.Writer) io.Writer {
	if u, ok := w.(WrappedWriter); ok {
		return Underlying(u.Underlying())
	}
	return w
}
