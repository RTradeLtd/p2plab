package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParser(t *testing.T) {
	//var edef metadata.ExperimentDefinition
	data, err := ioutil.ReadFile("../cue.mod/p2plab.cue")
	if err != nil {
		t.Fatal(err)
	}
	parser := NewParser([]string{string(data)})
	/* this currently fails
	 */data, err = ioutil.ReadFile("../cue.mod/p2plab_example.cue")
	if err != nil {
		t.Fatal(err)
	}
	pinst, err := parser.Compile("p2plab_example", string(data))
	if err != nil {
		t.Fatal(err)
	}
	_, err = pinst.TrialsToDefinition()
	if err != nil {
		t.Fatal(err)
	}
	return

	val := pinst.GetGroups()
	if val.Err() != nil {
		t.Fatal(err)
	}
	val = pinst.GetScenario()
	if val.Err() != nil {
		t.Fatal(err)
	}
	val = pinst.GetObjects()
	if val.Err() != nil {
		t.Fatal(err)
	}
	edf, err := pinst.ToExperimentDefinition()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("START\n%+v\nEND\n", edf)
}
