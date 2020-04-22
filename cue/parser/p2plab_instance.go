package parser

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/Netflix/p2plab/metadata"
)

// P2PLabInstance is a wrapper around a cue instance
// which exposes helper functions to reduce lookup verbosity
type P2PLabInstance struct {
	*cue.Instance
}

// ToExperimentDefinition takes a cue instance and returns
// the experiment definition needed to process the experiment
func (p *P2PLabInstance) ToExperimentDefinition() (*metadata.ExperimentDefinition, error) {
	var (
		cedf metadata.ClusterDefinition
		sedf metadata.ScenarioDefinition
	)
	data, err := p.GetCluster().MarshalJSON()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &cedf); err != nil {
		return nil, err
	}
	iter, err := p.GetObjects().List()
	if err != nil {
		return nil, err
	}
	var objData []byte
	for iter.Next() {
		data, err := iter.Value().MarshalJSON()
		if err != nil {
			return nil, err
		}
		objData = append(objData, data...)
	}
	if err := json.Unmarshal(objData, &sedf.Objects); err != nil {
		return nil, err
	}
	data, err = p.GetSeed().MarshalJSON()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &sedf.Seed); err != nil {
		return nil, err
	}
	data, err = p.GetBenchmark().MarshalJSON()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &sedf.Benchmark); err != nil {
		return nil, err
	}
	return &metadata.ExperimentDefinition{
		ClusterDefinition:  cedf,
		ScenarioDefinition: sedf,
	}, nil
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

// GetTrials returns the trials, a mapping of clusters and scenarios to run together
func (p *P2PLabInstance) GetTrials() cue.Value {
	return p.GetExperiment().Lookup("trials")
}

// TrialsToDefinition returns a metadata.TrialDefinition
func (p *P2PLabInstance) TrialsToDefinition() (metadata.TrialDefinition, error) {
	def := metadata.TrialDefinition{}
	val := p.GetTrials()
	if val.Err() != nil {
		return def, val.Err()
	}
	iter, err := val.List()
	if err != nil {
		return def, err
	}
	for iter.Next() {
		val := iter.Value()
		if val.Err() != nil {
			return def, val.Err()
		}
		/*len, err := val.Len().Uint64()
		if err != nil {
			return def, err
		}
		*/
		iter2, err := val.List()
		if err != nil {
			return def, err
		}
		for iter2.Next() {
			val2 := iter2.Value()
			grps := val2.Lookup("groups")
			if grps.Err() != nil {
				continue
			} else {
				data, err := grps.MarshalJSON()
				if err != nil {
					return def, err
				}
				fmt.Println(string(data))
			}
			obj := val2.Lookup("objects")
			if obj.Err() != nil {
				continue
			} else {
				data, err := obj.MarshalJSON()
				if err != nil {
					return def, err
				}
				fmt.Println(string(data))
			}
		}
	}
	return def, nil
}
