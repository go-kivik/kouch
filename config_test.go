package kouch

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
)

func TestDefaultCtx(t *testing.T) {
	tests := []struct {
		name     string
		conf     *Config
		expected *Context
		err      string
	}{
		{
			name: "no default context set",
			conf: &Config{},
			err:  "No default context",
		},
		{
			name: "No matching context",
			conf: &Config{DefaultContext: "foo"},
			err:  "Default context 'foo' not defined",
		},
		{
			name: "Success",
			conf: &Config{DefaultContext: "foo",
				Contexts: []NamedContext{
					{Name: "foo", Context: &Context{Root: "foo.com"}},
				}},
			expected: &Context{Root: "foo.com"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, err := test.conf.DefaultCtx()
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, ctx); d != nil {
				t.Error(d)
			}
		})
	}
}
