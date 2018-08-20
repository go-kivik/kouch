// Package util contains random junk, for internal use only. Thus I'm willing
// to use a grab-back naming scheme, until something better pops up
package util

import (
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

// ChttpGet performs an HTTP GET (or HEAD if body is nil), writing the header to
// head, and body to body. If either head or body is nil, that write is skipped.
func ChttpGet(ctx context.Context, path string, o *kouch.Options, head, body io.Writer) error {
	nilBody := reflect.ValueOf(body).IsNil()
	c, err := chttp.New(ctx, o.Root)
	if err != nil {
		return err
	}

	method := http.MethodGet
	if nilBody {
		method = http.MethodHead
	}

	res, err := c.DoReq(ctx, method, path, o.Options)
	if err != nil {
		return err
	}
	if err = chttp.ResponseError(res); err != nil {
		return err
	}
	defer res.Body.Close()
	if !reflect.ValueOf(head).IsNil() {
		if e := res.Header.Write(head); e != nil {
			return e
		}
	}
	if !nilBody {
		if _, e := io.Copy(body, res.Body); e != nil {
			return e
		}
	}
	return nil
}
