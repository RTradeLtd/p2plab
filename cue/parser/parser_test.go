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
	inst, err := parser.Compile("p2plab_example", string(data))
	if err != nil {
		t.Fatal(err)
	}
	val, err := parser.GetGroups(inst)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", val)
}
