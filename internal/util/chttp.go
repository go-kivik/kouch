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

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// ChttpDo performs an HTTP request (GET is downgraded to HEAD if
// body is nil), writing the header to head, and body to body. If either head or body is nil, that write is skipped.
func ChttpDo(ctx context.Context, method, path string, o *kouch.Options) error {
	head, body := kouch.HeadDumper(ctx), kouch.Output(ctx)
	nilBody := isNil(body)
	defer close(head) // nolint: errcheck
	defer close(body) // nolint: errcheck
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

	if e := writeHead(head, res, !sameFd(head, body)); e != nil {
		return e
	}

	if nilBody {
		return nil
	}

	return CopyAll(body, res.Body)
}

// when closeHead is true, head is closed before return
func writeHead(head io.WriteCloser, res *http.Response, closeHead bool) error {
	if isNil(head) {
		return nil
	}
	if e := res.Header.Write(head); e != nil {
		return e
	}
	if closeHead {
		return head.Close()
	}
	_, err := head.Write([]byte("\r\n"))
	return err
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
