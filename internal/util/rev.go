package util

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

// FetchRev fetches the document revision from the server
func FetchRev(ctx context.Context, o *kouch.Options) (string, error) {
	c, err := chttp.New(ctx, o.Root)
	if err != nil {
		return "", nil
	}
	res, err := c.DoReq(ctx, http.MethodHead, DocPath(o), o.Options)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	return strings.Trim(res.Header.Get("Etag"), "\""), nil
}
