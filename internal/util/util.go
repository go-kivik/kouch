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

// ChttpDo performs an HTTP request (GET is downgraded to HEAD if
// body is nil), writing the header to head, and body to body. If either head or body is nil, that write is skipped.
func ChttpDo(ctx context.Context, method, path string, o *kouch.Options, head, body io.Writer) error {
	nilBody := body == nil || reflect.ValueOf(body).IsNil()
	c, err := chttp.New(ctx, o.Root)
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
	defer res.Body.Close()
	if head != nil && !reflect.ValueOf(head).IsNil() {
		if e := res.Header.Write(head); e != nil {
			return e
		}
		if c, ok := head.(io.WriteCloser); ok {
			if e := c.Close(); e != nil {
				return e
			}
		}
	}
	if !nilBody {
		return CopyAll(body, res.Body)
	}
	return nil
}
