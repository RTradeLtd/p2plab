package parser

import (
	"cuelang.org/go/cue"
)

// Parser bundles the cue runtime with helper functions
// to enable parsing of cue source files
type Parser struct {
	runtime cue.Runtime
}

// Compile is used to compile the given cue source into our runtime
func (p *Parser) Compile(name string, data interface{}) (*cue.Instance, error) {
	return p.runtime.Compile(name, data)
}
