package parser

import "cuelang.org/go/cue"

// P2PLabInstance is a wrapper around a cue instance
// which exposes helper functions to reduce lookup verbosity
type P2PLabInstance struct {
	*cue.Instance
}

// GetExperiment returns the top-most cue value
func (p *P2PLabInstance) GetExperiment() cue.Value {
	return p.Lookup("experiment")
}

// GetCluster returns the cluster to be created as part of the benchmark
func (p *P2PLabInstance) GetCluster() cue.Value {
	return p.GetExperiment().Lookup("cluster")
}

// GetGroups returns the groups in a cluster for the given instance
func (p *P2PLabInstance) GetGroups() cue.Value {
	return p.GetCluster().Lookup("groups")
}

// GetScenario returns the scenario in an experiment fro the given instance
func (p *P2PLabInstance) GetScenario() cue.Value {
	return p.GetExperiment().Lookup("scenario")
}

// GetObjects retunrs the objects to be used in an experiment from the given instance
func (p *P2PLabInstance) GetObjects() cue.Value {
	return p.GetScenario().Lookup("objects")
}

// GetSeed returns the nodes to seed as part of the benchmark
func (p *P2PLabInstance) GetSeed() cue.Value {
	return p.GetScenario().Lookup("seed")
}

// GetBenchmark returns the benchmarks to run in p2plab
func (p *P2PLabInstance) GetBenchmark() cue.Value {
	return p.GetScenario().Lookup("benchmark")
}
