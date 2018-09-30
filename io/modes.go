package io

import "github.com/go-kivik/kouch/kouchio"

const defaultOutputMode = "json"

var outputModes = map[string]kouchio.OutputMode{
	defaultOutputMode: &jsonMode{},
	"yaml":            &yamlMode{},
	"raw":             &rawMode{},
	"template":        &tmplMode{},
}
