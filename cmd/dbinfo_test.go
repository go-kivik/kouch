package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/kivik"
	"github.com/flimzy/testy"
)

func TestGetInfo(t *testing.T) {
	type Test struct {
		name           string
		server, dbname string

		expected string
		status   int
		err      string
		finish   func(*testing.T)
	}
	tests := []Test{
		{
			name:   "invalid url",
			server: "1.2.3.4:1",
			status: kivik.StatusBadRequest,
			err:    "parse 1.2.3.4:1: first path segment in URL cannot contain colon",
		},
		{
			name:   "connection refused",
			server: "http://127.0.0.1:1/",
			status: kivik.StatusNetworkError,
			err:    "Get http://127.0.0.1:1: dial tcp 127.0.0.1:1: getsockopt: connection refused",
		},
		func() Test {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/receipts" {
					http.Error(w, fmt.Sprintf("Unexpected path: %s", r.URL.Path), 400)
					return
				}
				w.Write([]byte(`{"cluster":{"n":3,"q":8,"r":2,"w":2},"compact_running":false,"data_size":65031503,"db_name":"receipts","disk_format_version":6,"disk_size":137433211,"doc_count":6146,"doc_del_count":64637,"instance_start_time":"0","other":{"data_size": 6982448},"purge_seq":0,"sizes":{"active":65031503,"external":66982448,"file":137433211},"update_seq":"292786-g1AAAAF..."}`))
			}))
			return Test{
				name:     "doc example",
				server:   s.URL,
				dbname:   "receipts",
				expected: `{"cluster":{"n":3,"q":8,"r":2,"w":2},"compact_running":false,"data_size":65031503,"db_name":"receipts","disk_format_version":6,"disk_size":137433211,"doc_count":6146,"doc_del_count":64637,"instance_start_time":"0","other":{"data_size": 6982448},"purge_seq":0,"sizes":{"active":65031503,"external":66982448,"file":137433211},"update_seq":"292786-g1AAAAF..."}`,
			}
		}(),
	}
	for _, test := range tests {
		func(test Test) {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				if test.finish != nil {
					defer func() {
						test.finish(t)
					}()
				}
				result, err := getInfo(context.Background(), test.server, test.dbname)
				testy.StatusError(t, test.err, test.status, err)
				if d := diff.Text(test.expected, string(result)); d != nil {
					t.Error(d)
				}
			})
		}(test)
	}
}
