// Package util contains random junk, for internal use only. Thus I'm willing
// to use a grab-back naming scheme, until something better pops up
package util

import (
	"context"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	kio "github.com/go-kivik/kouch/io"
)

// ChttpDo performs an HTTP request (GET is downgraded to HEAD if
// body is nil), writing the header to head, and body to body. If either head or body is nil, that write is skipped.
func ChttpDo(ctx context.Context, method, path string, o *kouch.Options) error {
	head, body := kouch.HeadDumper(ctx), kouch.Output(ctx)
	defer close(head) // nolint: errcheck
	defer close(body) // nolint: errcheck
	nilBody := body == nil || reflect.ValueOf(body).IsNil()
	nilHead := head == nil || reflect.ValueOf(head).IsNil()
	c, err := o.NewClient()
	if err != nil {
		return err
	}

	if method == http.MethodGet && nilBody {
		method = http.MethodHead
	}

	res, err := c.DoReq(ctx, method, path, o.Options)
	if err != nil {
		return err
	}
	if err = chttp.ResponseError(res); err != nil {
		return err
	}
	defer res.Body.Close() // nolint: errcheck

	if !nilHead {
		if e := res.Header.Write(head); e != nil {
			return e
		}
		// If head and body go to the same place, output a blank line between them
		if sameFd(head, body) {
			if _, e := head.Write([]byte("\r\n")); e != nil {
				return e
			}
		} else {
			_ = close(head)
		}
	}

	if !nilBody {
		if e := CopyAll(body, res.Body); e != nil {
			return e
		}
	}

	return nil
}

func sameFd(w1, w2 io.Writer) bool {
	if w1 == nil || w2 == nil {
		return false
	}
	if w1 == w2 {
		return true
	}
	u1 := kio.Underlying(w1)
	u2 := kio.Underlying(w2)
	if u1 == u2 {
		return true
	}
	f1, _ := u1.(*os.File)
	if f1 == nil {
		return false
	}
	f2, _ := u2.(*os.File)
	if f2 == nil {
		return false
	}
	return f1.Fd() == f2.Fd()
}
