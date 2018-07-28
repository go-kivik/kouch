package io

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	AddFlags(cmd)

	testOptions(t, []string{"json-escape-html", "json-indent", "json-prefix", "output-format", "template", "template-file"}, cmd)
}
