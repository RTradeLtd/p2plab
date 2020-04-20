package parser

import (
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
	inst, err := parser.Compile("p2plab", data)
	if err != nil {
		t.Fatal(err)
	}
	if inst.Err != nil {
		t.Fatal(inst.Err)
	}
	/* this currently fails
	 */data, err = ioutil.ReadFile("../cue.mod/p2plab_example.cue")
	if err != nil {
		t.Fatal(err)
	}
	_, err = parser.Compile("p2plab_example", data)
	if err != nil {
		t.Fatal(err)
	}

}
