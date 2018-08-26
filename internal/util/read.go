package util

import "io"

// CopyAll copys from src to dst, and checks for any errors on close.
// If dst is nil, src is simply closed.
func CopyAll(dst io.WriteCloser, src io.ReadCloser) error {
	if dst == nil {
		return src.Close()
	}
	_, err := io.Copy(dst, src)
	if e := src.Close(); e != nil && err == nil {
		err = e
	}
	if e := dst.Close(); e != nil && err == nil {
		err = e
	}
	return err
}
