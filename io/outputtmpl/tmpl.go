package outputtmpl

import (
	"context"
	"html/template"
	"io"
	"path/filepath"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/io/outputcommon"
	"github.com/go-kivik/kouch/kouchio"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// TmplMode outputs based on a provided template.
type TmplMode struct{}

var _ kouchio.OutputMode = &TmplMode{}

// AddFlags adds template-related flags
func (m *TmplMode) AddFlags(flags *pflag.FlagSet) {
	flags.String(kouch.FlagTemplate, "", "Template string to use with -o=go-template. See [http://golang.org/pkg/text/template/#pkg-overview] for format documetation.")
	flags.String(kouch.FlagTemplateFile, "", "Template file to use with -o=go-template. Alternative to --template.")
}

// New returns a new template outputter.
func (m *TmplMode) New(ctx context.Context, w io.Writer) (io.Writer, error) {
	flags := kouch.Flags(ctx)
	tmpl, err := newTmpl(flags)
	if err != nil {
		return nil, err
	}
	return outputcommon.NewProcessor(w, func(o io.Writer, i interface{}) error {
		return tmpl.Execute(o, i)
	}), nil
}

func newTmpl(flags *pflag.FlagSet) (*template.Template, error) {
	templateString, err := flags.GetString(kouch.FlagTemplate)
	if err != nil {
		return nil, err
	}
	templateFile, err := flags.GetString(kouch.FlagTemplateFile)
	if err != nil {
		return nil, err
	}
	if templateString == "" && templateFile == "" {
		return nil, errors.Errorf("Must provide --%s or --%s option", kouch.FlagTemplate, kouch.FlagTemplateFile)
	}
	if templateString != "" && templateFile != "" {
		return nil, errors.Errorf("Both --%s and --%s specified; must provide only one.", kouch.FlagTemplate, kouch.FlagTemplateFile)
	}
	if templateString != "" {
		return template.New("").Parse(templateString)
	}
	return template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
}
