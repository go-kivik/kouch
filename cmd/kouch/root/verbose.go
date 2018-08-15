package root

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"os"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

func verbose(ctx context.Context, cmd *cobra.Command) (context.Context, error) {
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		return ctx, err
	}
	if !verbose {
		return ctx, nil
	}
	ctx = kouch.SetVerbose(ctx, true)
	ctx = trace(ctx, os.Stderr)
	return ctx, nil
}

func trace(ctx context.Context, out io.Writer) context.Context {
	ct := &clientTrace{out: out}
	ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		GetConn:      ct.getConn,
		GotConn:      ct.gotConn,
		WroteRequest: ct.wroteRequest,
	})
	ctx = chttp.WithClientTrace(ctx, &chttp.ClientTrace{
		HTTPRequestBody: ct.httpRequestBody,
		HTTPResponse:    ct.httpResponse,
	})
	return ctx
}

type clientTrace struct {
	out        io.Writer
	uploadSize int64
	req        *http.Request
}

func (ct *clientTrace) getConn(hostPort string) {
	fmt.Fprintf(ct.out, "*   Trying %s...\n", hostPort)
}

func (ct *clientTrace) gotConn(info httptrace.GotConnInfo) {
	host, port, _ := net.SplitHostPort(info.Conn.RemoteAddr().String())
	fmt.Fprintf(ct.out, "* Connected to %s port %s\n", host, port)
}

func (ct *clientTrace) httpRequestBody(r *http.Request) {
	ct.req = r
	if r.Body != nil {
		defer r.Body.Close()
		ct.uploadSize, _ = io.Copy(ioutil.Discard, r.Body)
	}
}

func (ct *clientTrace) wroteRequest(i httptrace.WroteRequestInfo) {
	if i.Err != nil {
		fmt.Fprintf(ct.out, "Error writing request: %s\n", i.Err)
		return
	}
	dump, err := httputil.DumpRequest(ct.req, false)
	if err != nil {
		fmt.Fprintf(ct.out, "ERROR: %s\n", err)
	}
	ct.dump(">", dump)
	if ct.uploadSize > 0 {
		fmt.Fprintf(ct.out, "* upload completely sent off: %d of %[1]d bytes\n", ct.uploadSize)
	}
}

func (ct *clientTrace) httpResponse(r *http.Response) {
	dump, err := httputil.DumpResponse(r, false)
	if err != nil {
		fmt.Fprintf(ct.out, "ERROR: %s\n", err)
	}
	ct.dump("<", dump)
	fmt.Fprintf(ct.out, "* Closing connection\n")
}

func (ct *clientTrace) dump(prefix string, body []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		if text := scanner.Text(); text != "" {
			fmt.Fprintf(ct.out, "%s %s\n", prefix, text)
		} else {
			fmt.Fprintln(ct.out, prefix)
		}
	}
}
