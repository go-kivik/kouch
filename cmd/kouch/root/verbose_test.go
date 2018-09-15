package root

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

func TestVerbose(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
		err      string
	}{
		{
			name:     "defaults",
			expected: false,
		},
		{
			name:     "verbose enabled",
			args:     []string{"--" + flagVerbose},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := rootCmd("1.2.3")
			cmd.ParseFlags(test.args)
			ctx, err := verbose(kouch.GetContext(cmd), cmd)
			testy.Error(t, test.err, err)
			if verbose := kouch.Verbose(ctx); verbose != test.expected {
				t.Errorf("Unexpected result: %t\n", verbose)
			}
		})
	}
}

func TestClientTrace(t *testing.T) {
	s := testy.ServeResponse(&http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Date": []string{"Wed, 15 Aug 2018 17:52:20 GMT"},
		},
		Body: ioutil.NopCloser(strings.NewReader("Test body")),
	})
	defer s.Close()

	buf := &bytes.Buffer{}
	ctx := trace(context.Background(), buf)

	target := &kouch.Target{Root: s.URL}
	c, err := target.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.DoReq(ctx, http.MethodGet, "/_testing", &chttp.Options{
		Body: ioutil.NopCloser(strings.NewReader("foo")),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = chttp.ResponseError(res); err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	sURL, _ := url.Parse(s.URL)
	host, port, _ := net.SplitHostPort(sURL.Host)

	expected := fmt.Sprintf(`*   Trying %s...
* Connected to %s port %s
> GET /_testing HTTP/1.1
> Host: %[1]s
> Accept: application/json
> Content-Type: application/json
> User-Agent: Kivik chttp/`+chttp.Version+` (Language=`+runtime.Version()+`; Platform=`+runtime.GOARCH+`/`+runtime.GOOS+`) Kouch/`+kouch.Version+`
>
* upload completely sent off: 3 of 3 bytes
< HTTP/1.1 200 OK
< Content-Length: 9
< Content-Type: text/plain; charset=utf-8
< Date: Wed, 15 Aug 2018 17:52:20 GMT
<
* Closing connection
`,
		sURL.Host, host, port)
	if d := diff.Text(expected, buf.String()); d != nil {
		t.Error(d)
	}
}
