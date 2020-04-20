package parser

import "cuelang.org/go/cue"

// P2PLabInstance is a wrapper around a cue instance
// which exposes helper functions to reduce lookup verbosity
type P2PLabInstance struct {
	*cue.Instance
}

// GetGroups returns the groups in a cluster for the given instance
func (p *P2PLabInstance) GetGroups() (cue.Value, error) {
	value := p.Lookup("experiment").Lookup("cluster").Lookup("groups")
	return value, value.Err()
}

// GetScenario returns the scenario in an experiment fro the given instance
func (p *P2PLabInstance) GetScenario() (cue.Value, error) {
	value := p.Lookup("experiment").Lookup("scenario")
	return value, value.Err()
}

// GetObjects retunrs the objects to be used in an experiment from the given instance
func (p *P2PLabInstance) GetObjects() (cue.Value, error) {
	value := p.Lookup("experiment").Lookup("scenario").Lookup("objects")
	return value, value.Err()
}
