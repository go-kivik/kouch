package util

import "io"

// CopyAll copys from src to dst, and checks for any errors on close.
// If dst is nil, src is simply closed.
func CopyAll(dst io.Writer, src io.Reader) error {
	if dst == nil {
		return close(src)
	}
	_, err := io.Copy(dst, src)
	if e := close(src); e != nil && err == nil {
		err = e
	}
	if e := close(dst); e != nil && err == nil {
		err = e
	}
	return err
}

func close(x interface{}) error {
	if c, ok := x.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
