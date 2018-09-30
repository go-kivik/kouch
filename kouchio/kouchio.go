package kouchio

import "io"

// CloseWriter closes w if it is an io.WriteCloser; else it does nothing
func CloseWriter(w io.Writer) error {
	if wc, ok := w.(io.WriteCloser); ok {
		return wc.Close()
	}
	return nil
}
