package parser

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
)

// Parser bundles the cue runtime with helper functions
// to enable parsing of cue source files
type Parser struct {
	entrypoints    []string
	buildInstances []*build.Instance
	runtime        *cue.Runtime
}

// NewParser returns a ready to use cue parser
func NewParser(entrypoints []string) *Parser {
	p := &Parser{
		entrypoints: entrypoints,
		runtime:     new(cue.Runtime),
	}
	p.init()
	return p
}

func (p *Parser) init() {
	p.buildInstances = load.Instances(p.entrypoints, nil)
	for _, instance := range p.buildInstances {
		if instance.Incomplete {
			panic("failed to load build instances")
		}
	}
}

// Compile is used to compile the given cue source into our runtime
func (p *Parser) Compile(name string, data interface{}) (*cue.Instance, error) {
	return p.runtime.Compile(name, data)
}
