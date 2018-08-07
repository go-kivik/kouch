package attachments

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

/*
func TestGetAttachments(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		resp    *http.Response
		content string
		err     string
		status  int
	}{
		{
			name:   "no filename",
			err:    "Must provide exactly one filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "too many filenames",
			args:   []string{"foo.txt", "bar.jpg"},
			err:    "Must provide exactly one filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "duplicate filenames",
			args:   []string{"--" + FlagFilename, "foo.txt", "foo.txt"},
			err:    "Must use --" + FlagFilename + " and pass separate filename",
			status: chttp.ExitFailedToInitialize,
		},
		// {
		// 	name: "no doc id provided",
		// 	args: []string{"foo.txt"},
		// 	err:  "document id must be provided as part of the filename argument or with the --" + FlagDocID + " flag",
		// },
		// {
		//     name: "Filename with slash",
		//     args: []string{"--filename", "foo/bar.txt"},
		//     err: "xxx",
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := &testy.ResponseHandler{Response: test.resp}
				h.ServeHTTP(w, r)
			}))
			conf := &kouch.Config{
				Contexts: []kouch.NamedContext{
					{
						Name:    "default",
						Context: &kouch.Context{Root: s.URL},
					},
				},
				DefaultContext: "default",
			}
			buf := &bytes.Buffer{}
			cx := &kouch.CmdContext{
				Conf:   conf,
				Output: buf,
			}
			root := &cobra.Command{}
			get := &cobra.Command{Use: "get"}
			get.AddCommand(attCmd(cx))
			root.AddCommand(get)
			root.ParseFlags(test.args)
			attCmd := attachmentCmd(cx)
			err := attCmd(root, nil)
			testy.ExitStatusErrorRE(t, test.err, test.status, err)
			if d := diff.Text(test.content, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}
*/

func TestGetAttachmentOpts(t *testing.T) {
	tests := []struct {
		name     string
		conf     *kouch.Config
		cmd      *cobra.Command
		args     []string
		expected *getAttOpts
		err      string
		status   int
	}{
		{
			name:   "no filename",
			args:   nil,
			err:    "Must provide exactly one filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "too many filenames",
			args:   []string{"foo.txt", "bar.jpg"},
			err:    "Must provide exactly one filename",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "duplicate filenames",
			args:   []string{"--" + FlagFilename, "foo.txt", "foo.txt"},
			err:    "Must use --" + FlagFilename + " and pass separate filename",
			status: chttp.ExitFailedToInitialize,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cx := &attCmdCtx{&kouch.CmdContext{
				Conf: test.conf,
			}}
			cmd := &cobra.Command{Use: "kouch"}
			cmd.AddCommand(attCmd(cx.CmdContext))
			cmd.ParseFlags(test.args)
			opts, err := cx.getAttachmentOpts(test.cmd, cmd.Flags().Args())
			testy.ExitStatusError(t, test.err, test.status, err)
			if d := diff.Interface(test.expected, opts); d != nil {
				t.Error(d)
			}
		})
	}
}
