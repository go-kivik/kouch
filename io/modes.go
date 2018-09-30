package io

import (
	"github.com/go-kivik/kouch/io/outputjson"
	"github.com/go-kivik/kouch/kouchio"
)

const defaultOutputMode = "json"

var outputModes = map[string]kouchio.OutputMode{
	defaultOutputMode: &outputjson.JSONMode{},
	"yaml":            &yamlMode{},
	"raw":             &rawMode{},
	"template":        &tmplMode{},
}
