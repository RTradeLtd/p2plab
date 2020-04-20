package parser

import (
	"cuelang.org/go/cue"
)

// Parser bundles the cue runtime with helper functions
// to enable parsing of cue source files
type Parser struct {
	entrypoints []string
	runtime     *cue.Runtime
}

// NewParser returns a ready to use cue parser
func NewParser(entrypoints []string) *Parser {
	p := &Parser{
		// entrypoints are the actual cue source files
		// we want to use when building p2plab cue files.
		// essentially they are the actual definitions
		// used to validate incoming cue source files
		entrypoints: entrypoints,
		runtime:     new(cue.Runtime),
	}
	return p
}

// Compile is used to compile the given cue source into our runtime
func (p *Parser) Compile(name string, cueSource string) (*cue.Instance, error) {
	// this is a temporary work around
	// until we can properly figure out the cue api
	for _, point := range p.entrypoints {
		cueSource += point
	}
	return p.runtime.Compile(name, cueSource)
}

// GetGroups returns the groups in a cluster for the given instance
func (p *Parser) GetGroups(inst *cue.Instance) (cue.Value, error) {
	value := inst.Lookup("experiment").Lookup("cluster").Lookup("groups")
	return value, value.Err()
}
