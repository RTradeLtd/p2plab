package parser

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
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

// to prevent any issues we cope the byte slice
func (p *Parser) SourceFromBytes(data []byte) load.Source {
	return load.FromBytes(append(data[0:0:0], data...))
}
